package manifests

// Options is a set of options to use when building manifests such as resource sizes, etc.
// Most of this should be provided - either directly or indirectly - by the user. This will
// probably be converted from the CR.
type Options struct {
	Name      string
	Namespace string

	Ingester Ingester `default:"{}"`
	Querier  Querier  `default:"{}"`
}

// Storage defines PVC settings for StatefulSets
type Storage struct {
	ClassName     string `default:"-"`
	SizeRequested string `default:"10Gi"`
}

// Querier is options for the Querier
type Querier struct {
	Storage Storage
}

// Ingester is options for the Ingester
type Ingester struct {
	Storage Storage
}
