/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"github.com/ViaQ/loki-operator/controllers/internal/management/state"
	"github.com/ViaQ/loki-operator/internal/handlers"
	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
)

// LokiStackReconciler reconciles a LokiStack object
type LokiStackReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=loki.openshift.io,resources=lokistacks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=loki.openshift.io,resources=lokistacks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=loki.openshift.io,resources=lokistacks/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods;nodes;services;endpoints;configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings;clusterroles,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Compare the state specified by the LokiStack object against the actual cluster state,
// and then perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *LokiStackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ok, err := state.IsManaged(ctx, req, r.Client)
	if err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}
	if !ok {
		r.Log.Info("Skipping reconciliation for unmanaged lokistack resource", "name", req.NamespacedName)
		// Stop requeueing for unmanaged LokiStack custom resources
		return ctrl.Result{}, nil
	}

	err = handlers.CreateOrUpdateLokiStack(ctx, req, r.Client, r.Scheme)
	if err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LokiStackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	createOrUpdateOnlyPred := builder.WithPredicates(predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Update only if generation changes, filter out anything else.
			// We only need to check generation here, because it is only
			// updated on spec changes. On the other hand RevisionVersion
			// changes also on status changes. We want to omit reconciliation
			// for status updates for now.
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		CreateFunc:  func(e event.CreateEvent) bool { return true },
		DeleteFunc:  func(e event.DeleteEvent) bool { return false },
		GenericFunc: func(e event.GenericEvent) bool { return false },
	})
	deleteOnlyPred := builder.WithPredicates(predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool { return false },
		CreateFunc: func(e event.CreateEvent) bool { return false },
		DeleteFunc: func(e event.DeleteEvent) bool {
			// DeleteStateUnknown evaluates to false only if the object
			// has been confirmed as deleted by the api server.
			return !e.DeleteStateUnknown
		},
		GenericFunc: func(e event.GenericEvent) bool { return false },
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1beta1.LokiStack{}, createOrUpdateOnlyPred).
		Owns(&corev1.ConfigMap{}, deleteOnlyPred).
		Owns(&corev1.Service{}, deleteOnlyPred).
		Owns(&appsv1.Deployment{}, deleteOnlyPred).
		Owns(&appsv1.StatefulSet{}, deleteOnlyPred).
		Complete(r)
}
