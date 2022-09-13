package system

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/bentoml/yatai-common/consts"
)

func GetNetworkConfigConfigMap(ctx context.Context, cliset *kubernetes.Clientset) (configMap *corev1.ConfigMap, err error) {
	configMapCli := cliset.CoreV1().ConfigMaps(GetNamespace())
	configMap, err = configMapCli.Get(ctx, consts.KubeConfigMapNameNetworkConfig, metav1.GetOptions{})
	return
}
