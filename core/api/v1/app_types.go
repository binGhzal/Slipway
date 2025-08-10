// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image.ref"
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.url"
type App struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppSpec   `json:"spec,omitempty"`
	Status AppStatus `json:"status,omitempty"`
}

// AppSpec describes build and runtime
// +kubebuilder:validation:XPreserveUnknownFields=false
type AppSpec struct {
	Source      *AppSource      `json:"source,omitempty"`
	Image       *AppImage       `json:"image,omitempty"`
	Runtime     AppRuntime      `json:"runtime"`
	Route       *AppRoute       `json:"route,omitempty"`
	Autoscaling *AppAutoscaling `json:"autoscaling,omitempty"`
}

type AppSource struct {
	Repo   string    `json:"repo"`
	Branch string    `json:"branch,omitempty"`
	Build  *AppBuild `json:"build,omitempty"`
}

type AppBuild struct {
	// +kubebuilder:validation:Enum=buildkit;kaniko;buildpack;nixpacks
	// +kubebuilder:default=buildkit
	Strategy string `json:"strategy,omitempty"`
	// +kubebuilder:default=./Dockerfile
	Dockerfile string `json:"dockerfile,omitempty"`
	// +kubebuilder:default=.
	Context string `json:"context,omitempty"`
}

type AppImage struct {
	Ref string `json:"ref"`
}

type AppRuntime struct {
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	Replicas  int32                       `json:"replicas,omitempty"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	Env       []corev1.EnvVar             `json:"env,omitempty"`
	Secrets   []AppSecretRef              `json:"secrets,omitempty"`
	Ports     []corev1.ContainerPort      `json:"ports,omitempty"`
}

type AppSecretRef struct {
	Name string `json:"name"`
}

type AppRoute struct {
	Host string `json:"host"`
	// +kubebuilder:default=false
	TLS       bool    `json:"tls,omitempty"`
	ClassName *string `json:"className,omitempty"`
}

type AppAutoscaling struct {
	// +kubebuilder:default=false
	Enabled bool                  `json:"enabled,omitempty"`
	Min     int32                 `json:"min,omitempty"`
	Max     int32                 `json:"max,omitempty"`
	Policy  *AppAutoscalingPolicy `json:"policy,omitempty"`
}

type AppAutoscalingPolicy struct {
	// +kubebuilder:validation:Enum=cpu;memory
	// +kubebuilder:default=cpu
	Type string `json:"type,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	Target int32 `json:"target"`
}

// AppStatus provides app URL and release
// +kubebuilder:validation:XPreserveUnknownFields=false
type AppStatus struct {
	URL        string             `json:"url,omitempty"`
	Release    string             `json:"release,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []App `json:"items"`
}

func init() {
	SchemeBuilder.Register(&App{}, &AppList{})
}

func (in *App) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(App)
	*out = *in
	return out
}
func (in *AppList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(AppList)
	*out = *in
	return out
}
