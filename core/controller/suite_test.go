// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package controller

import (
	testing "testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = Describe("stub", func() {
	It("runs", func() {
		Expect(true).To(BeTrue())
	})
})
