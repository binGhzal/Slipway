// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package render

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/slipway/slipway/core/api/v1"
)

// Deployment renders a Deployment for the App
func Deployment(app *v1.App, image string) *appsv1.Deployment {
	replicas := int32(1)
	if app.Spec.Runtime.Replicas > 0 {
		replicas = app.Spec.Runtime.Replicas
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": app.Name}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": app.Name}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: app.Name, Image: image, Ports: app.Spec.Runtime.Ports, Env: app.Spec.Runtime.Env}}},
			},
		},
	}
	return dep
}
