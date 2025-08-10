// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package render

import (
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/slipway/slipway/core/api/v1"
)

// HPA renders a HorizontalPodAutoscaler if enabled
func HPA(app *v1.App) *autoscalingv2.HorizontalPodAutoscaler {
	if app.Spec.Autoscaling == nil || !app.Spec.Autoscaling.Enabled {
		return nil
	}
	h := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{APIVersion: "apps/v1", Kind: "Deployment", Name: app.Name},
			MinReplicas:    &app.Spec.Autoscaling.Min,
			MaxReplicas:    app.Spec.Autoscaling.Max,
		},
	}
	return h
}
