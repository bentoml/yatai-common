package system

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/bentoml/yatai-common/consts"
)

func GetIngressAnnotations(ctx context.Context, cliset *kubernetes.Clientset) (annotations map[string]string, err error) {
	configMap, err := GetNetworkConfigConfigMap(ctx, cliset)
	if err != nil {
		err = errors.Wrapf(err, "failed to get configmap %s", consts.KubeConfigMapNameNetworkConfig)
		return
	}

	annotations = make(map[string]string)

	annotationsStr := strings.TrimSpace(configMap.Data[consts.KubeConfigMapKeyNetworkConfigIngressAnnotations])

	if annotationsStr == "" {
		return
	}

	err = json.Unmarshal([]byte(annotationsStr), &annotations)
	if err != nil {
		err = errors.Wrapf(err, "failed to json unmarshal %s in configmap %s: %s", consts.KubeConfigMapKeyNetworkConfigIngressAnnotations, consts.KubeConfigMapNameNetworkConfig, annotationsStr)
		return
	}

	return
}

func GetIngressClassName(ctx context.Context, cliset *kubernetes.Clientset) (ingressClassName *string, err error) {
	configMap, err := GetNetworkConfigConfigMap(ctx, cliset)
	if err != nil {
		err = errors.Wrapf(err, "failed to get configmap %s", consts.KubeConfigMapNameNetworkConfig)
		return
	}

	ingressClassName_ := strings.TrimSpace(configMap.Data[consts.KubeConfigMapKeyNetworkConfigIngressClass])

	if ingressClassName_ != "" {
		ingressClassName = &ingressClassName_
	}

	return
}

func GetIngressIP(ctx context.Context, cliset *kubernetes.Clientset) (ip string, err error) {
	ingressClassName, err := GetIngressClassName(ctx, cliset)
	if err != nil {
		return
	}

	ingressAnnotations, err := GetIngressAnnotations(ctx, cliset)
	if err != nil {
		return
	}

	ingressCli := cliset.NetworkingV1().Ingresses(GetNamespace())

	ingName := "default-domain-"
	pathType := networkingv1.PathTypeImplementationSpecific

	podName := os.Getenv("POD_NAME")
	if podName == "" {
		// random string
		guid := xid.New()
		podName = fmt.Sprintf("a%s", strings.ToLower(guid.String()))
	}

	logrus.Infof("Creating ingress %s to get a ingress IP automatically", ingName)
	ing, err := ingressCli.Create(ctx, &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: ingName,
			Namespace:    GetNamespace(),
			Annotations:  ingressAnnotations,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: ingressClassName,
			Rules: []networkingv1.IngressRule{{
				Host: fmt.Sprintf("%s.this-is-yatai-in-order-to-generate-the-default-domain-suffix.yeah", podName),
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: "default-domain-service",
										Port: networkingv1.ServiceBackendPort{
											Number: consts.BentoServicePort,
										},
									},
								},
							},
						},
					},
				},
			}},
		},
	}, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		err = errors.Wrapf(err, "failed to create ingress %s", ingName)
		return
	}
	defer func() {
		_ = ingressCli.Delete(ctx, ing.Name, metav1.DeleteOptions{})
	}()

	// Interval to poll for objects.
	pollInterval := 10 * time.Second
	// How long to wait for objects.
	waitTimeout := 20 * time.Minute

	logrus.Infof("Waiting for ingress %s to be ready", ing.Name)
	// Wait for the Ingress to be Ready.
	if err = wait.PollImmediate(pollInterval, waitTimeout, func() (done bool, err error) {
		ing, err = ingressCli.Get(
			ctx, ing.Name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return len(ing.Status.LoadBalancer.Ingress) > 0, nil
	}); err != nil {
		err = errors.Wrapf(err, "failed to wait for ingress %s to be ready", ing.Name)
		return
	}
	logrus.Infof("Ingress %s is ready", ing.Name)

	address := ing.Status.LoadBalancer.Ingress[0]

	ip = address.IP
	if ip == "" {
		if address.Hostname == "" {
			err = errors.Errorf("the ingress %s status has no IP or hostname", ing.Name)
			return
		}
		var ipAddr *net.IPAddr
		ipAddr, err = net.ResolveIPAddr("ip4", address.Hostname)
		if err != nil {
			err = errors.Wrapf(err, "failed to resolve ip address for hostname %s", address.Hostname)
			return
		}
		ip = ipAddr.String()
	}

	return
}

func GetDomainSuffix(ctx context.Context, cliset *kubernetes.Clientset) (domainSuffix string, err error) {
	configMap, err := GetNetworkConfigConfigMap(ctx, cliset)
	if err != nil {
		err = errors.Wrapf(err, "failed to get configmap %s", consts.KubeConfigMapNameNetworkConfig)
		return
	}

	domainSuffix = strings.TrimSpace(configMap.Data[consts.KubeConfigMapKeyNetworkConfigDomainSuffix])
	if domainSuffix != "" {
		logrus.Infof("The %s in the network config has already set to `%s`", consts.KubeConfigMapKeyNetworkConfigDomainSuffix, domainSuffix)
		return
	}

	magicDNS := GetMagicDNS()

	var ip string

	ip, err = GetIngressIP(ctx, cliset)
	if err != nil {
		return
	}

	domainSuffix = fmt.Sprintf("%s.%s", ip, magicDNS)

	logrus.Infof("you have not set the %s in the network config, so use magic DNS to generate a domain suffix automatically: `%s`, and set it to the network config", consts.KubeConfigMapKeyNetworkConfigDomainSuffix, domainSuffix)

	configMapCli := cliset.CoreV1().ConfigMaps(configMap.Namespace)
	_, err = configMapCli.Patch(ctx, configMap.Name, types.MergePatchType, []byte(fmt.Sprintf(`{"data":{"%s":"%s"}}`, consts.KubeConfigMapKeyNetworkConfigDomainSuffix, domainSuffix)), metav1.PatchOptions{})
	if err != nil {
		err = errors.Wrapf(err, "failed to patch configmap %s", consts.KubeConfigMapNameNetworkConfig)
		return
	}

	return
}
