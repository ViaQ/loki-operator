package manifests

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

const (
	// LokiGatewayMountDir is the path that is mounted from the configmap
	LokiGatewayMountDir = "/etc/observatorium"
)

// BuildLokiStackGateway returns a list of k8s objects for Loki Stack Gateway
func BuildLokiStackGateway(opts Options) *appsv1.Deployment {
	return NewLokiStackGatewayDeployment(opts)
}

// NewLokiStackGatewayDeployment creates a deployment object for a lokiStack-gateway
func NewLokiStackGatewayDeployment(opts Options) *appsv1.Deployment {
	podSpec := corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: "rbac",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "observatorium-api",
						},
					},
				},
			},
			{
				Name: "tenants",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "observatorium-api",
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:  "lokistack-gateway",
				Image: DefaultLokiStackGatewayImage,
				Args: []string{
					"--debug.name=lokistack-gateway",
					"--web.listen=0.0.0.0:8080",
					"--web.internal.listen=0.0.0.0:8081",
					"--log.level=debug",
					fmt.Sprintf("--logs.read.endpoint=http://%s:%d", fqdn(serviceNameQueryFrontendHTTP(opts.Name), opts.Namespace), httpPort),
					fmt.Sprintf("--logs.tail.endpoint=http://%s:%d", fqdn(serviceNameQuerierHTTP(opts.Name), opts.Namespace), httpPort),
					fmt.Sprintf("--logs.write.endpoint=http://%s:%d", fqdn(serviceNameDistributorHTTP(opts.Name), opts.Namespace), httpPort),
					"--rbac.config=/etc/observatorium/rbac.yaml",
					"--tenants.config=/etc/observatorium/tenants.yaml",
				},
				WorkingDir: LokiGatewayMountDir,
				Ports: []corev1.ContainerPort{
					{
						Name:          "internal",
						ContainerPort: 8081,
					},
					{
						Name:          "public",
						ContainerPort: 8080,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "rbac",
						ReadOnly:  true,
						MountPath: fmt.Sprintf("%s/rbac.yaml", LokiGatewayMountDir),
						SubPath:   "rbac.yaml",
					},
					{
						Name:      "tenants",
						ReadOnly:  true,
						MountPath: fmt.Sprintf("%s/tenants.yaml", LokiGatewayMountDir),
						SubPath:   "tenants.yaml",
					},
				},
				LivenessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path:   "/live",
							Port:   intstr.FromInt(8081),
							Scheme: corev1.URISchemeHTTP,
						},
					},
					TimeoutSeconds:   2,
					PeriodSeconds:    30,
					FailureThreshold: 10,
				},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path:   "/ready",
							Port:   intstr.FromInt(8081),
							Scheme: corev1.URISchemeHTTP,
						},
					},
					TimeoutSeconds:   1,
					PeriodSeconds:    5,
					FailureThreshold: 12,
				},
			},
		},
	}

	l := ComponentLabels(LabelLokiStackGatewayComponent, opts.Name)
	a := commonAnnotations(opts.ConfigSHA1)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   LokiStackGatewayName(opts.Name),
			Labels: l,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels.Merge(l, GossipLabels()),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        LokiStackGatewayName(opts.Name),
					Labels:      labels.Merge(l, GossipLabels()),
					Annotations: a,
				},
				Spec: podSpec,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
		},
	}
}
