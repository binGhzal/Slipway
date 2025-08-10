// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package capability

import (
	context "context"

	v1 "github.com/slipway/slipway/core/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Has returns whether a capability is available and any details
func Has(c client.Client, ns string, cap string) (bool, map[string]string, error) {
	ctx := context.Background()
	var list v1.PluginList
	if err := c.List(ctx, &list); err != nil {
		return false, nil, err
	}
	for _, p := range list.Items {
		if !p.Status.Ready {
			continue
		}
		for _, prov := range p.Spec.Provides {
			if prov.Capability == cap {
				return true, prov.Details, nil
			}
		}
	}
	return false, nil, nil
}
