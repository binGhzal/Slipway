// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package build

// Builder defines an image build interface
type Builder interface {
	Build(spec Spec) (string, error)
}

type Spec struct {
	Repo       string
	Dockerfile string
	Context    string
}
