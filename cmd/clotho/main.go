// Clotho - Distributed Task Orchestration & Workflow Engine
// License: MIT
//
// Copyright (c) 2026 Clotho contributors

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/agent"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/auth"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/cron"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/gateway"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/migrate"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/onboard"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/skills"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/status"
	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal/version"
)

func NewClothoCommand() *cobra.Command {
	short := fmt.Sprintf("%s clotho - Distributed Task Orchestration v%s\n\n", internal.Logo, internal.GetVersion())

	cmd := &cobra.Command{
		Use:     "clotho",
		Short:   short,
		Example: "clotho list",
	}

	cmd.AddCommand(
		onboard.NewOnboardCommand(),
		agent.NewAgentCommand(),
		auth.NewAuthCommand(),
		gateway.NewGatewayCommand(),
		status.NewStatusCommand(),
		cron.NewCronCommand(),
		migrate.NewMigrateCommand(),
		skills.NewSkillsCommand(),
		version.NewVersionCommand(),
	)

	return cmd
}

const (
	colorBlue = "\033[1;38;2;62;93;185m"
	colorRed  = "\033[1;38;2;213;70;70m"
	banner    = "\r\n" +
		colorBlue + "██████╗ ██╗ ██████╗ ██████╗ " + colorRed + " ██████╗██╗      █████╗ ██╗    ██╗\n" +
		colorBlue + "██╔══██╗██║██╔════╝██╔═══██╗" + colorRed + "██╔════╝██║     ██╔══██╗██║    ██║\n" +
		colorBlue + "██████╔╝██║██║     ██║   ██║" + colorRed + "██║     ██║     ███████║██║ █╗ ██║\n" +
		colorBlue + "██╔═══╝ ██║██║     ██║   ██║" + colorRed + "██║     ██║     ██╔══██║██║███╗██║\n" +
		colorBlue + "██║     ██║╚██████╗╚██████╔╝" + colorRed + "╚██████╗███████╗██║  ██║╚███╔███╔╝\n" +
		colorBlue + "╚═╝     ╚═╝ ╚═════╝ ╚═════╝ " + colorRed + " ╚═════╝╚══════╝╚═╝  ╚═╝ ╚══╝╚══╝\n " +
		"\033[0m\r\n"
)

func main() {
	fmt.Printf("%s", banner)
	cmd := NewClothoCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
