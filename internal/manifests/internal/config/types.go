package config

// Loki reprsents the loki configuration file.
type Loki struct {
	Ingester      Ingester      `yaml:"ingester"`
	LimitsConfig  LimitsConfig  `yaml:"limits_config"`
	StorageConfig StorageConfig `yaml:"storage_config"`
}

// Ingester represents the ingester configuration section.
type Ingester struct {
	Lifecycler Lifecycler `yaml:"lifecycler"`
}

// Lifecycler represents the ingester lifecycler configuration section.n
type Lifecycler struct {
	Ring Ring `yaml:"ring"`
}

// Ring represents the ingester ring configuration section.
type Ring struct {
	ReplicationFactor int32 `yaml:"replication_factor"`
}

// LimitsConfig represents the limits_config section.
type LimitsConfig struct {
	IngestionRate           int32 `yaml:"ingestion_rate_mb"`
	IngestionBurstSize      int32 `yaml:"ingestion_burst_size_mb"`
	MaxLabelNameLength      int32 `yaml:"max_label_name_length"`
	MaxLabelValueLength     int32 `yaml:"max_label_value_length"`
	MaxLabelNamesPerSeries  int32 `yaml:"max_label_names_per_series"`
	MaxStreamsPerUser       int32 `yaml:"max_streams_per_user"`
	MaxLineSize             int32 `yaml:"max_line_size"`
	MaxGlobalStreamsPerUser int32 `yaml:"max_global_streams_per_user"`
	MaxEntriesLimitPerQuery int32 `yaml:"max_entries_limit_per_query"`
	MaxQuerySeries          int32 `yaml:"max_query_series"`
}

// StorageConfig represents the storage_config section.
type StorageConfig struct {
	AWS AWS `yaml:"aws"`
}

// AWS represents the aws storage_config section.
type AWS struct {
	S3              string `yaml:"s3"`
	BucketNames     string `yaml:"bucketnames"`
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}
