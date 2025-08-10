// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package format

import "fmt"

// Human table stub
func Table(headers []string, rows [][]string) {
	fmt.Println(headers)
	for _, r := range rows {
		fmt.Println(r)
	}
}
