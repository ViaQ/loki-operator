package manifests

import (
	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/ViaQ/loki-operator/internal/manifests/internal"
	"github.com/ViaQ/loki-operator/internal/manifests/internal/config"
)

// Options is a set of configuration values to use when building manifests such as resource sizes, etc.
// Most of this should be provided - either directly or indirectly - by the user.
type Options struct {
	Name      string
	Namespace string
	Image     string

	Config Config

	Stack                lokiv1beta1.LokiStackSpec
	ResourceRequirements internal.ResourceRequirements

	ObjectStorage ObjectStorage
}

// Config for config map contents.
type Config struct {
	Config        []byte
	RuntimeConfig []byte
	CompareResult config.CompareResult
}

// ObjectStorage for storage config.
type ObjectStorage struct {
	Endpoint        string
	Region          string
	Buckets         string
	AccessKeyID     string
	AccessKeySecret string
}
