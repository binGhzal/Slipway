// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package controller

import (
	context "context"

	v1 "github.com/slipway/slipway/core/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ProjectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	proj := &v1.Project{}
	if err := r.Get(ctx, req.NamespacedName, proj); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	ns := &corev1.Namespace{}
	if err := r.Get(ctx, types.NamespacedName{Name: proj.Name}, ns); err != nil {
		if errors.IsNotFound(err) {
			ns.Name = proj.Name
			ns.Labels = map[string]string{"app.kubernetes.io/managed-by": "slipway"}
			if err := r.Create(ctx, ns); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{}, err
		}
	}

	proj.Status.Namespace = proj.Name
	meta.SetStatusCondition(&proj.Status.Conditions, metav1.Condition{Type: "Ready", Status: metav1.ConditionTrue, Reason: "Reconciled"})
	if err := r.Status().Update(ctx, proj); err != nil {
		logger.Error(err, "update status")
	}
	return ctrl.Result{}, nil
}

func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&v1.Project{}).Complete(r)
}
