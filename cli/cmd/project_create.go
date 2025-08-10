// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var projectCreateCmd = &cobra.Command{
	Use:   "project create <name>",
	Args:  cobra.ExactArgs(1),
	Short: "Create a Project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Project %s created (stub)\n", args[0])
		return nil
	},
}

func init() { rootCmd.AddCommand(projectCreateCmd) }
