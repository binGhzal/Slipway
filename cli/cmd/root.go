// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slipway",
	Short: "Slipway CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
