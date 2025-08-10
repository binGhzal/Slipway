// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appLogsCmd = &cobra.Command{
	Use:   "app logs <name>",
	Args:  cobra.ExactArgs(1),
	Short: "Stream app logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Streaming logs for %s (stub)\n", args[0])
		return nil
	},
}

func init() { rootCmd.AddCommand(appLogsCmd) }
