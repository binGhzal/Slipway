// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready"
type Plugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PluginSpec   `json:"spec,omitempty"`
	Status PluginStatus `json:"status,omitempty"`
}

// PluginSpec describes plugin source and capabilities
type PluginSpec struct {
	Source   PluginSource         `json:"source"`
	Values   runtime.RawExtension `json:"values,omitempty"`
	Provides []PluginCapability   `json:"provides,omitempty"`
	// +kubebuilder:default=slipway-system
	Namespace string `json:"namespace,omitempty"`
}

type PluginSource struct {
	OCI     string `json:"oci,omitempty"`
	Repo    string `json:"repo,omitempty"`
	Chart   string `json:"chart,omitempty"`
	Version string `json:"version,omitempty"`
}

type PluginCapability struct {
	// +kubebuilder:validation:Enum=ingress;certs;dns;autoscaling;secrets;observability;mesh;serverless;db;gitops
	Capability string            `json:"capability"`
	Details    map[string]string `json:"details,omitempty"`
}

// PluginStatus shows readiness
// +kubebuilder:validation:XPreserveUnknownFields=false
type PluginStatus struct {
	Ready        bool               `json:"ready,omitempty"`
	Resources    int32              `json:"resources,omitempty"`
	Capabilities []string           `json:"capabilities,omitempty"`
	Conditions   []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
type PluginList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Plugin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Plugin{}, &PluginList{})
}

func (in *Plugin) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Plugin)
	*out = *in
	return out
}
func (in *PluginList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(PluginList)
	*out = *in
	return out
}
