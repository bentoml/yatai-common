package consts

const (
	EnvYataiEndpoint    = "YATAI_ENDPOINT"
	EnvYataiClusterName = "YATAI_CLUSTER_NAME"
	// nolint: gosec
	EnvYataiApiToken = "YATAI_API_TOKEN"

	EnvBentoServicePort = "PORT"

	// tracking envars
	EnvYataiVersion       = "YATAI_T_VERSION"
	EnvYataiOrgUID        = "YATAI_T_ORG_UID"
	EnvYataiDeploymentUID = "YATAI_T_DEPLOYMENT_UID"
	EnvYataiClusterUID    = "YATAI_T_CLUSTER_UID"

	EnvYataiBentoDeploymentName      = "YATAI_BENTO_DEPLOYMENT_NAME"
	EnvYataiBentoDeploymentNamespace = "YATAI_BENTO_DEPLOYMENT_NAMESPACE"

	EnvS3Endpoint   = "S3_ENDPOINT"
	EnvS3Region     = "S3_REGION"
	EnvS3BucketName = "S3_BUCKET_NAME"
	EnvS3AccessKey  = "S3_ACCESS_KEY"
	// nolint:gosec
	EnvS3SecretKey = "S3_SECRET_KEY"
	EnvS3Secure    = "S3_SECURE"

	EnvDockerRegistryServer          = "DOCKER_REGISTRY_SERVER"
	EnvDockerRegistryInClusterServer = "DOCKER_REGISTRY_IN_CLUSTER_SERVER"
	EnvDockerRegistryUsername        = "DOCKER_REGISTRY_USERNAME"
	// nolint:gosec
	EnvDockerRegistryPassword            = "DOCKER_REGISTRY_PASSWORD"
	EnvDockerRegistrySecure              = "DOCKER_REGISTRY_SECURE"
	EnvDockerRegistryBentoRepositoryName = "DOCKER_REGISTRY_BENTO_REPOSITORY_NAME"
	EnvDockerRegistryModelRepositoryName = "DOCKER_REGISTRY_MODEL_REPOSITORY_NAME"

	InternalImagesCurl               = "INTERNAL_IMAGES_CURL"
	InternalImagesKaniko             = "INTERNAL_IMAGES_KANIKO"
	InternalImagesMetricsTransformer = "INTERNAL_IMAGES_METRICS_TRANSFORMER"

	InternalImagesCurlDefault               = "quay.io/bentoml/curl"
	InternalImagesKanikoDefault             = "quay.io/bentoml/kaniko"
	InternalImagesMetricsTransformerDefault = "quay.io/bentoml/yatai-bento-metrics-transformer"
)
