package manifests

import (
	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/ViaQ/loki-operator/internal/manifests/internal"
)

// Options is a set of configuration values to use when building manifests such as resource sizes, etc.
// Most of this should be provided - either directly or indirectly - by the user.
type Options struct {
	Name       string
	Namespace  string
	Image      string
	ConfigSHA1 string

	Flags FeatureFlags

	Stack                lokiv1beta1.LokiStackSpec
	ResourceRequirements internal.ComponentResources

	ObjectStorage ObjectStorage
}

// ObjectStorage for storage config.
type ObjectStorage struct {
	Endpoint        string
	Region          string
	Buckets         string
	AccessKeyID     string
	AccessKeySecret string
}

// FeatureFlags contains flags that activate various features
type FeatureFlags struct {
	EnableCertificateSigningService   bool
	EnableServiceMonitors             bool
	EnableTLSServiceMonitorConfig     bool
	EnableLokiStackGateway            bool
	EnableLokiStackGatewayTLSListener bool
}
