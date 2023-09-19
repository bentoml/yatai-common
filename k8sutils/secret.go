package k8sutils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/bentoml/yatai-common/config"
	"github.com/bentoml/yatai-common/consts"
)

func MakeSureDockerRegcred(ctx context.Context, secretGetter func(ctx context.Context, namespace, name string) (*corev1.Secret, error), cliset *kubernetes.Clientset, namespace string) (secret *corev1.Secret, err error) {
	dockerRegistry, err := config.GetDockerRegistryConfig(ctx, secretGetter)
	if err != nil {
		return
	}

	if dockerRegistry.Username == "" {
		return
	}

	secret, err = secretGetter(ctx, namespace, consts.KubeSecretNameRegcred)
	isNotFound := k8serrors.IsNotFound(err)
	if err != nil && !isNotFound {
		return
	}

	dockerConfig := struct {
		Auths map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
	}{
		Auths: map[string]struct {
			Auth string `json:"auth"`
		}{
			dockerRegistry.Server: {
				Auth: base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", dockerRegistry.Username, dockerRegistry.Password))),
			},
		},
	}

	var dockerConfigContent []byte
	dockerConfigContent, err = json.Marshal(&dockerConfig)
	if err != nil {
		return
	}

	if isNotFound {
		secret = &corev1.Secret{
			Type: corev1.SecretTypeDockerConfigJson,
			ObjectMeta: metav1.ObjectMeta{
				Name:      consts.KubeSecretNameRegcred,
				Namespace: namespace,
			},
			Data: map[string][]byte{
				".dockerconfigjson": dockerConfigContent,
			},
		}
		secret, err = cliset.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return
		}
	} else if string(secret.Data[".dockerconfigjson"]) != string(dockerConfigContent) {
		secret.Data[".dockerconfigjson"] = dockerConfigContent
		secret, err = cliset.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			return
		}
	}
	return
}
