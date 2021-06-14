package manifests

import (
	"fmt"

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
		NewDistributorServiceMonitor(opt.Name, opt.Namespace, opt.EnableTLSServiceMonitorConfig),
		NewIngesterServiceMonitor(opt.Name, opt.Namespace, opt.EnableTLSServiceMonitorConfig),
		NewQuerierServiceMonitor(opt.Name, opt.Namespace, opt.EnableTLSServiceMonitorConfig),
		NewCompactorServiceMonitor(opt.Name, opt.Namespace, opt.EnableTLSServiceMonitorConfig),
		NewQueryFrontendServiceMonitor(opt.Name, opt.Namespace, opt.EnableTLSServiceMonitorConfig),
	}
}

// NewDistributorServiceMonitor creates a k8s service monitor for the distributor component
func NewDistributorServiceMonitor(stackName, namespace string, useTLSConfig bool) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelDistributorComponent, stackName)

	serviceMonitorName := fmt.Sprintf("monitor-%s", DistributorName(stackName))
	serviceName := serviceNameDistributorHTTP(stackName)
	lokiEndpoint := serviceMonitorLokiEndPoint(stackName, serviceName, namespace, useTLSConfig)

	return newServiceMonitor(namespace, serviceMonitorName, l, lokiEndpoint)
}

// NewIngesterServiceMonitor creates a k8s service monitor for the ingester component
func NewIngesterServiceMonitor(stackName, namespace string, useTLSConfig bool) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelIngesterComponent, stackName)

	serviceMonitorName := fmt.Sprintf("monitor-%s", IngesterName(stackName))
	serviceName := serviceNameIngesterHTTP(stackName)
	lokiEndpoint := serviceMonitorLokiEndPoint(stackName, serviceName, namespace, useTLSConfig)

	return newServiceMonitor(namespace, serviceMonitorName, l, lokiEndpoint)
}

// NewQuerierServiceMonitor creates a k8s service monitor for the querier component
func NewQuerierServiceMonitor(stackName, namespace string, useTLSConfig bool) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelQuerierComponent, stackName)

	serviceMonitorName := fmt.Sprintf("monitor-%s", QuerierName(stackName))
	serviceName := serviceNameQuerierHTTP(stackName)
	lokiEndpoint := serviceMonitorLokiEndPoint(stackName, serviceName, namespace, useTLSConfig)

	return newServiceMonitor(namespace, serviceMonitorName, l, lokiEndpoint)
}

// NewCompactorServiceMonitor creates a k8s service monitor for the compactor component
func NewCompactorServiceMonitor(stackName, namespace string, useTLSConfig bool) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelCompactorComponent, stackName)

	serviceMonitorName := fmt.Sprintf("monitor-%s", CompactorName(stackName))
	serviceName := serviceNameCompactorHTTP(stackName)
	lokiEndpoint := serviceMonitorLokiEndPoint(stackName, serviceName, namespace, useTLSConfig)

	return newServiceMonitor(namespace, serviceMonitorName, l, lokiEndpoint)
}

// NewQueryFrontendServiceMonitor creates a k8s service monitor for the query-frontend component
func NewQueryFrontendServiceMonitor(stackName, namespace string, useTLSConfig bool) *monitoringv1.ServiceMonitor {
	l := ComponentLabels(LabelQueryFrontendComponent, stackName)

	serviceMonitorName := fmt.Sprintf("monitor-%s", QueryFrontendName(stackName))
	serviceName := serviceNameQueryFrontendHTTP(stackName)
	lokiEndpoint := serviceMonitorLokiEndPoint(stackName, serviceName, namespace, useTLSConfig)

	return newServiceMonitor(namespace, serviceMonitorName, l, lokiEndpoint)
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

func configureHTTPServiceCertSigning(service *corev1.Service, serviceName string) error {
	annotations := map[string]string{}
	annotations["service.beta.openshift.io/serving-cert-secret-name"] = signingServiceSecretName(serviceName)

	if err := mergo.Merge(&service.Annotations, annotations); err != nil {
		return kverrors.Wrap(err, "failed to merge http service annotations")
	}

	return nil
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
			// "-server.http-tls-ca-path=/etc/proxy/secrets/ca-bundle.crt",
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
