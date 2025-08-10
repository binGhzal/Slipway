// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package build

import "fmt"

// BuildKitBuilder is a stub implementation using BuildKit
// TODO: implement real build job

type BuildKitBuilder struct{}

func (b *BuildKitBuilder) Build(spec Spec) (string, error) {
	// In MVP return fake digest
	return "sha256:DEADBEEF", fmt.Errorf("buildkit not implemented")
}
