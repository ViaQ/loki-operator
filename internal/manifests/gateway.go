package manifests

import (
	"crypto/sha1"
	"fmt"
	"path"

	"github.com/ViaQ/logerr/kverrors"
	"github.com/imdario/mergo"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ViaQ/loki-operator/internal/manifests/internal/gateway"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

// BuildLokiStackGateway returns a list of k8s objects for Loki Stack Gateway
func BuildLokiStackGateway(opts Options) ([]client.Object, string, error) {
	gatewayConfigMap, sha1C, err := GatewayConfigMap(opts)
	if err != nil {
		return nil, "", err
	}

	deployment := NewLokiStackGatewayDeployment(opts)
	if opts.Flags.EnableTLSLokiStackGateway {
		if err := configureLokiStackGatewayPKI(&deployment.Spec.Template.Spec); err != nil {
			return nil, "", err
		}
	}

	return []client.Object{
		gatewayConfigMap,
		deployment,
	}, sha1C, nil
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
							Name: LabelLokiStackGatewayComponent,
						},
					},
				},
			},
			{
				Name: "tenants",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: LabelLokiStackGatewayComponent,
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:  LabelLokiStackGatewayComponent,
				Image: DefaultLokiStackGatewayImage,
				Args: []string{
					fmt.Sprintf("--debug.name=%s", LabelLokiStackGatewayComponent),
					"--web.listen=0.0.0.0:8080",
					"--web.internal.listen=0.0.0.0:8081",
					"--log.level=debug",
					fmt.Sprintf("--logs.read.endpoint=http://%s:%d", fqdn(serviceNameQueryFrontendHTTP(opts.Name), opts.Namespace), httpPort),
					fmt.Sprintf("--logs.tail.endpoint=http://%s:%d", fqdn(serviceNameQueryFrontendHTTP(opts.Name), opts.Namespace), httpPort),
					fmt.Sprintf("--logs.write.endpoint=http://%s:%d", fqdn(serviceNameDistributorHTTP(opts.Name), opts.Namespace), httpPort),
					fmt.Sprintf("--rbac.config=%s", path.Join(gateway.LokiGatewayMountDir, gateway.LokiGatewayRbacFileName)),
					fmt.Sprintf("--tenants.config=%s", path.Join(gateway.LokiGatewayMountDir, gateway.LokiGatewayTenantFileName)),
				},
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
						MountPath: path.Join(gateway.LokiGatewayMountDir, gateway.LokiGatewayRbacFileName),
						SubPath:   "rbac.yaml",
					},
					{
						Name:      "tenants",
						ReadOnly:  true,
						MountPath: path.Join(gateway.LokiGatewayMountDir, gateway.LokiGatewayTenantFileName),
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

// GatewayConfigMap creates a configMap for rbac.yaml and tenants.yaml
func GatewayConfigMap(opt Options) (*corev1.ConfigMap, string, error) {
	cfg := GatewayConfigOptions(opt)
	rbacConfig, tenantsConfig, err := gateway.Build(cfg)
	if err != nil {
		return nil, "", err
	}

	s := sha1.New()
	_, err = s.Write(rbacConfig)
	if err != nil {
		return nil, "", err
	}
	sha1C := fmt.Sprintf("%x", s.Sum(nil))

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   LabelLokiStackGatewayComponent,
			Labels: commonLabels(opt.Name),
		},
		BinaryData: map[string][]byte{
			gateway.LokiGatewayRbacFileName:   rbacConfig,
			gateway.LokiGatewayTenantFileName: tenantsConfig,
		},
	}, sha1C, nil
}

// GatewayConfigOptions converts Options to gateway.Options
func GatewayConfigOptions(opt Options) gateway.Options {
	return gateway.Options{
		Stack:     opt.Stack,
		Namespace: opt.Namespace,
		Name:      opt.Name,
	}
}

func configureLokiStackGatewayPKI(podSpec *corev1.PodSpec) error {
	secretVolumeSpec := corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: "tls-secret",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: LabelLokiStackGatewayComponent,
					},
				},
			},
			{
				Name: "tls-configmap",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: LabelLokiStackGatewayComponent,
						},
					},
				},
			},
		},
	}
	secretContainerSpec := corev1.Container{
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "tls-secret",
				ReadOnly:  true,
				MountPath: path.Join(gateway.LokiGatewayTLSDir, "cert"),
				SubPath:   "cert",
			},
			{
				Name:      "tls-secret",
				ReadOnly:  true,
				MountPath: path.Join(gateway.LokiGatewayTLSDir, "key"),
				SubPath:   "key",
			},
			{
				Name:      "tls-configmap",
				ReadOnly:  true,
				MountPath: path.Join(gateway.LokiGatewayTLSDir, "ca"),
				SubPath:   "ca",
			},
		},
		Args: []string{
			fmt.Sprintf("--tls.server.cert-file=%s", path.Join(gateway.LokiGatewayTLSDir, "cert")),
			fmt.Sprintf("--tls.server.key-file=%s", path.Join(gateway.LokiGatewayTLSDir, "key")),
			fmt.Sprintf("--tls.healthchecks.server-ca-file=%s", path.Join(gateway.LokiGatewayTLSDir, "ca")),
		},
	}

	if err := mergo.Merge(podSpec, secretVolumeSpec, mergo.WithAppendSlice); err != nil {
		return kverrors.Wrap(err, "failed to merge volumes")
	}

	if err := mergo.Merge(&podSpec.Containers[0], secretContainerSpec, mergo.WithAppendSlice); err != nil {
		return kverrors.Wrap(err, "failed to merge container")
	}

	return nil
}
