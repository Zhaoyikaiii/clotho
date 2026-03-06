// Copyright (c) 2026 Clotho contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cron

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal"
)

func NewCronCommand() *cobra.Command {
	var storePath string

	cmd := &cobra.Command{
		Use:     "cron",
		Aliases: []string{"c"},
		Short:   "Manage scheduled tasks",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
		// Resolve storePath at execution time so it reflects the current config
		// and is shared across all subcommands.
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := internal.LoadConfig()
			if err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}
			storePath = filepath.Join(cfg.WorkspacePath(), "cron", "jobs.json")
			return nil
		},
	}

	cmd.AddCommand(
		newListCommand(func() string { return storePath }),
		newAddCommand(func() string { return storePath }),
		newRemoveCommand(func() string { return storePath }),
		newEnableCommand(func() string { return storePath }),
		newDisableCommand(func() string { return storePath }),
	)

	return cmd
}
