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

package migrate

import (
	"github.com/spf13/cobra"

	"github.com/Zhaoyikaiii/clotho/pkg/migrate"
)

func NewMigrateCommand() *cobra.Command {
	var opts migrate.Options

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate from xxxclaw(openclaw, etc.) to clotho",
		Args:  cobra.NoArgs,
		Example: `  clotho migrate
  clotho migrate --from openclaw
  clotho migrate --dry-run
  clotho migrate --refresh
  clotho migrate --force`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			m := migrate.NewMigrateInstance(opts)
			result, err := m.Run(opts)
			if err != nil {
				return err
			}
			if !opts.DryRun {
				m.PrintSummary(result)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.DryRun, "dry-run", false,
		"Show what would be migrated without making changes")
	cmd.Flags().StringVar(&opts.Source, "from", "openclaw",
		"Source to migrate from (e.g., openclaw)")
	cmd.Flags().BoolVar(&opts.Refresh, "refresh", false,
		"Re-sync workspace files from OpenClaw (repeatable)")
	cmd.Flags().BoolVar(&opts.ConfigOnly, "config-only", false,
		"Only migrate config, skip workspace files")
	cmd.Flags().BoolVar(&opts.WorkspaceOnly, "workspace-only", false,
		"Only migrate workspace files, skip config")
	cmd.Flags().BoolVar(&opts.Force, "force", false,
		"Skip confirmation prompts")
	cmd.Flags().StringVar(&opts.SourceHome, "source-home", "",
		"Override source home directory (default: ~/.openclaw)")
	cmd.Flags().StringVar(&opts.TargetHome, "target-home", "",
		"Override target home directory (default: ~/.clotho)")

	return cmd
}
