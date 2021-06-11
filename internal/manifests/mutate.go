package manifests

import (
	"reflect"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	"github.com/ViaQ/logerr/kverrors"
	// "github.com/imdario/mergo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// MutateFuncFor returns a mutate function based on the
// existing resource's concrete type. It supports currently
// only the following types or else panics:
// - ConfigMap
// - Service
// - Deployment
// - StatefulSet
// - ServiceMonitor
func MutateFuncFor(existing, desired client.Object) controllerutil.MutateFn {
	return func() error {
		existing.SetAnnotations(desired.GetAnnotations())
		existing.SetLabels(desired.GetLabels())

		switch existing.(type) {
		case *corev1.ConfigMap:
			cm := existing.(*corev1.ConfigMap)
			wantCm := desired.(*corev1.ConfigMap)
			mutateConfigMap(cm, wantCm)

		case *corev1.Service:
			svc := existing.(*corev1.Service)
			wantSvc := desired.(*corev1.Service)
			mutateService(svc, wantSvc)

		case *appsv1.Deployment:
			dpl := existing.(*appsv1.Deployment)
			wantDpl := desired.(*appsv1.Deployment)
			mutateDeployment(dpl, wantDpl)

		case *appsv1.StatefulSet:
			sts := existing.(*appsv1.StatefulSet)
			wantSts := desired.(*appsv1.StatefulSet)
			mutateStatefulSet(sts, wantSts)

		case *monitoringv1.ServiceMonitor:
			svcMonitor := existing.(*monitoringv1.ServiceMonitor)
			wantSvcMonitor := desired.(*monitoringv1.ServiceMonitor)
			mutateServiceMonitor(svcMonitor, wantSvcMonitor)

		default:
			t := reflect.TypeOf(existing).String()
			return kverrors.New("missing mutate implementation for resource type", "type", t)
		}
		return nil
	}
}

func mutateConfigMap(existing, desired *corev1.ConfigMap) {
	existing.BinaryData = desired.BinaryData
}

func mutateService(existing, desired *corev1.Service) {
	existing.Spec.Ports = desired.Spec.Ports
	existing.Spec.Selector = desired.Spec.Selector
}

func mutateDeployment(existing, desired *appsv1.Deployment) {
	// Deployment selector is immutable so we set this value only if
	// a new object is going to be created
	if existing.CreationTimestamp.IsZero() {
		existing.Spec.Selector = desired.Spec.Selector
	}
	existing.Spec.Replicas = desired.Spec.Replicas
	existing.Spec.Template = desired.Spec.Template
	existing.Spec.Strategy = desired.Spec.Strategy
}

func mutateStatefulSet(existing, desired *appsv1.StatefulSet) {
	// StatefulSet selector is immutable so we set this value only if
	// a new object is going to be created
	if existing.CreationTimestamp.IsZero() {
		existing.Spec.Selector = desired.Spec.Selector
	}
	existing.Spec.PodManagementPolicy = desired.Spec.PodManagementPolicy
	existing.Spec.Replicas = desired.Spec.Replicas
	existing.Spec.Template = desired.Spec.Template
	existing.Spec.VolumeClaimTemplates = desired.Spec.VolumeClaimTemplates
}

func mutateServiceMonitor(existing, desired *monitoringv1.ServiceMonitor) {
	// ServiceMonitor selector is immutable so we set this value only if
	// a new object is going to be created
}

func mutatePodSpecForTLSEnablement(podSpec *corev1.PodSpec, serviceName string) {
	secretName := signingServiceSecretName(serviceName)

	podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
		Name: secretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	})
	podSpec.Containers[0].VolumeMounts = append(podSpec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      secretName,
		ReadOnly:  false,
		MountPath: "/etc/proxy/secrets",
	})
	// podSpec.Containers[0].Args = append(podSpec.Containers[0].Args, "-server.http-tls-ca-path=/etc/proxy/secrets/ca-bundle.crt")
	podSpec.Containers[0].Args = append(podSpec.Containers[0].Args, "-server.http-tls-cert-path=/etc/proxy/secrets/tls.crt")
	podSpec.Containers[0].Args = append(podSpec.Containers[0].Args, "-server.http-tls-key-path=/etc/proxy/secrets/tls.key")

	// tlsSpec := corev1.PodSpec{
	// 	Volumes: []corev1.Volume{
	// 		,
	// 	},
	// 	Containers: []corev1.Container{
	// 		{
	// 			VolumeMounts: []corev1.VolumeMount{
	// 				{
	// 					Name:      secretName,
	// 					ReadOnly:  false,
	// 					MountPath: "/etc/proxy/secrets",
	// 				},
	// 			},
	// 			Args: []string{
	// 				"-server.http-tls-ca-path=/etc/proxy/secrets/ca-bundle.crt",
	// 				"-server.http-tls-cert-path=/etc/proxy/secrets/tls.crt",
	// 				"-server.http-tls-key-path=/etc/proxy/secrets/tls.key",
	// 			},
	// 		},
	// 	},
	// }

	// if err := mergo.Merge(podSpec.Volumes, tlsSpec.Volumes, mergo.WithAppendSlice); err != nil {
	// 	return kverrors.Wrap(err, "failed to merge volumes")
	// }
	//
	// if err := mergo.Merge(podSpec.Containers[0], tlsSpec.Containers[0], mergo.WithAppendSlice); err != nil {
	// 	return kverrors.Wrap(err, "failed to merge containers")
	// }
	//
	// return nil
}
