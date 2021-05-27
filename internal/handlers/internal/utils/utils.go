package utils

import (
	"github.com/ViaQ/loki-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// AppendDefaultLabel appends the label "cluster-name"=clusterName if not already present
func AppendDefaultLabel(clusterName string, labels map[string]string) map[string]string {
	if _, ok := labels["cluster-name"]; ok {
		return labels
	}
	if labels == nil {
		labels = map[string]string{}
	}
	labels["cluster-name"] = clusterName
	return labels
}

// AddOwnerRefTo appends the LokiStack object as an OwnerReference to the passed object
func AddOwnerRefTo(ls *v1beta1.LokiStack, o metav1.Object) {
	ref := metav1.OwnerReference{
		APIVersion: v1beta1.GroupVersion.String(),
		Kind:       "Loki",
		Name:       ls.Name,
		UID:        ls.UID,
		Controller: pointer.BoolPtr(true),
	}
	if (metav1.OwnerReference{}) != ref {
		o.SetOwnerReferences(append(o.GetOwnerReferences(), ref))
	}
}
