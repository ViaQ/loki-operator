package handlers

import (
	"context"
	"errors"
	"testing"

	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/ViaQ/loki-operator/internal/external/k8s/k8sfakes"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestCreateOrUpdateServiceMonitor_WhenGetReturnsNotFound_DoesNotError(t *testing.T) {
	k := &k8sfakes.FakeClient{}
	r := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "my-stack",
			Namespace: "some-ns",
		},
	}

	k.GetStub = func(ctx context.Context, name types.NamespacedName, object client.Object) error {
		return apierrors.NewNotFound(schema.GroupResource{}, "something wasn't found")
	}

	err := CreateOrUpdateServiceMonitor(context.TODO(), r, k)
	require.NoError(t, err)

	// make sure create was NOT called because the Get failed
	require.Zero(t, k.CreateCallCount())
}

func TestCreateOrUpdateServiceMonitor_WhenGetReturnsAnErrorOtherThanNotFound_ReturnsTheError(t *testing.T) {
	k := &k8sfakes.FakeClient{}
	r := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "my-stack",
			Namespace: "some-ns",
		},
	}

	badRequestErr := apierrors.NewBadRequest("doesn't belong here")
	k.GetStub = func(ctx context.Context, name types.NamespacedName, object client.Object) error {
		return badRequestErr
	}

	err := CreateOrUpdateServiceMonitor(context.TODO(), r, k)

	require.Equal(t, badRequestErr, errors.Unwrap(err))

	// make sure create was NOT called because the Get failed
	require.Zero(t, k.CreateCallCount())
}

func TestCreateOrUpdateServiceMonitor_WhenCreateReturnsNoError_ReturnNil(t *testing.T) {
	k := &k8sfakes.FakeClient{}
	r := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "my-stack",
			Namespace: "some-ns",
		},
	}

	ls := lokiv1beta1.LokiStack{
		TypeMeta: metav1.TypeMeta{
			Kind: "LokiStack",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "someStack",
			Namespace: "some-ns",
			UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
		},
	}

	svcMonitor := monitoringv1.ServiceMonitor{
		TypeMeta: metav1.TypeMeta{
			Kind:       monitoringv1.ServiceMonitorsKind,
			APIVersion: monitoringv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "monitor-someStack-cluster",
			Namespace: "some-ns",
			Labels: map[string]string{
				"cluster-name":   "someStack",
				"scrape-metrics": "enabled",
			},
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			JobLabel: labelJobName,
			Endpoints: []monitoringv1.Endpoint{
				{
					Port:            "someStack",
					Path:            "/metrics",
					Scheme:          "https",
					BearerTokenFile: bearerTokenFile,
					TLSConfig: &monitoringv1.TLSConfig{
						SafeTLSConfig: monitoringv1.SafeTLSConfig{
							ServerName: "loki-metrics.openshift-logging.svc",
						},
						CAFile: prometheusCAFile,
					},
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"cluster-name":   "someStack",
					"scrape-metrics": "enabled",
				},
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{"some-ns"},
			},
		},
	}

	// Create looks up the CR first, so we need to return our fake stack
	k.GetStub = func(_ context.Context, name types.NamespacedName, object client.Object) error {
		if r.Name == name.Name && r.Namespace == name.Namespace {
			k.SetClientObject(object, &ls)
		}
		if svcMonitor.Name == name.Name && svcMonitor.Namespace == name.Namespace {
			k.SetClientObject(object, &svcMonitor)
		}
		return nil
	}

	err := CreateOrUpdateServiceMonitor(context.TODO(), r, k)
	require.NoError(t, err)

	// make sure create was called
	require.NotZero(t, k.CreateCallCount())
}
