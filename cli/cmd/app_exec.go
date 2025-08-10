// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appExecCmd = &cobra.Command{
	Use:   "app exec <name> -- cmd",
	Args:  cobra.MinimumNArgs(1),
	Short: "Exec into app container",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Exec into %s (stub)\n", args[0])
		return nil
	},
}

func init() { rootCmd.AddCommand(appExecCmd) }
