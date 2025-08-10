// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package plugins

import "k8s.io/apimachinery/pkg/runtime"

// Manager controls plugin lifecycle; stub

type Manager struct {
	Scheme *runtime.Scheme
}

func NewManager(s *runtime.Scheme) *Manager { return &Manager{Scheme: s} }
