package system

import (
	"context"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/bentoml/yatai-common/consts"
	"github.com/bentoml/yatai-common/k8sutils"
)

func GetNetworkConfigConfigMap(ctx context.Context, cliset *kubernetes.Clientset) (configMap *corev1.ConfigMap, err error) {
	configMapCli := cliset.CoreV1().ConfigMaps(GetNamespace())
	configMap, err = configMapCli.Get(ctx, consts.KubeConfigMapNameNetworkConfig, metav1.GetOptions{})
	isNotFound := k8serrors.IsNotFound(err)
	if err != nil && !isNotFound {
		err = errors.Wrapf(err, "failed to get configmap %s", consts.KubeConfigMapNameNetworkConfig)
		return
	}
	if isNotFound {
		ns := GetNamespace()
		err = k8sutils.MakesureNamespaceExists(ctx, cliset, ns)
		if err != nil {
			return
		}

		configMap, err = configMapCli.Create(ctx, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      consts.KubeConfigMapNameNetworkConfig,
				Namespace: ns,
			},
			Data: map[string]string{},
		}, metav1.CreateOptions{})
		if err != nil {
			err = errors.Wrapf(err, "failed to create configmap %s", consts.KubeConfigMapNameNetworkConfig)
			return
		}
	}
	return
}
