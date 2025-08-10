// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Release struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReleaseSpec   `json:"spec,omitempty"`
	Status ReleaseStatus `json:"status,omitempty"`
}

type ReleaseSpec struct {
	AppRef     NamespacedName `json:"appRef"`
	Image      string         `json:"image"`
	ConfigHash string         `json:"configHash"`
	RenderedAt metav1.Time    `json:"renderedAt"`
}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type ReleaseStatus struct {
	// +kubebuilder:validation:Enum=Pending;Deployed;Failed
	Phase      string             `json:"phase,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
type ReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Release `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Release{}, &ReleaseList{})
}

func (in *Release) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Release)
	*out = *in
	return out
}
func (in *ReleaseList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(ReleaseList)
	*out = *in
	return out
}
