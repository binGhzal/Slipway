// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package gitops

// Writer defines interface for emitting manifests to Git. Not implemented.

type Writer interface {
	Write(objects []byte) error
}
