// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package render

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/slipway/slipway/core/api/v1"
)

// Service builds a Kubernetes Service from App spec
func Service(app *v1.App) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace},
	}
	var ports []corev1.ServicePort
	for _, p := range app.Spec.Runtime.Ports {
		ports = append(ports, corev1.ServicePort{Name: p.Name, Port: p.ContainerPort})
	}
	svc.Spec.Selector = map[string]string{"app": app.Name}
	svc.Spec.Ports = ports
	return svc
}
