// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Slipway core",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Initializing Slipway (stub)")
		return nil
	},
}

func init() { rootCmd.AddCommand(initCmd) }
