package manifests

import (
	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/ViaQ/loki-operator/internal/manifests/internal"
)

// UseCertificateSigningService is a flag to allow certificate signing for the http services.
// This flag is set to true when the operator is built for the OpenShift platform.
var UseCertificateSigningService bool = false

// Options is a set of configuration values to use when building manifests such as resource sizes, etc.
// Most of this should be provided - either directly or indirectly - by the user.
type Options struct {
	Name       string
	Namespace  string
	Image      string
	ConfigSHA1 string

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
