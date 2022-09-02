package consts

const (
	DefaultNewsURL                         = "https://raw.githubusercontent.com/bentoml/yatai-homepage-news/main/news.json"
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

	YataiApiTokenPrefixYataiDeploymentOperator = "yatai-deployment-operator"

	// nolint: gosec
	YataiK8sBotApiTokenName = "yatai-k8s-bot"
)
