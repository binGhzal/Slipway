// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package controller

import (
	context "context"

	v1 "github.com/slipway/slipway/core/api/v1"
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type PluginReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *PluginReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	pl := &v1.Plugin{}
	if err := r.Get(ctx, req.NamespacedName, pl); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODO: integrate Helm install
	ready := true
	pl.Status.Ready = ready
	var caps []string
	for _, c := range pl.Spec.Provides {
		caps = append(caps, c.Capability)
	}
	pl.Status.Capabilities = caps
	meta.SetStatusCondition(&pl.Status.Conditions, metav1.Condition{Type: "Ready", Status: metav1.ConditionTrue, Reason: "Stub"})
	if err := r.Status().Update(ctx, pl); err != nil {
		logger.Error(err, "update status")
	}
	return ctrl.Result{}, nil
}

func (r *PluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&v1.Plugin{}).Complete(r)
}
