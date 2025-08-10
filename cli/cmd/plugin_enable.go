// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pluginEnableCmd = &cobra.Command{
	Use:   "plugin enable <name>",
	Args:  cobra.ExactArgs(1),
	Short: "Enable a plugin from catalog",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Enabling plugin %s (stub)\n", args[0])
		return nil
	},
}

func init() { rootCmd.AddCommand(pluginEnableCmd) }
