package config

import (
	"context"
	"os"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/bentoml/yatai-common/consts"
)

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
		err = errors.Wrapf(consts.ErrNotFound, "the environment variable %s is not set", consts.EnvDockerRegistryServer)
	}

	return
}

type YataiConfig struct {
	Endpoint    string `yaml:"endpoint"`
	ClusterName string `yaml:"cluster_name"`
	ApiToken    string `yaml:"api_token"`
}

func GetYataiConfig(ctx context.Context, cliset *kubernetes.Clientset, namespace string, ignoreEnv bool) (conf *YataiConfig, err error) {
	conf = &YataiConfig{}
	if !ignoreEnv {
		conf.Endpoint = os.Getenv(consts.EnvYataiEndpoint)
		conf.ClusterName = os.Getenv(consts.EnvYataiClusterName)
		conf.ApiToken = os.Getenv(consts.EnvYataiApiToken)
	}

	if conf.ApiToken == "" {
		secretName := "env"
		var secret *corev1.Secret
		secret, err = cliset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				err = errors.Errorf("the secret %s in namespace %s does not exist", secretName, namespace)
			} else {
				err = errors.Wrapf(err, "failed to get secret %s in namespace %s", secretName, namespace)
			}
			return
		}
		conf.ApiToken = string(secret.Data[consts.EnvYataiApiToken])
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
