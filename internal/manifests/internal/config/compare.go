package config

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

// CompareResult represents the loki config comparison result type.
// It defines a mapping from component names to config hashes.
type CompareResult map[string]string

const (
	// CompareCompactorKey is the key to access the compactor component hash
	// from a loki configuration comparison result.
	CompareCompactorKey string = "compactor"
	// CompareDistributorKey is the key to access the distributor component hash
	// from a loki configuration comparison result.
	CompareDistributorKey string = "distributor"
	// CompareIngesterKey is the key to access the ingester component hash
	// from a loki configuration comparison result.
	CompareIngesterKey string = "ingester"
	// CompareQuerierKey is the key to access the querier component hash
	// from a loki configuration comparison result.
	CompareQuerierKey string = "querier"
	// CompareQueryFrontendKey is the key to access the query-frontend component hash
	// from a loki configuration comparison result.
	CompareQueryFrontendKey string = "query-frontend"
)

var configToCompopnentMap = map[string][]string{
	"LimitsConfig":  {CompareDistributorKey, CompareIngesterKey, CompareQuerierKey, CompareQueryFrontendKey},
	"Ingester":      {CompareIngesterKey},
	"StorageConfig": {CompareCompactorKey, CompareDistributorKey, CompareIngesterKey, CompareQuerierKey, CompareQueryFrontendKey},
}

// Compare loads and compares two loki configuration byte slices and returns
// a mapping of affected components to the corresponding config SHA1 hash.
// A component is getting a new hash only if the diff between the byte slices
// affects the component itself. Currently the compare result supports changes:
// - Limits config: Affect all component except compactor.
// - Ingester config: Affects only ingester component.
// - Storage config: Affects all components.
func Compare(old, new []byte) (CompareResult, error) {
	// Unmarshal old and new config
	// contents into Loki objects
	var o, n Loki
	err := yaml.Unmarshal(old, &o)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(new, &n)
	if err != nil {
		return nil, err
	}

	// Hash old and new config contents
	ohash, err := Sha1sum(old)
	if err != nil {
		return nil, err
	}
	nhash, err := Sha1sum(new)
	if err != nil {
		return nil, err
	}

	// Return all components with old hash
	// if config contents are equal.
	var r PathReporter
	ok := cmp.Equal(o, n, cmp.Reporter(&r))
	if ok {
		return all(ohash), nil
	}

	// Add new hash to affected components
	res := CompareResult{}
	for _, rr := range r.Roots() {
		for _, c := range configToCompopnentMap[rr] {
			res[c] = nhash
		}
	}

	// Add old hash to non-affected components
	for k, v := range all(ohash) {
		_, ok := res[k]
		if !ok {
			res[k] = v
		}
	}
	return res, nil
}

func all(h string) CompareResult {
	return CompareResult{
		CompareCompactorKey:     h,
		CompareDistributorKey:   h,
		CompareIngesterKey:      h,
		CompareQuerierKey:       h,
		CompareQueryFrontendKey: h,
	}
}
