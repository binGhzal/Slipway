// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package controller

import (
	context "context"
	fmt "fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/slipway/slipway/core/api/v1"
	"github.com/slipway/slipway/core/pkg/render"
)

type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	app := &v1.App{}
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	image := "busybox"
	if app.Spec.Image != nil {
		image = app.Spec.Image.Ref
	}

	dep := render.Deployment(app, image)
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, dep, func() error {
		dep.Spec = render.Deployment(app, image).Spec
		return controllerutil.SetControllerReference(app, dep, r.Scheme)
	}); err != nil {
		return ctrl.Result{}, err
	}

	if len(app.Spec.Runtime.Ports) > 0 {
		svc := render.Service(app)
		if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
			svc.Spec = render.Service(app).Spec
			return controllerutil.SetControllerReference(app, svc, r.Scheme)
		}); err != nil {
			return ctrl.Result{}, err
		}
	}

	if app.Spec.Route != nil {
		ing := render.Ingress(app)
		if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, ing, func() error {
			*ing = *render.Ingress(app)
			return controllerutil.SetControllerReference(app, ing, r.Scheme)
		}); err != nil {
			return ctrl.Result{}, err
		}
		scheme := "http"
		if app.Spec.Route.TLS {
			scheme = "https"
		}
		app.Status.URL = fmt.Sprintf("%s://%s", scheme, app.Spec.Route.Host)
	}

	app.Status.Release = ""
	if err := r.Status().Update(ctx, app); err != nil {
		logger.Error(err, "status update")
	}
	return ctrl.Result{}, nil
}

func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&v1.App{}).Owns(&appsv1.Deployment{}).Owns(&corev1.Service{}).Owns(&networkingv1.Ingress{}).Complete(r)
}
