// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Slipway status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Status (stub)")
		return nil
	},
}

func init() { rootCmd.AddCommand(statusCmd) }
