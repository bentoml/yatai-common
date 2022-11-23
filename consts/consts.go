package consts

const (
	DefaultETCDTimeoutSeconds              = 5
	DefaultETCDDialKeepaliveTimeSeconds    = 30
	DefaultETCDDialKeepaliveTimeoutSeconds = 10

	HPADefaultMaxReplicas = 10

	HPACPUDefaultAverageUtilization = 80

	YataiDebugImg             = "yatai.ai/yatai-infras/debug"
	YataiKubectlNamespace     = "default"
	YataiKubectlContainerName = "main"
	YataiKubectlImage         = "yatai.ai/yatai-infras/k8s"

	TracingContextKey = "tracing-context"
	// nolint: gosec
	YataiApiTokenHeaderName = "X-YATAI-API-TOKEN"

	YataiOrganizationHeaderName = "X-Yatai-Organization"

	BentoServicePort          = 3000
	BentoContainerDefaultPort = 3000
	BentoServicePortName      = "http"
	BentoContainerPortName    = "http"

	NoneStr = "None"

	AmazonS3Endpoint = "s3.amazonaws.com"

	YataiImageBuilderComponentName = "yatai-image-builder"
	YataiDeploymentComponentName   = "yatai-deployment"

	// nolint: gosec
	YataiK8sBotApiTokenName = "yatai-k8s-bot"

	YataiBentoDeploymentComponentApiServer = "api-server"
	YataiBentoDeploymentComponentRunner    = "runner"

	InternalImagesBentoDownloaderDefault    = "quay.io/bentoml/bento-downloader:0.0.1"
	InternalImagesCurlDefault               = "quay.io/bentoml/curl:0.0.1"
	InternalImagesKanikoDefault             = "quay.io/bentoml/kaniko:1.9.1"
	InternalImagesMetricsTransformerDefault = "quay.io/bentoml/yatai-bento-metrics-transformer:0.0.2"
)
