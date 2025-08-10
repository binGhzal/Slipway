// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrlscheme "sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	GroupVersion  = schema.GroupVersion{Group: "paas.slipway.dev", Version: "v1"}
	SchemeBuilder = &ctrlscheme.Builder{GroupVersion: GroupVersion}
	AddToScheme   = SchemeBuilder.AddToScheme
)
