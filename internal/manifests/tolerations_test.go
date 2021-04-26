package manifests

import (
	"testing"

	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestTolerationsAreSetForEachComponent(t *testing.T) {
	expected := []corev1.Toleration{{
		Key:      "type",
		Operator: corev1.TolerationOpEqual,
		Value:    "storage",
		Effect:   corev1.TaintEffectNoSchedule,
	}}
	opts := Options{
		Name:      uuid.New().String(),
		Namespace: uuid.New().String(),
		Image:     uuid.New().String(),
		Stack: lokiv1beta1.LokiStackSpec{
			Template: lokiv1beta1.LokiTemplateSpec{
				Compactor: lokiv1beta1.LokiComponentSpec{
					Tolerations: expected,
				},
				Distributor: lokiv1beta1.LokiComponentSpec{
					Tolerations: expected,
				},
				Ingester: lokiv1beta1.LokiComponentSpec{
					Tolerations: expected,
				},
				Querier: lokiv1beta1.LokiComponentSpec{
					Tolerations: expected,
				},
				QueryFrontend: lokiv1beta1.LokiComponentSpec{
					Tolerations: expected,
				},
			},
		},
		ObjectStorage: ObjectStorage{},
	}


	t.Run("distributor", func(t *testing.T) {
		d := NewDistributorDeployment(opts)
		assert.Equal(t, expected, d.Spec.Template.Spec.Tolerations)
	})


	t.Run("query_frontend", func(t *testing.T) {
		q := NewQueryFrontendDeployment(opts)
		assert.Equal(t, expected, q.Spec.Template.Spec.Tolerations)
	})

	t.Run("querier", func(t *testing.T) {
		q := NewQuerierStatefulSet(opts)
		assert.Equal(t, expected, q.Spec.Template.Spec.Tolerations)
	})

	t.Run("ingester", func(t *testing.T) {
		i := NewIngesterStatefulSet(opts)
		assert.Equal(t, expected, i.Spec.Template.Spec.Tolerations)
	})

	t.Run("compactor", func(t *testing.T) {
		c := NewCompactorStatefulSet(opts)
		assert.Equal(t, expected, c.Spec.Template.Spec.Tolerations)
	})
}

