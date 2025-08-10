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

type ReleaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ReleaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	rel := &v1.Release{}
	if err := r.Get(ctx, req.NamespacedName, rel); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	meta.SetStatusCondition(&rel.Status.Conditions, metav1.Condition{Type: "Observed", Status: metav1.ConditionTrue})
	if rel.Status.Phase == "" {
		rel.Status.Phase = "Pending"
	}
	if err := r.Status().Update(ctx, rel); err != nil {
		logger.Error(err, "update status")
	}
	return ctrl.Result{}, nil
}

func (r *ReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&v1.Release{}).Complete(r)
}
