package config_test

import (
	"testing"

	"github.com/ViaQ/loki-operator/internal/manifests/internal/config"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestPathExporter_WhenEqual_ReportNothing(t *testing.T) {
	var (
		r    config.PathReporter
		a, b config.Loki
	)
	ok := cmp.Equal(a, b, cmp.Reporter(&r))
	require.True(t, ok)
	require.Empty(t, r.String())
}

func TestPathExporter_WhenNotEqual_ReportDiffPathOnly(t *testing.T) {
	var (
		r config.PathReporter
		a = config.Loki{
			Ingester: config.Ingester{
				Lifecycler: config.Lifecycler{
					Ring: config.Ring{
						ReplicationFactor: 1,
					},
				},
			},
		}
		b = config.Loki{
			Ingester: config.Ingester{
				Lifecycler: config.Lifecycler{
					Ring: config.Ring{
						ReplicationFactor: 2,
					},
				},
			},
		}
	)
	ok := cmp.Equal(a, b, cmp.Reporter(&r))
	require.False(t, ok)
	require.Equal(t, "Ingester.Lifecycler.Ring.ReplicationFactor", r.String())
}

func TestPathExporter_WhenNotEqual_ReturnRootsOnly(t *testing.T) {
	var (
		r config.PathReporter
		a = config.Loki{
			Ingester: config.Ingester{
				Lifecycler: config.Lifecycler{
					Ring: config.Ring{
						ReplicationFactor: 1,
					},
				},
			},
			StorageConfig: config.StorageConfig{
				AWS: config.AWS{
					Region: "us-nowhere",
				},
			},
		}
		b = config.Loki{
			Ingester: config.Ingester{
				Lifecycler: config.Lifecycler{
					Ring: config.Ring{
						ReplicationFactor: 2,
					},
				},
			},
			StorageConfig: config.StorageConfig{
				AWS: config.AWS{
					Region: "us-somewhere",
				},
			},
		}
	)
	ok := cmp.Equal(a, b, cmp.Reporter(&r))
	require.False(t, ok)
	expected := []string{"Ingester", "StorageConfig"}
	require.ElementsMatch(t, expected, r.Roots())
}
