// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Client returns a controller-runtime client
func Client() (client.Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	return client.New(cfg, client.Options{})
}

// Ping verifies connection
func Ping(c client.Client) error { return c.List(context.Background(), &corev1.NamespaceList{}) }
