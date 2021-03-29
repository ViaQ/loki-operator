package manifests_test

import (
	"testing"

	"github.com/ViaQ/loki-operator/internal/manifests"
	"github.com/stretchr/testify/require"
)

func TestNewIngesterStatefulSet_SelectorMatchesLabels(t *testing.T) {
	// You must set the .spec.selector field of a StatefulSet to match the labels of
	// its .spec.template.metadata.labels. Prior to Kubernetes 1.8, the
	// .spec.selector field was defaulted when omitted. In 1.8 and later versions,
	// failing to specify a matching Pod Selector will result in a validation error
	// during StatefulSet creation.
	// See https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#pod-selector
	ss, err := manifests.NewIngesterStatefulSet(manifests.Options{
		Name:      "abcd",
		Namespace: "efgh",
		Ingester: manifests.Ingester{
			Storage: manifests.Storage{
				ClassName:     "standard",
				SizeRequested: "1Gi",
			},
		},
	})
	require.NoError(t, err)
	l := ss.Spec.Template.GetObjectMeta().GetLabels()
	for key, value := range ss.Spec.Selector.MatchLabels {
		require.Contains(t, l, key)
		require.Equal(t, l[key], value)
	}
}
