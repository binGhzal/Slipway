// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package render

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/slipway/slipway/core/api/v1"
)

// Ingress returns an ingress if route is specified
func Ingress(app *v1.App) *networkingv1.Ingress {
	if app.Spec.Route == nil {
		return nil
	}
	ing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{{
				Host: app.Spec.Route.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{Paths: []networkingv1.HTTPIngressPath{{
						Path:     "/",
						PathType: pathTypePtr(networkingv1.PathTypePrefix),
						Backend: networkingv1.IngressBackend{Service: &networkingv1.IngressServiceBackend{
							Name: app.Name,
							Port: networkingv1.ServiceBackendPort{Number: app.Spec.Runtime.Ports[0].ContainerPort},
						}},
					}}},
				},
			}},
		},
	}
	if app.Spec.Route.TLS {
		ing.Spec.TLS = []networkingv1.IngressTLS{{Hosts: []string{app.Spec.Route.Host}, SecretName: app.Name + "-tls"}}
	}
	if app.Spec.Route.ClassName != nil {
		ing.Spec.IngressClassName = app.Spec.Route.ClassName
	}
	return ing
}

func pathTypePtr(p networkingv1.PathType) *networkingv1.PathType { return &p }
