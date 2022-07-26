package consts

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	KubeIngressClassName = "yatai-ingress"

	KubeLabelYataiSelector        = "yatai.ai/selector"
	KubeLabelYataiBentoRepository = "yatai.ai/bento-repository"
	KubeLabelYataiBento           = "yatai.ai/bento"
	KubeLabelYataiModelRepository = "yatai.ai/model-repository"
	KubeLabelYataiModel           = "yatai.ai/model"

	KubeHPAQPSMetric = "http_request"
	KubeHPAGPUMetric = "container_accelerator_duty_cycle"

	KubeNamespaceYataiBentoImageBuilder = "yatai-builders"
	KubeNamespaceYataiModelImageBuilder = "yatai-builders"
	KubeNamespaceYataiDeployment        = "yatai"
	KubeNamespaceYataiOperators         = "yatai-operators"
	KubeNamespaceYataiComponents        = "yatai-components"

	KubeNamespaceYataiSystem              = "yatai-system"
	KubeNamespaceYataiDeploymentComponent = "yatai-deployment"

	KubeLabelMcdInfraCli               = "mcd-infra-cli"
	KubeLabelMcdKubectl                = "mcd-kubectl"
	KubeLabelMcdUser                   = "mcd-user"
	KubeLabelMcdAppPool                = "mcd-app-pool"
	KubeLabelYataiDeployment           = "yatai.ai/deployment"
	KubeLabelYataiDeploymentId         = "yatai.ai/deployment-id"
	KubeLabelYataiDeploymentTargetType = "yatai.ai/deployment-target-type"
	KubeLabelYataiBentoRunner          = "yatai.ai/bento-runner"
	KubeLabelYataiIsBentoApiServer     = "yatai.ai/is-bento-api-server"
	KubeLabelBentoRepository           = "yatai.ai/bento-repository"
	KubeLabelBentoVersion              = "yatai.ai/bento-version"
	KubeLabelCreator                   = "creator"
	// nolint: gosec
	KubeLabelYataiDeployToken = "yatai.ai/deploy-token"

	KubeLabelMcdAppCompType = "mcd-app-comp-type"
	KubeLabelMcdAppCompName = "mcd-app-comp-name"

	KubeLabelYataiOwnerReference = "yatai.ai/owner-reference"

	KubeLabelGPUAccelerator = "gpu-accelerator"

	KubeLabelHostName = "kubernetes.io/hostname"
	KubeLabelArch     = "kubernetes.io/arch"

	KubeLabelMcdNodePool       = "mcd.io/node-pool"
	KubeLabelAlibabaEdgeWorker = "alibabacloud.com/is-edge-worker"
	KubeLabelMcdEdgeWorker     = "mcd.io/is-edge-worker"
	KubeLabelFalse             = "false"
	KubeLabelTrue              = "true"

	KubeLabelManagedBy    = "app.kubernetes.io/managed-by"
	KubeLabelHelmHeritage = "heritage"
	KubeLabelHelmRelease  = "release"

	KubeAnnotationBentoRepository   = "yatai.ai/bento-repository"
	KubeAnnotationBentoVersion      = "yatai.ai/bento-version"
	KubeAnnotationYataiDeploymentId = "yatai.ai/deployment-id"
	KubeAnnotationHelmReleaseName   = "meta.helm.sh/release-name"

	KubeAnnotationPrometheusScrape = "prometheus.io/scrape"
	KubeAnnotationPrometheusPort   = "prometheus.io/port"
	KubeAnnotationPrometheusPath   = "prometheus.io/path"

	KubeAnnotationARMSAutoEnable = "armsPilotAutoEnable"
	KubeAnnotationARMSAppName    = "armsPilotCreateAppName"

	KubeCreator = "yatai"

	KubeVolumeNamePermdir                            = "permdir"
	KubeVolumeNameFastPermdir                        = "fast-permdir"
	KubeVolumeNameHostTimezone                       = "host-timezone"
	KubeVolumeNameMcdTracingAgentDir                 = "mcd-tracing"
	KubeVolumeNameMcdJmxAgentDir                     = "mcd-jmx"
	KubeVolumeMountPathPermdir                       = "/permdir"
	KubeVolumeMountPathFastPermdir                   = "/fast_permdir"
	KubeVolumeNameDockerSock                         = "mcd-docker-sock"
	KubeVolumeMountPathDockerSock                    = "/var/run/docker.sock"
	KubeVolumeNameDockerGraphStorage                 = "mcd-docker-graph-storage"
	KubeVolumeMountPathDockerGraphStorage            = "/var/lib/docker"
	KubeVolumeNameVarRun                             = "mcd-var-run"
	KubeVolumeMountPathVarRun                        = "/var/run"
	KubePersistentVolumeClaimNamePermdir             = "mcd-app-permdir"
	KubePersistentVolumeClaimNameFastPermdir         = "mcd-app-fast-permdir"
	KubePersistentVolumeClaimPermdirStorageClass     = "mcd-nfs"
	KubePersistentVolumeClaimFastPermdirStorageClass = "mcd-fast-nfs"
	KubeAliCouldStorageClassProvisioner              = "nasplugin.csi.alibabacloud.com"

	KubeIngressCanaryHeader      = "mcd-canary"
	KubeIngressCanaryHeaderValue = "always"

	KubeNameMcdDns = "mcd-dns"

	KubeStorageClassNameMcd       = "mcd"
	KubeStorageClassNameLocalPath = "local-path"

	KubeResourceGPUNvidia = "nvidia.com/gpu"

	KubeEventResourceKindPod        = "Pod"
	KubeEventResourceKindHPA        = "HorizontalPodAutoscaler"
	KubeEventResourceKindReplicaSet = "ReplicaSet"

	KubeTaintKeyDedicatedNodeGroup = "mcd.io/dedicated-node-group"
	KubeLabelDedicatedNodeGroup    = "mcd.io/dedicated-node-group"

	KubeLabelMcdESEnable  = "mcd-es-enable"
	KubeLabelMcdESSaveDay = "mcd-es-save-day"

	KubeImageCSIDriver          = "image.csi.k8s.io"
	KubeImageCSIDriverWarmMetal = "csi-image.warm-metal.tech"

	KubeDefaultMcdResourceQuotaName = "mcd"

	KubeLabelNodeResourceResizeCPU    = "mcd.io/resize-node-cpu"
	KubeLabelNodeResourceResizePods   = "mcd.io/resize-node-pods"
	KubeLabelNodeResourceResizeMemory = "mcd.io/resize-node-memory"

	KubeConfigMapNameNetworkConfig = "network"

	KubeConfigMapKeyNetworkConfigDomainSuffix = "domain-suffix"
	KubeConfigMapKeyNetworkConfigIngressClass = "ingress-class"

	KubeConfigMapNameS3Config = "s3"

	KubeConfigMapKeyS3ConfigEndpoint            = "endpoint"
	KubeConfigMapKeyS3ConfigAccessKeySecretName = "access-key-secret-name"
	KubeConfigMapKeyS3ConfigAccessKeySecretKey  = "access-key-secret-key"
	KubeConfigMapKeyS3ConfigSecretKeySecretName = "secret-key-secret-name"
	KubeConfigMapKeyS3ConfigSecretKeySecretKey  = "secret-key-secret-key"
	KubeConfigMapKeyS3ConfigRegion              = "region"
	KubeConfigMapKeyS3ConfigBucketName          = "bucket-name"
	KubeConfigMapKeyS3ConfigSecure              = "secure"

	KubeConfigMapNameDockerRegistryConfig = "docker-registry"

	KubeConfigMapKeyDockerRegistryConfigBentoRepositoryName = "bento-repository-name"
	KubeConfigMapKeyDockerRegistryConfigModelRepositoryName = "model-repository-name"
	KubeConfigMapKeyDockerRegistryConfigServer              = "server"
	KubeConfigMapKeyDockerRegistryConfigUsername            = "username"
	KubeConfigMapKeyDockerRegistryConfigPasswordSecretName  = "password-secret-name"
	KubeConfigMapKeyDockerRegistryConfigPasswordSecretKey   = "password-secret-key"
	KubeConfigMapKeyDockerRegistryConfigSecure              = "secure"

	KubeConfigMapNameDockerImageBuilderConfig = "docker-image-builder"

	KubeConfigMapKeyDockerImageBuilderConfigPrivileged = "privileged"

	KubeConfigMapNameYataiConfig = "yatai"

	KubeConfigMapKeyYataiConfigEndpoint           = "endpoint"
	KubeConfigMapKeyYataiConfigClusterName        = "cluster-name"
	KubeConfigMapKeyYataiConfigApiTokenSecretName = "api-token-secret-name"
	KubeConfigMapKeyYataiConfigApiTokenSecretKey  = "api-token-secret-key"

	// nolint: gosec
	KubeSecretNameRegcred = "yatai-regcred"
)

var KubeListEverything = metav1.ListOptions{
	LabelSelector: labels.Everything().String(),
	FieldSelector: fields.Everything().String(),
}
