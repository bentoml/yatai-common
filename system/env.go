package system

import (
	"fmt"
	"os"
)

const (
	// NamespaceEnvKey is the environment variable that specifies the system namespace.
	NamespaceEnvKey = "SYSTEM_NAMESPACE"
	// ResourceLabelEnvKey is the environment variable that specifies the system resource
	// label.
	ResourceLabelEnvKey = "SYSTEM_RESOURCE_LABEL"

	MagicDNSEnvKey  = "MAGIC_DNS"
	DefaultMagicDNS = "sslip.io"
)

// GetNamespace returns the name of the K8s namespace where our system components
// run.
func GetNamespace() string {
	if ns := os.Getenv(NamespaceEnvKey); ns != "" {
		return ns
	}

	panic(fmt.Sprintf(`The environment variable %q is not set

If this is a process running on Kubernetes, then it should be using the downward
API to initialize this variable via:

  env:
  - name: %s
    valueFrom:
      fieldRef:
        fieldPath: metadata.namespace
`, NamespaceEnvKey, NamespaceEnvKey))
}

// GetResourceLabel returns the label key identifying K8s objects our system
// components source their configuration from.
func GetResourceLabel() string {
	return os.Getenv(ResourceLabelEnvKey)
}

func GetMagicDNS() string {
	magicDNS := os.Getenv(MagicDNSEnvKey)
	if magicDNS == "" {
		magicDNS = DefaultMagicDNS
	}
	return magicDNS
}
