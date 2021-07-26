package gateway

import (
	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
)

// Options is used to render the loki-config.yaml file template
type Options struct {
	Stack lokiv1beta1.LokiStackSpec

	Namespace        string
	Name             string
	StorageDirectory string
}
