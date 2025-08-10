// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package plugins

import "context"

// Install chart via Helm - TODO
func (m *Manager) Install(ctx context.Context, name, oci string, values map[string]interface{}) error {
	// TODO implement helm install using SDK
	return nil
}
