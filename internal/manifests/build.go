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

	cm, sha1C, mapErr := LokiConfigMap(opt)
	if mapErr != nil {
		return nil, mapErr
	}
	opt.ConfigSHA1 = sha1C

	res = append(res, cm)
	res = append(res, BuildLokiGossipRingService(opt.Name))

	distributorDeployment, distributorServices := BuildDistributor(opt)
	res = append(res, distributorServices...)

	ingesterStatefulSet, ingesterServices := BuildIngester(opt)
	res = append(res, ingesterServices...)

	querierStatefulSet, querierServices := BuildQuerier(opt)
	res = append(res, querierServices...)

	compactorStatefulSet, compactorServices := BuildCompactor(opt)
	res = append(res, compactorServices...)

	queryFrontendDeployment, queryFrontendServices := BuildQueryFrontend(opt)
	res = append(res, queryFrontendServices...)

	if opt.EnableTLSServiceMonitorConfig {
		if err := configureDistributorServiceMonitorPKI(distributorDeployment, opt.Name); err != nil {
			return nil, err
		}

		if err := configureIngesterServiceMonitorPKI(ingesterStatefulSet, opt.Name); err != nil {
			return nil, err
		}

		if err := configureQuerierServiceMonitorPKI(querierStatefulSet, opt.Name); err != nil {
			return nil, err
		}

		if err := configureCompactorServiceMonitorPKI(compactorStatefulSet, opt.Name); err != nil {
			return nil, err
		}

		if err := configureQueryFrontendServiceMonitorPKI(queryFrontendDeployment, opt.Name); err != nil {
			return nil, err
		}
	}

	res = append(res, distributorDeployment)
	res = append(res, ingesterStatefulSet)
	res = append(res, querierStatefulSet)
	res = append(res, compactorStatefulSet)
	res = append(res, queryFrontendDeployment)

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
