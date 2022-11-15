package config

import (
	"context"
	"os"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/bentoml/yatai-common/consts"
)

func GetYataiSystemNamespaceFromEnv() string {
	return getEnv(consts.EnvYataiSystemNamespace, consts.DefaultKubeNamespaceYataiSystem)
}

func GetYataiImageBuilderNamespace(ctx context.Context, cliset *kubernetes.Clientset) (namespace string, err error) {
	namespace = os.Getenv(consts.EnvYataiImageBuilderNamespace)
	if namespace != "" {
		return
	}

	yataiSystemNamespace := GetYataiSystemNamespaceFromEnv()
	yataiImageBuilderSharedEnvSecretName := consts.KubeSecretNameYataiImageBuilderSharedEnv

	secret, err := cliset.CoreV1().Secrets(yataiSystemNamespace).Get(ctx, yataiImageBuilderSharedEnvSecretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = errors.Wrapf(err, "secret %s not found in namespace %s", yataiImageBuilderSharedEnvSecretName, yataiSystemNamespace)
		}
		return
	}

	namespace = string(secret.Data[consts.EnvYataiImageBuilderNamespace])
	if namespace == "" {
		namespace = consts.DefaultKubeNamespaceYataiImageBuilderComponent
	}

	return
}

func GetYataiDeploymentNamespace(ctx context.Context, cliset *kubernetes.Clientset) (namespace string, err error) {
	namespace = os.Getenv(consts.EnvYataiDeploymentNamespace)
	if namespace != "" {
		return
	}

	yataiSystemNamespace := GetYataiSystemNamespaceFromEnv()
	yataiDeploymentSharedEnvSecretName := consts.KubeSecretNameYataiDeploymentSharedEnv

	secret, err := cliset.CoreV1().Secrets(yataiSystemNamespace).Get(ctx, yataiDeploymentSharedEnvSecretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = errors.Wrapf(err, "secret %s not found in namespace %s", yataiDeploymentSharedEnvSecretName, yataiSystemNamespace)
		}
		return
	}

	namespace = string(secret.Data[consts.EnvYataiDeploymentNamespace])
	if namespace == "" {
		namespace = consts.DefaultKubeNamespaceYataiDeploymentComponent
	}

	return
}

func GetImageBuildersNamespace(ctx context.Context, cliset *kubernetes.Clientset) (namespace string, err error) {
	namespace = os.Getenv(consts.EnvImageBuildersNamespace)
	if namespace != "" {
		return
	}

	yataiSystemNamespace := GetYataiSystemNamespaceFromEnv()
	yataiImageBuilderSharedEnvSecretName := consts.KubeSecretNameYataiImageBuilderSharedEnv

	secret, err := cliset.CoreV1().Secrets(yataiSystemNamespace).Get(ctx, yataiImageBuilderSharedEnvSecretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = errors.Wrapf(err, "secret %s not found in namespace %s", yataiImageBuilderSharedEnvSecretName, yataiSystemNamespace)
		}
		return
	}

	namespace = string(secret.Data[consts.EnvImageBuildersNamespace])
	if namespace == "" {
		namespace = consts.DefaultKubeNamespaceImageBuilders
	}

	return
}

func GetBentoDeploymentNamespaces(ctx context.Context, cliset *kubernetes.Clientset) (namespaces []string, err error) {
	namespaces_ := os.Getenv(consts.EnvBentoDeploymentNamespaces)
	if namespaces_ != "" {
		namespaces = strings.Split(namespaces_, ",")
		return
	}

	yataiSystemNamespace := GetYataiSystemNamespaceFromEnv()
	yataiDeploymentSharedEnvSecretName := consts.KubeSecretNameYataiDeploymentSharedEnv

	secret, err := cliset.CoreV1().Secrets(yataiSystemNamespace).Get(ctx, yataiDeploymentSharedEnvSecretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = errors.Wrapf(err, "secret %s not found in namespace %s", yataiDeploymentSharedEnvSecretName, yataiSystemNamespace)
		}
		return
	}

	namespaces_ = string(secret.Data[consts.EnvBentoDeploymentNamespaces])
	if namespaces_ == "" {
		namespaces = []string{consts.DefaultKubeNamespaceBentoDeployment}
	} else {
		namespaces = strings.Split(namespaces_, ",")
	}

	return
}

type S3Config struct {
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	Region     string `yaml:"region"`
	BucketName string `yaml:"bucket_name"`
	Secure     bool   `yaml:"secure"`
}

func GetS3Config(ctx context.Context) (conf *S3Config, err error) {
	conf = &S3Config{}
	conf.Endpoint = os.Getenv(consts.EnvS3Endpoint)
	conf.AccessKey = os.Getenv(consts.EnvS3AccessKey)
	conf.SecretKey = os.Getenv(consts.EnvS3SecretKey)
	conf.Region = os.Getenv(consts.EnvS3Region)
	conf.BucketName = os.Getenv(consts.EnvS3BucketName)
	conf.Secure = os.Getenv(consts.EnvS3Secure) == "true"

	if conf.Endpoint == "" {
		err = errors.Wrapf(consts.ErrNotFound, "the environment variable %s is not set", consts.EnvS3Endpoint)
	}

	return
}

type DockerRegistryConfig struct {
	BentoRepositoryName string `yaml:"bento_repository_name"`
	ModelRepositoryName string `yaml:"model_repository_name"`
	Server              string `yaml:"server"`
	InClusterServer     string `yaml:"in_cluster_server"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	Secure              bool   `yaml:"secure"`
}

func GetDockerRegistryConfig(ctx context.Context, cliset *kubernetes.Clientset) (conf *DockerRegistryConfig, err error) {
	conf = &DockerRegistryConfig{}
	conf.BentoRepositoryName = os.Getenv(consts.EnvDockerRegistryBentoRepositoryName)
	conf.ModelRepositoryName = os.Getenv(consts.EnvDockerRegistryModelRepositoryName)
	conf.Server = os.Getenv(consts.EnvDockerRegistryServer)
	conf.InClusterServer = os.Getenv(consts.EnvDockerRegistryInClusterServer)
	conf.Username = os.Getenv(consts.EnvDockerRegistryUsername)
	conf.Password = os.Getenv(consts.EnvDockerRegistryPassword)
	conf.Secure = os.Getenv(consts.EnvDockerRegistrySecure) == "true"

	if conf.Server == "" {
		yataiSystemNamespace := GetYataiSystemNamespaceFromEnv()
		yataiImageBuilderSharedEnvSecretName := consts.KubeSecretNameYataiImageBuilderSharedEnv

		var secret *corev1.Secret

		secret, err = cliset.CoreV1().Secrets(yataiSystemNamespace).Get(ctx, yataiImageBuilderSharedEnvSecretName, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				err = errors.Wrapf(err, "secret %s not found in namespace %s", yataiImageBuilderSharedEnvSecretName, yataiSystemNamespace)
			}
			return
		}

		conf.BentoRepositoryName = string(secret.Data[consts.EnvDockerRegistryBentoRepositoryName])
		conf.ModelRepositoryName = string(secret.Data[consts.EnvDockerRegistryModelRepositoryName])
		conf.Server = string(secret.Data[consts.EnvDockerRegistryServer])
		conf.InClusterServer = string(secret.Data[consts.EnvDockerRegistryInClusterServer])
		conf.Username = string(secret.Data[consts.EnvDockerRegistryUsername])
		conf.Password = string(secret.Data[consts.EnvDockerRegistryPassword])
		conf.Secure = string(secret.Data[consts.EnvDockerRegistrySecure]) == "true"

	}

	if conf.Server == "" {
		err = errors.Wrapf(consts.ErrNotFound, "the environment variable %s is not set", consts.EnvDockerRegistryServer)
	}

	return
}

type YataiConfig struct {
	Endpoint    string `yaml:"endpoint"`
	ClusterName string `yaml:"cluster_name"`
	ApiToken    string `yaml:"api_token"`
}

func GetYataiConfig(ctx context.Context, cliset *kubernetes.Clientset, yataiComponentName string, ignoreEnv bool) (conf *YataiConfig, err error) {
	conf = &YataiConfig{}
	if !ignoreEnv {
		conf.Endpoint = os.Getenv(consts.EnvYataiEndpoint)
		conf.ClusterName = os.Getenv(consts.EnvYataiClusterName)
		conf.ApiToken = os.Getenv(consts.EnvYataiApiToken)
	}

	yataiSystemNamespace := GetYataiSystemNamespaceFromEnv()

	if conf.Endpoint == "" {
		var secret *corev1.Secret
		secret, err = cliset.CoreV1().Secrets(yataiSystemNamespace).Get(ctx, consts.KubeSecretNameYataiCommonEnv, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				err = errors.Wrapf(err, "secret %s not found in namespace %s", consts.KubeSecretNameYataiCommonEnv, yataiSystemNamespace)
			}
			return
		}
		conf.Endpoint = string(secret.Data[consts.EnvYataiEndpoint])
		conf.ClusterName = string(secret.Data[consts.EnvYataiClusterName])
	}

	if conf.ApiToken == "" {
		var secret *corev1.Secret
		var secretName string
		var secretNamespace string
		if yataiComponentName == consts.YataiImageBuilderComponentName {
			secretName = consts.KubeSecretNameYataiImageBuilderEnv
			secretNamespace, err = GetYataiImageBuilderNamespace(ctx, cliset)
			if err != nil {
				err = errors.Wrapf(err, "failed to get namespace for %s", yataiComponentName)
				return
			}
		} else if yataiComponentName == consts.YataiDeploymentComponentName {
			secretName = consts.KubeSecretNameYataiDeploymentEnv
			secretNamespace, err = GetYataiDeploymentNamespace(ctx, cliset)
			if err != nil {
				err = errors.Wrapf(err, "failed to get namespace for %s", yataiComponentName)
				return
			}
		} else {
			err = errors.Errorf("invalid yatai component name %s", yataiComponentName)
			return
		}
		secret, err = cliset.CoreV1().Secrets(secretNamespace).Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				err = errors.Errorf("the secret %s in namespace %s does not exist", secretName, secretNamespace)
			} else {
				err = errors.Wrapf(err, "failed to get secret %s in namespace %s", secretName, secretNamespace)
			}
			return
		}
		conf.ApiToken = string(secret.Data[consts.EnvYataiApiToken])
	}

	return
}

// if key found in environ return value else return fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type InternalImages struct {
	Curl               string `yaml:"curl"`
	Kaniko             string `yaml:"kaniko"`
	MetricsTransformer string `yaml:"metrics-transformer"`
}

func GetInternalImages() (conf *InternalImages) {
	conf = &InternalImages{}
	conf.Curl = getEnv(consts.EnvInternalImagesCurl, consts.InternalImagesCurlDefault)
	conf.Kaniko = getEnv(consts.EnvInternalImagesKaniko, consts.InternalImagesKanikoDefault)
	conf.MetricsTransformer = getEnv(consts.EnvInternalImagesMetricsTransformer, consts.InternalImagesMetricsTransformerDefault)

	return
}
