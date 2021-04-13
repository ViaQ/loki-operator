package manifests

import (
	"fmt"
	"strings"

	"github.com/ViaQ/logerr/kverrors"
	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/ViaQ/loki-operator/internal/manifests/internal/config"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LokiConfigMap creates the single configmap containing the loki configuration for the whole cluster
func LokiConfigMap(opt Options) (*corev1.ConfigMap, error) {
	var cfg config.Options
	if err := mergo.Merge(&cfg, configForSize(opt.Name, opt.Namespace, opt.Stack.Size)); err != nil {
		return nil, kverrors.Wrap(err, "failed to merge configs")
	}

	b, err := config.Build(cfg)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   lokiConfigMapName(opt.Name),
			Labels: commonLabels(opt.Name),
		},
		BinaryData: map[string][]byte{
			config.LokiConfigFileName: b,
		},
	}, nil
}

func lokiConfigMapName(stackName string) string {
	return fmt.Sprintf("loki-config-%s", stackName)
}

func configForSize(name, namespace string, sizeType lokiv1beta1.LokiStackSizeType) config.Options {
	// TODO switch on size
	return config.Options{
		Namespace: namespace,
		Name:      name,
		FrontendWorker: config.Address{
			FQDN: "",
			Port: 0,
		},
		GossipRing: config.Address{
			FQDN: fqdn(LokiGossipRingService(name).GetName(), namespace),
			Port: gossipPort,
		},
		Querier: config.Address{
			FQDN: serviceNameQuerierHTTP(name),
			Port: httpPort,
		},
		StorageDirectory: strings.TrimRight(dataDirectory, "/"),
		Spec: lokiv1beta1.LokiStackSpec{
			Size:              sizeType,
			ReplicationFactor: 2,
			Limits: lokiv1beta1.LimitsSpec{
				Global: lokiv1beta1.LimitsTemplateSpec{
					IngestionLimits: lokiv1beta1.IngestionLimitSpec{
						IngestionRate:      20,
						IngestionBurstSize: 10,
						MaxStreamsPerUser:  25000,
					},
					QueryLimits: lokiv1beta1.QueryLimitSpec{
						MaxEntriesPerQuery: 0,
						MaxChunksPerQuery:  0,
						MaxQuerySeries:     0,
					},
				},
			},
			Template: lokiv1beta1.LokiTemplateSpec{
				Compactor: lokiv1beta1.LokiComponentSpec{
					Replicas: 3,
				},
				Distributor: lokiv1beta1.LokiComponentSpec{
					Replicas: 3,
				},
				Ingester: lokiv1beta1.LokiComponentSpec{
					Replicas: 3,
				},
				Querier: lokiv1beta1.LokiComponentSpec{
					Replicas: 3,
				},
				QueryFrontend: lokiv1beta1.LokiComponentSpec{
					Replicas: 3,
				},
			},
		},
	}
}
