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

	"github.com/spf13/cobra"

	"github.com/Zhaoyikaiii/clotho/pkg/cron"
)

func newAddCommand(storePath func() string) *cobra.Command {
	var (
		name    string
		message string
		every   int64
		cronExp string
		deliver bool
		channel string
		to      string
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new scheduled job",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if every <= 0 && cronExp == "" {
				return fmt.Errorf("either --every or --cron must be specified")
			}

			var schedule cron.CronSchedule
			if every > 0 {
				everyMS := every * 1000
				schedule = cron.CronSchedule{Kind: "every", EveryMS: &everyMS}
			} else {
				schedule = cron.CronSchedule{Kind: "cron", Expr: cronExp}
			}

			cs := cron.NewCronService(storePath(), nil)
			job, err := cs.AddJob(name, schedule, message, deliver, channel, to)
			if err != nil {
				return fmt.Errorf("error adding job: %w", err)
			}

			fmt.Printf("✓ Added job '%s' (%s)\n", job.Name, job.ID)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Job name")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Message for agent")
	cmd.Flags().Int64VarP(&every, "every", "e", 0, "Run every N seconds")
	cmd.Flags().StringVarP(&cronExp, "cron", "c", "", "Cron expression (e.g. '0 9 * * *')")
	cmd.Flags().BoolVarP(&deliver, "deliver", "d", false, "Deliver response to channel")
	cmd.Flags().StringVar(&to, "to", "", "Recipient for delivery")
	cmd.Flags().StringVar(&channel, "channel", "", "Channel for delivery")

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("message")
	cmd.MarkFlagsMutuallyExclusive("every", "cron")

	return cmd
}
