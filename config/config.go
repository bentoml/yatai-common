package config

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/bentoml/yatai-common/consts"
)

func getValueFromSecret(ctx context.Context, cliset *kubernetes.Clientset, configMap *corev1.ConfigMap, configMapKeySecretName, configMapKeySecretKey string) (value string, err error) {
	secretName := string(configMap.Data[configMapKeySecretName])
	if secretName == "" {
		err = errors.Errorf("the config map %s in namespace %s does not contain key %s", configMap.Name, configMap.Namespace, configMapKeySecretName)
		return
	}
	secretKey := string(configMap.Data[configMapKeySecretKey])
	if secretKey == "" {
		err = errors.Errorf("the config map %s in namespace %s does not contain key %s", configMap.Name, configMap.Namespace, configMapKeySecretKey)
		return
	}
	secret, err := cliset.CoreV1().Secrets(configMap.Namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		err = errors.Wrapf(err, "failed to get secret %s in namespace %s", secretName, configMap.Namespace)
		return
	}
	valueRaw := string(secret.Data[secretKey])
	var value_ []byte
	value_, err = base64.StdEncoding.DecodeString(valueRaw)
	if err != nil {
		err = errors.Wrapf(err, "failed to decode the field %s from secret %s in namespace %s", secretKey, secretName, configMap.Namespace)
		return
	}
	value = string(value_)
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

func GetS3Config(ctx context.Context, cliset *kubernetes.Clientset) (conf *S3Config, err error) {
	configMapName := consts.KubeConfigMapNameS3Config
	namespace := consts.KubeNamespaceYataiSystem

	conf = &S3Config{}
	conf.Endpoint = os.Getenv(consts.EnvS3Endpoint)
	conf.AccessKey = os.Getenv(consts.EnvS3AccessKey)
	conf.SecretKey = os.Getenv(consts.EnvS3SecretKey)
	conf.Region = os.Getenv(consts.EnvS3Region)
	conf.BucketName = os.Getenv(consts.EnvS3BucketName)
	secure, secureEnvExists := os.LookupEnv(consts.EnvS3Secure)
	if secureEnvExists {
		conf.Secure = secure == "true"
	}

	configMapCli := cliset.CoreV1().ConfigMaps(namespace)
	configMap, err := configMapCli.Get(ctx, configMapName, metav1.GetOptions{})
	isNotFound := k8serrors.IsNotFound(err)
	if err != nil && !isNotFound {
		err = errors.Wrapf(err, "failed to get config map %s in namespace %s", configMapName, namespace)
		return
	}

	if isNotFound {
		return
	}

	if conf.Endpoint == "" {
		conf.Endpoint = string(configMap.Data[consts.KubeConfigMapKeyS3ConfigEndpoint])
	}
	if conf.AccessKey == "" {
		conf.AccessKey, err = getValueFromSecret(ctx, cliset, configMap, consts.KubeConfigMapKeyS3ConfigAccessKeySecretName, consts.KubeConfigMapKeyS3ConfigAccessKeySecretKey)
		if err != nil {
			return
		}
	}
	if conf.SecretKey == "" {
		conf.SecretKey, err = getValueFromSecret(ctx, cliset, configMap, consts.KubeConfigMapKeyS3ConfigSecretKeySecretName, consts.KubeConfigMapKeyS3ConfigSecretKeySecretKey)
		if err != nil {
			return
		}
	}
	if conf.Region == "" {
		conf.Region = string(configMap.Data[consts.KubeConfigMapKeyS3ConfigRegion])
	}
	if conf.BucketName == "" {
		conf.BucketName = string(configMap.Data[consts.KubeConfigMapKeyS3ConfigBucketName])
	}
	if !secureEnvExists {
		conf.Secure = string(configMap.Data[consts.KubeConfigMapKeyS3ConfigSecure]) == "true"
	}

	return
}

type DockerRegistryConfig struct {
	BentoRepositoryName string `yaml:"bento_repository_name"`
	ModelRepositoryName string `yaml:"model_repository_name"`
	Server              string `yaml:"server"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	Secure              bool   `yaml:"secure"`
}

func GetDockerRegistryConfig(ctx context.Context, cliset *kubernetes.Clientset) (conf *DockerRegistryConfig, err error) {
	configMapName := consts.KubeConfigMapNameDockerRegistryConfig
	namespace := consts.KubeNamespaceYataiDeploymentComponent

	conf = &DockerRegistryConfig{}
	conf.BentoRepositoryName = os.Getenv(consts.EnvDockerRegistryBentoRepositoryName)
	conf.ModelRepositoryName = os.Getenv(consts.EnvDockerRegistryModelRepositoryName)
	conf.Server = os.Getenv(consts.EnvDockerRegistryServer)
	conf.Username = os.Getenv(consts.EnvDockerRegistryUsername)
	conf.Password = os.Getenv(consts.EnvDockerRegistryPassword)
	secure, secureEnvExists := os.LookupEnv(consts.EnvDockerRegistrySecure)
	if secureEnvExists {
		conf.Secure = secure == "true"
	}

	configMapCli := cliset.CoreV1().ConfigMaps(namespace)
	configMap, err := configMapCli.Get(ctx, configMapName, metav1.GetOptions{})
	isNotFound := k8serrors.IsNotFound(err)
	if err != nil && !isNotFound {
		err = errors.Wrapf(err, "failed to get config map %s in namespace %s", configMapName, namespace)
	}
	if isNotFound {
		return
	}

	if conf.BentoRepositoryName == "" {
		conf.BentoRepositoryName = string(configMap.Data[consts.KubeConfigMapKeyDockerRegistryConfigBentoRepositoryName])
	}
	if conf.ModelRepositoryName == "" {
		conf.ModelRepositoryName = string(configMap.Data[consts.KubeConfigMapKeyDockerRegistryConfigModelRepositoryName])
	}
	if conf.Server == "" {
		conf.Server = string(configMap.Data[consts.KubeConfigMapKeyDockerRegistryConfigServer])
	}
	if conf.Username == "" {
		conf.Username = string(configMap.Data[consts.KubeConfigMapKeyDockerRegistryConfigUsername])
	}
	if conf.Password == "" {
		conf.Password, err = getValueFromSecret(ctx, cliset, configMap, consts.KubeConfigMapKeyDockerRegistryConfigPasswordSecretName, consts.KubeConfigMapKeyDockerRegistryConfigPasswordSecretKey)
		if err != nil {
			return
		}
	}
	if !secureEnvExists {
		conf.Secure = string(configMap.Data[consts.KubeConfigMapKeyDockerRegistryConfigSecure]) == "true"
	}
	return
}

type YataiConfig struct {
	Endpoint    string `yaml:"endpoint"`
	ClusterName string `yaml:"cluster_name"`
	ApiToken    string `yaml:"api_token"`
}

func GetYataiConfig(ctx context.Context, cliset *kubernetes.Clientset, ignoreEnv bool) (conf *YataiConfig, err error) {
	configMapName := consts.KubeConfigMapNameYataiConfig
	namespace := consts.KubeNamespaceYataiDeploymentComponent

	conf = &YataiConfig{}
	if !ignoreEnv {
		conf.Endpoint = os.Getenv(consts.EnvYataiEndpoint)
		conf.ClusterName = os.Getenv(consts.EnvYataiClusterName)
		conf.ApiToken = os.Getenv(consts.EnvYataiApiToken)
	}

	configMapCli := cliset.CoreV1().ConfigMaps(namespace)
	configMap, err := configMapCli.Get(ctx, configMapName, metav1.GetOptions{})
	isNotFound := k8serrors.IsNotFound(err)
	if err != nil && !isNotFound {
		err = errors.Wrapf(err, "failed to get config map %s in namespace %s", configMapName, namespace)
		return
	}

	if isNotFound {
		return
	}

	if conf.Endpoint == "" {
		conf.Endpoint = string(configMap.Data[consts.KubeConfigMapKeyYataiConfigEndpoint])
	}
	if conf.ClusterName == "" {
		conf.ClusterName = string(configMap.Data[consts.KubeConfigMapKeyYataiConfigClusterName])
	}
	if conf.ApiToken == "" {
		conf.ApiToken, err = getValueFromSecret(ctx, cliset, configMap, consts.KubeConfigMapKeyYataiConfigApiTokenSecretName, consts.KubeConfigMapKeyYataiConfigApiTokenSecretKey)
		if err != nil {
			return
		}
	}

	return
}

type DockerImageBuilderConfig struct {
	Privileged bool `yaml:"privileged"`
}

func GetDockerImageBuilderConfig(ctx context.Context, cliset *kubernetes.Clientset) (conf *DockerImageBuilderConfig, err error) {
	configMapName := consts.KubeConfigMapNameDockerImageBuilderConfig
	namespace := consts.KubeNamespaceYataiDeploymentComponent

	conf = &DockerImageBuilderConfig{}
	privileged, privilegedEnvExists := os.LookupEnv(consts.EnvDockerImageBuilderPrivileged)
	if privilegedEnvExists {
		conf.Privileged = privileged == "true"
	}

	configMapCli := cliset.CoreV1().ConfigMaps(namespace)
	configMap, err := configMapCli.Get(ctx, configMapName, metav1.GetOptions{})
	isNotFound := k8serrors.IsNotFound(err)
	if err != nil && !isNotFound {
		err = errors.Wrapf(err, "failed to get config map %s in namespace %s", configMapName, namespace)
		return
	}

	if isNotFound {
		return
	}

	if !privilegedEnvExists {
		conf.Privileged = string(configMap.Data[consts.KubeConfigMapKeyDockerImageBuilderConfigPrivileged]) == "true"
	}

	return
}
