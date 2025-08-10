// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=proj
// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".status.namespace"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

// ProjectSpec defines desired state
// +kubebuilder:validation:XPreserveUnknownFields=false
type ProjectSpec struct {
	Owners []string `json:"owners,omitempty"`

	Quotas *ProjectQuotas `json:"quotas,omitempty"`

	// +kubebuilder:validation:Enum=defaultDeny;allowAll
	// +kubebuilder:default=defaultDeny
	NetworkPolicy string `json:"networkPolicy,omitempty"`

	Domains ProjectDomains `json:"domains"`
}

// ProjectQuotas for cpu/memory
// +kubebuilder:validation:XPreserveUnknownFields=false
type ProjectQuotas struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// ProjectDomains defines domain settings
type ProjectDomains struct {
	Base string `json:"base"`
}

// ProjectStatus holds observed state
type ProjectStatus struct {
	Namespace  string             `json:"namespace,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// ProjectList contains a list of Projects
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Project{}, &ProjectList{})
}

func (in *Project) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Project)
	*out = *in
	return out
}
func (in *ProjectList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(ProjectList)
	*out = *in
	return out
}
