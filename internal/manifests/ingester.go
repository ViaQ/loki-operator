package manifests

import (
	"fmt"
	"path"

	"github.com/ViaQ/logerr/kverrors"
	"github.com/ViaQ/loki-operator/internal/manifests/internal/config"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildIngester builds the k8s objects required to run Loki Ingester
func BuildIngester(opts Options) ([]client.Object, error) {
	ss, err := NewIngesterStatefulSet(opts)
	if err != nil {
		return nil, err
	}

	return []client.Object{
		ss,
		NewIngesterGRPCService(opts),
		NewIngesterHTTPService(opts),
	}, nil
}

// NewIngesterStatefulSet creates a deployment object for an ingester
func NewIngesterStatefulSet(opts Options) (*apps.StatefulSet, error) {
	podSpec := core.PodSpec{
		Volumes: []core.Volume{
			{
				Name: configVolumeName,
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: lokiConfigMapName(opts.Name),
						},
					},
				},
			},
		},
		Containers: []core.Container{
			{
				Image: containerImage,
				Name:  "loki-ingester",
				Args: []string{
					"-target=ingester",
					fmt.Sprintf("-config.file=%s", path.Join(config.LokiConfigMountDir, config.LokiConfigFileName)),
				},
				ReadinessProbe: &core.Probe{
					Handler: core.Handler{
						HTTPGet: &core.HTTPGetAction{
							Path:   "/ready",
							Port:   intstr.FromInt(httpPort),
							Scheme: core.URISchemeHTTP,
						},
					},
					InitialDelaySeconds: 15,
					TimeoutSeconds:      1,
				},
				LivenessProbe: &core.Probe{
					Handler: core.Handler{
						HTTPGet: &core.HTTPGetAction{
							Path:   "/metrics",
							Port:   intstr.FromInt(httpPort),
							Scheme: core.URISchemeHTTP,
						},
					},
					TimeoutSeconds:   2,
					PeriodSeconds:    30,
					FailureThreshold: 10,
				},
				Ports: []core.ContainerPort{
					{
						Name:          "metrics",
						ContainerPort: httpPort,
					},
					{
						Name:          "grpc",
						ContainerPort: grpcPort,
					},
					{
						Name:          "gossip-ring",
						ContainerPort: gossipPort,
					},
				},
				// Resources: core.ResourceRequirements{
				// 	Limits: core.ResourceList{
				// 		core.ResourceMemory: resource.MustParse("1Gi"),
				// 		core.ResourceCPU:    resource.MustParse("1000m"),
				// 	},
				// 	Requests: core.ResourceList{
				// 		core.ResourceMemory: resource.MustParse("50m"),
				// 		core.ResourceCPU:    resource.MustParse("50m"),
				// 	},
				// },
				VolumeMounts: []core.VolumeMount{
					{
						Name:      configVolumeName,
						ReadOnly:  false,
						MountPath: config.LokiConfigMountDir,
					},
					{
						Name:      storageVolumeName,
						ReadOnly:  false,
						MountPath: dataDirectory,
					},
				},
			},
		},
	}

	ingesterLabels := ComponentLabels("ingester", opts.Name)

	storageRequests, err := resource.ParseQuantity(opts.Ingester.Storage.SizeRequested)
	if err != nil {
		return nil, kverrors.Wrap(err, "failed to parse quantity specified in Options", "field", "Ingester.StorageClass")
	}

	return &apps.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: apps.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("loki-ingester-%s", opts.Name),
			Labels: ingesterLabels,
		},
		Spec: apps.StatefulSetSpec{
			PodManagementPolicy:  apps.OrderedReadyPodManagement,
			RevisionHistoryLimit: pointer.Int32Ptr(10),
			Replicas:             pointer.Int32Ptr(int32(3)),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels.Merge(ingesterLabels, GossipLabels()),
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   fmt.Sprintf("loki-ingester-%s", opts.Name),
					Labels: labels.Merge(ingesterLabels, GossipLabels()),
				},
				Spec: podSpec,
			},
			VolumeClaimTemplates: []core.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Labels: ingesterLabels,
						Name:   storageVolumeName,
					},
					Spec: core.PersistentVolumeClaimSpec{
						AccessModes: []core.PersistentVolumeAccessMode{
							// TODO: should we verify that this is possible with the given storage class first?
							core.ReadWriteOnce,
						},
						Resources: core.ResourceRequirements{
							Requests: map[core.ResourceName]resource.Quantity{
								core.ResourceStorage: storageRequests,
							},
						},
						StorageClassName: pointer.StringPtr(opts.Ingester.Storage.ClassName),
					},
				},
			},
		},
	}, nil
}

// NewIngesterGRPCService creates a k8s service for the ingester GRPC endpoint
func NewIngesterGRPCService(opts Options) *core.Service {
	l := ComponentLabels("ingester", opts.Name)
	return &core.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: core.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceNameIngesterGRPC(opts.Name),
			Labels: l,
		},
		Spec: core.ServiceSpec{
			ClusterIP: "None",
			Ports: []core.ServicePort{
				{
					Name: "grpc",
					Port: grpcPort,
				},
			},
			Selector: l,
		},
	}
}

// NewIngesterHTTPService creates a k8s service for the ingester HTTP endpoint
func NewIngesterHTTPService(opts Options) *core.Service {
	l := ComponentLabels("ingester", opts.Name)
	return &core.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: core.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceNameIngesterHTTP(opts.Name),
			Labels: l,
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Name: "metrics",
					Port: httpPort,
				},
			},
			Selector: l,
		},
	}
}
