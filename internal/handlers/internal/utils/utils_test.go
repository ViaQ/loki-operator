package utils

import (
	"testing"

	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAppendDefaultLabel_WhenLabelsNil_AppendNewLabels(t *testing.T) {
	clusterName := "loki-stack"

	ls := lokiv1beta1.LokiStack{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-stack",
			Namespace: "some-ns",
			Labels:    nil,
		},
	}

	defaultLabel := AppendDefaultLabel(clusterName, ls.Labels)
	require.Equal(t, clusterName, defaultLabel["cluster-name"])
}

func TestAppendDefaultLabel_WhenLabelsAlreadyPresent_SimplyReturnLabels(t *testing.T) {
	clusterName := "loki-stack"

	ls := lokiv1beta1.LokiStack{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-stack",
			Namespace: "some-ns",
			Labels: map[string]string{
				"cluster-name": clusterName,
			},
		},
	}

	defaultLabel := AppendDefaultLabel(clusterName, ls.Labels)
	require.Equal(t, clusterName, defaultLabel["cluster-name"])
	// To verify that already present label isn't appended once more
	require.Equal(t, 1, len(defaultLabel))
}
