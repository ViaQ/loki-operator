package config_test

import (
	"testing"

	"github.com/ViaQ/loki-operator/internal/manifests/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCompare_WhenLimitsConfigNotEqual_ReturnAllWithCompactor(t *testing.T) {
	var (
		a = []byte(`
limits_config:
  ingestion_rate_mb: 100
`)
		b = []byte(`
limits_config:
  ingestion_rate_mb: 200
`)
		ah, _ = config.Sha1sum(a)
		bh, _ = config.Sha1sum(b)
	)

	res, err := config.Compare(a, b)
	require.NoError(t, err)

	expected := config.CompareResult{
		config.CompareCompactorKey: ah,
		// Rest is new hash
		config.CompareDistributorKey:   bh,
		config.CompareIngesterKey:      bh,
		config.CompareQuerierKey:       bh,
		config.CompareQueryFrontendKey: bh,
	}
	require.Exactly(t, expected, res)
}

func TestCompare_WhenIngesterConfigNotEqual_ReturnOnlyIngester(t *testing.T) {
	var (
		a = []byte(`
ingester:
  lifecycler:
    ring:
      replication_factor: 2
`)
		b = []byte(`
ingester:
  lifecycler:
    ring:
     replication_factor: 1
`)
		ah, _ = config.Sha1sum(a)
		bh, _ = config.Sha1sum(b)
	)

	res, err := config.Compare(a, b)
	require.NoError(t, err)

	expected := config.CompareResult{
		config.CompareIngesterKey: bh,
		// Rest is old hash
		config.CompareCompactorKey:     ah,
		config.CompareDistributorKey:   ah,
		config.CompareQuerierKey:       ah,
		config.CompareQueryFrontendKey: ah,
	}
	require.Exactly(t, expected, res)
}

func TestCompare_WhenStorageConfigNotEqual_ReturnAll(t *testing.T) {
	var (
		a = []byte(`
storage_config:
  aws:
    region: us-nowhere-42
`)
		b = []byte(`
storage_config:
  aws:
    region: us-somewhere-42
`)
		bh, _ = config.Sha1sum(b)
	)

	res, err := config.Compare(a, b)
	require.NoError(t, err)

	expected := config.CompareResult{
		config.CompareCompactorKey:     bh,
		config.CompareDistributorKey:   bh,
		config.CompareIngesterKey:      bh,
		config.CompareQuerierKey:       bh,
		config.CompareQueryFrontendKey: bh,
	}
	require.Exactly(t, expected, res)
}
