package manifests

import (
	"github.com/ViaQ/logerr/kverrors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/imdario/mergo"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

// BuildServiceMonitors builds the service monitors
func BuildServiceMonitors(opt Options) []client.Object {
	return []client.Object{
		NewDistributorServiceMonitor(opt),
		NewIngesterServiceMonitor(opt),
		NewQuerierServiceMonitor(opt),
		NewCompactorServiceMonitor(opt),
		NewQueryFrontendServiceMonitor(opt),
	}
}

// NewDistributorServiceMonitor creates a k8s service monitor for the distributor component
func NewDistributorServiceMonitor(opt Options) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelDistributorComponent, opt.Name)

	serviceMonitorName := serviceMonitorName(DistributorName(opt.Name))
	serviceName := serviceNameDistributorHTTP(opt.Name)
	lokiEndpoint := serviceMonitorLokiEndPoint(opt.Name, serviceName, opt.Namespace, opt.EnableTLSServiceMonitorConfig)

	return newServiceMonitor(opt.Namespace, serviceMonitorName, l, lokiEndpoint)
}

// NewIngesterServiceMonitor creates a k8s service monitor for the ingester component
func NewIngesterServiceMonitor(opt Options) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelIngesterComponent, opt.Name)

	serviceMonitorName := serviceMonitorName(IngesterName(opt.Name))
	serviceName := serviceNameIngesterHTTP(opt.Name)
	lokiEndpoint := serviceMonitorLokiEndPoint(opt.Name, serviceName, opt.Namespace, opt.EnableTLSServiceMonitorConfig)

	return newServiceMonitor(opt.Namespace, serviceMonitorName, l, lokiEndpoint)
}

// NewQuerierServiceMonitor creates a k8s service monitor for the querier component
func NewQuerierServiceMonitor(opt Options) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelQuerierComponent, opt.Name)

	serviceMonitorName := serviceMonitorName(QuerierName(opt.Name))
	serviceName := serviceNameQuerierHTTP(opt.Name)
	lokiEndpoint := serviceMonitorLokiEndPoint(opt.Name, serviceName, opt.Namespace, opt.EnableTLSServiceMonitorConfig)

	return newServiceMonitor(opt.Namespace, serviceMonitorName, l, lokiEndpoint)
}

// NewCompactorServiceMonitor creates a k8s service monitor for the compactor component
func NewCompactorServiceMonitor(opt Options) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelCompactorComponent, opt.Name)

	serviceMonitorName := serviceMonitorName(CompactorName(opt.Name))
	serviceName := serviceNameCompactorHTTP(opt.Name)
	lokiEndpoint := serviceMonitorLokiEndPoint(opt.Name, serviceName, opt.Namespace, opt.EnableTLSServiceMonitorConfig)

	return newServiceMonitor(opt.Namespace, serviceMonitorName, l, lokiEndpoint)
}

// NewQueryFrontendServiceMonitor creates a k8s service monitor for the query-frontend component
func NewQueryFrontendServiceMonitor(opt Options) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelQueryFrontendComponent, opt.Name)

	serviceMonitorName := serviceMonitorName(QueryFrontendName(opt.Name))
	serviceName := serviceNameQueryFrontendHTTP(opt.Name)
	lokiEndpoint := serviceMonitorLokiEndPoint(opt.Name, serviceName, opt.Namespace, opt.EnableTLSServiceMonitorConfig)

	return newServiceMonitor(opt.Namespace, serviceMonitorName, l, lokiEndpoint)
}

func newServiceMonitor(namespace, serviceMonitorName string, labels labels.Set, endpoint monitoringv1.Endpoint) *monitoringv1.ServiceMonitor {
	return &monitoringv1.ServiceMonitor{
		TypeMeta: metav1.TypeMeta{
			Kind:       monitoringv1.ServiceMonitorsKind,
			APIVersion: monitoringv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceMonitorName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			JobLabel:  labelJobComponent,
			Endpoints: []monitoringv1.Endpoint{endpoint},
			Selector: metav1.LabelSelector{
				MatchLabels: labels,
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{namespace},
			},
		},
	}
}

func configureServiceMonitorPKI(podSpec *corev1.PodSpec, serviceName string) error {
	secretName := signingServiceSecretName(serviceName)
	secretVolumeSpec := corev1.PodSpec{
		Volumes: []corev1.Volume{
			{
				Name: secretName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: secretName,
					},
				},
			},
		},
	}
	secretContainerSpec := corev1.Container{
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      secretName,
				ReadOnly:  false,
				MountPath: secretDirectory,
			},
		},
		Args: []string{
			"-server.http-tls-cert-path=/etc/proxy/secrets/tls.crt",
			"-server.http-tls-key-path=/etc/proxy/secrets/tls.key",
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
