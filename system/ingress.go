package system

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/bentoml/yatai-common/consts"
	"github.com/bentoml/yatai-common/utils"
)

func GetIngressClassName(ctx context.Context, cliset *kubernetes.Clientset) (ingressClassName string, err error) {
	configMap, err := GetNetworkConfigConfigMap(ctx, cliset)
	if err != nil {
		err = errors.Wrapf(err, "failed to get configmap %s", consts.KubeConfigMapNameNetworkConfig)
		return
	}

	ingressClassName = strings.TrimSpace(configMap.Data[consts.NetworkConfigKeyIngressClass])
	if ingressClassName != "" {
		return
	}

	ingressClassName = consts.DefaultIngressClassName

	configMapCli := cliset.CoreV1().ConfigMaps(GetNamespace())
	_, err = configMapCli.Patch(ctx, consts.KubeConfigMapNameNetworkConfig, types.StrategicMergePatchType, []byte(fmt.Sprintf(`{"data":{"%s":"%s"}}`, consts.NetworkConfigKeyIngressClass, ingressClassName)), metav1.PatchOptions{})
	if err != nil {
		err = errors.Wrapf(err, "failed to patch configmap %s", consts.KubeConfigMapNameNetworkConfig)
		return
	}

	return
}

func GetIngressIP(ctx context.Context, cliset *kubernetes.Clientset) (ip string, err error) {
	ingressClassName, err := GetIngressClassName(ctx, cliset)
	if err != nil {
		return
	}

	ingressCli := cliset.NetworkingV1().Ingresses(GetNamespace())

	ingName := "default-domain-"
	pathType := networkingv1.PathTypeImplementationSpecific

	podName := os.Getenv("POD_NAME")
	if podName == "" {
		// random string
		podName = strings.ToLower(utils.RandString(10))
	}

	logrus.Infof("Creating ingress %s to get a ingress IP automatically", ingName)
	ing, err := ingressCli.Create(ctx, &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: ingName,
			Namespace:    GetNamespace(),
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClassName,
			Rules: []networkingv1.IngressRule{{
				Host: fmt.Sprintf("%s.default-domain.invalid", podName),
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
	if err != nil {
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

	domainSuffix = strings.TrimSpace(configMap.Data[consts.NetworkConfigKeyDomainSuffix])
	if domainSuffix != "" {
		logrus.Infof("The domain suffix has already set to %s", domainSuffix)
		return
	}

	magicDNS := GetMagicDNS()

	var ip string

	ip, err = GetIngressIP(ctx, cliset)
	if err != nil {
		return
	}

	domainSuffix = fmt.Sprintf("%s.%s", ip, magicDNS)

	configMapCli := cliset.CoreV1().ConfigMaps(GetNamespace())

	logrus.Infof("Setting domain suffix to %s", domainSuffix)
	_, err = configMapCli.Patch(ctx, consts.KubeConfigMapNameNetworkConfig, types.MergePatchType, []byte(fmt.Sprintf(`{"data":{"%s":"%s"}}`, consts.NetworkConfigKeyDomainSuffix, domainSuffix)), metav1.PatchOptions{})
	if err != nil {
		err = errors.Wrapf(err, "failed to patch configmap %s", consts.KubeConfigMapNameNetworkConfig)
		return
	}
	logrus.Infof("Domain suffix has been set to %s", domainSuffix)

	return
}
