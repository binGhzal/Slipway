// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2025 Slipway Authors

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appApplyCmd = &cobra.Command{
	Use:   "app apply -f file.yaml",
	Short: "Apply an App spec",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Applying app (stub)")
		return nil
	},
}

func init() { rootCmd.AddCommand(appApplyCmd) }
