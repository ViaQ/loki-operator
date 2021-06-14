package manifests

import (
	"github.com/ViaQ/logerr/kverrors"
	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/ViaQ/loki-operator/internal/manifests/internal"

	"github.com/imdario/mergo"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildAll builds all manifests required to run a Loki Stack
func BuildAll(opt Options) ([]client.Object, error) {
	res := make([]client.Object, 0)

	cm, sha1C, err := LokiConfigMap(opt)
	if err != nil {
		return nil, err
	}
	opt.ConfigSHA1 = sha1C

	res = append(res, cm)
	res = append(res, BuildLokiGossipRingService(opt.Name))

	objects, buildErr := BuildDistributor(opt)
	if buildErr != nil {
		return nil, buildErr
	}
	res = append(res, objects...)
	res = append(res, BuildIngester(opt)...)
	res = append(res, BuildQuerier(opt)...)
	res = append(res, BuildCompactor(opt)...)
	res = append(res, BuildQueryFrontend(opt)...)

	if opt.EnableServiceMonitors {
		res = append(res, BuildServiceMonitors(opt)...)
	}

	return res, nil
}

// DefaultLokiStackSpec returns the default configuration for a LokiStack of
// the specified size
func DefaultLokiStackSpec(size lokiv1beta1.LokiStackSizeType) *lokiv1beta1.LokiStackSpec {
	defaults := internal.StackSizeTable[size]
	return (&defaults).DeepCopy()
}

// ApplyDefaultSettings manipulates the options to conform to
// build specifications
func ApplyDefaultSettings(opt *Options) error {
	spec := DefaultLokiStackSpec(opt.Stack.Size)

	if err := mergo.Merge(spec, opt.Stack, mergo.WithOverride); err != nil {
		return kverrors.Wrap(err, "failed merging stack user options", "name", opt.Name)
	}

	strictOverrides := lokiv1beta1.LokiStackSpec{
		Template: &lokiv1beta1.LokiTemplateSpec{
			Compactor: &lokiv1beta1.LokiComponentSpec{
				// Compactor is a singelton application.
				// Only one replica allowed!!!
				Replicas: 1,
			},
		},
	}

	if err := mergo.Merge(spec, strictOverrides, mergo.WithOverride); err != nil {
		return kverrors.Wrap(err, "failed to merge strict defaults")
	}

	opt.ResourceRequirements = internal.ResourceRequirementsTable[opt.Stack.Size]
	opt.Stack = *spec

	return nil
}
