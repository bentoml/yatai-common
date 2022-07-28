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

func GetS3Config(ctx context.Context) (conf *S3Config, err error) {
	conf = &S3Config{}
	conf.Endpoint = os.Getenv(consts.EnvS3Endpoint)
	conf.AccessKey = os.Getenv(consts.EnvS3AccessKey)
	conf.SecretKey = os.Getenv(consts.EnvS3SecretKey)
	conf.Region = os.Getenv(consts.EnvS3Region)
	conf.BucketName = os.Getenv(consts.EnvS3BucketName)
	conf.Secure = os.Getenv(consts.EnvS3Secure) == "true"

	if conf.Endpoint == "" {
		err = errors.Errorf("the environment variable %s is not set", consts.EnvS3Endpoint)
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

func GetDockerRegistryConfig(ctx context.Context) (conf *DockerRegistryConfig, err error) {
	conf = &DockerRegistryConfig{}
	conf.BentoRepositoryName = os.Getenv(consts.EnvDockerRegistryBentoRepositoryName)
	conf.ModelRepositoryName = os.Getenv(consts.EnvDockerRegistryModelRepositoryName)
	conf.Server = os.Getenv(consts.EnvDockerRegistryServer)
	conf.InClusterServer = os.Getenv(consts.EnvDockerRegistryInClusterServer)
	conf.Username = os.Getenv(consts.EnvDockerRegistryUsername)
	conf.Password = os.Getenv(consts.EnvDockerRegistryPassword)
	conf.Secure = os.Getenv(consts.EnvDockerRegistrySecure) == "true"

	if conf.Server == "" {
		err = errors.Errorf("the environment variable %s is not set", consts.EnvDockerRegistryServer)
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

func GetDockerImageBuilderConfig() (conf *DockerImageBuilderConfig) {
	conf = &DockerImageBuilderConfig{}
	conf.Privileged = os.Getenv(consts.EnvDockerImageBuilderPrivileged) == "true"

	return
}
