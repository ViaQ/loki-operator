package handlers

import (
	"context"
	"fmt"

	"github.com/ViaQ/logerr/kverrors"
	"github.com/ViaQ/logerr/log"
	"github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/ViaQ/loki-operator/internal/external/k8s"
	"github.com/ViaQ/loki-operator/internal/handlers/internal/utils"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	prometheusCAFile = "/etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt"
	bearerTokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"

	labelJobName = "monitor-loki"
)

// CreateOrUpdateServiceMonitor ensures the existence of ServiceMonitors for Loki cluster
func CreateOrUpdateServiceMonitor(ctx context.Context, req ctrl.Request, k k8s.Client) error {
	var ls v1beta1.LokiStack
	if err := k.Get(ctx, req.NamespacedName, &ls); err != nil {
		if apierrors.IsNotFound(err) {
			// maybe the user deleted it before we could react? Either way this isn't an issue
			log.Error(err, "could not find the requested loki stack", "name", req.NamespacedName)
			return nil
		}
		return kverrors.Wrap(err, "failed to lookup lokistack", "name", req.NamespacedName)
	}

	serviceMonitorName := fmt.Sprintf("monitor-%s-%s", ls.Name, "cluster")

	labelsWithDefault := utils.AppendDefaultLabel(ls.Name, ls.Labels)
	labelsWithDefault["scrape-metrics"] = "enabled"

	lokiSvcMonitor := createServiceMonitor(serviceMonitorName, ls.Name, ls.Namespace, labelsWithDefault)
	utils.AddOwnerRefTo(&ls, lokiSvcMonitor)

	err := k.Create(context.TODO(), lokiSvcMonitor)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return kverrors.Wrap(err, "failed to construct loki ServiceMonitor")
	}

	return nil
}

func createServiceMonitor(serviceMonitorName, clusterName, namespace string, labels map[string]string) *monitoringv1.ServiceMonitor {
	svcMonitor := newServiceMonitor(serviceMonitorName, namespace, labels)
	labelSelector := metav1.LabelSelector{
		MatchLabels: labels,
	}
	tlsConfig := monitoringv1.TLSConfig{
		SafeTLSConfig: monitoringv1.SafeTLSConfig{
			// ServerName can be e.g. loki-metrics.openshift-logging.svc
			ServerName: fmt.Sprintf("%s-%s.%s.svc", clusterName, "metrics", namespace),
		},
		CAFile: prometheusCAFile,
	}
	lokiEndpoint := monitoringv1.Endpoint{
		Port:            clusterName,
		Path:            "/metrics",
		Scheme:          "https",
		BearerTokenFile: bearerTokenFile,
		TLSConfig:       &tlsConfig,
	}
	svcMonitor.Spec = monitoringv1.ServiceMonitorSpec{
		JobLabel:  labelJobName,
		Endpoints: []monitoringv1.Endpoint{lokiEndpoint},
		Selector:  labelSelector,
		NamespaceSelector: monitoringv1.NamespaceSelector{
			MatchNames: []string{namespace},
		},
	}
	return svcMonitor
}

func newServiceMonitor(serviceMonitorName string, namespace string, labels map[string]string) *monitoringv1.ServiceMonitor {
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
	}
}
