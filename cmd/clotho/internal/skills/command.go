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

package skills

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Zhaoyikaiii/clotho/cmd/clotho/internal"
	"github.com/Zhaoyikaiii/clotho/pkg/skills"
)

type deps struct {
	workspace    string
	installer    *skills.SkillInstaller
	skillsLoader *skills.SkillsLoader
}

func NewSkillsCommand() *cobra.Command {
	var d deps

	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Manage skills",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := internal.LoadConfig()
			if err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}

			d.workspace = cfg.WorkspacePath()
			d.installer = skills.NewSkillInstaller(d.workspace)

			// get global config directory and builtin skills directory
			globalDir := filepath.Dir(internal.GetConfigPath())
			globalSkillsDir := filepath.Join(globalDir, "skills")
			builtinSkillsDir := filepath.Join(globalDir, "clotho", "skills")
			d.skillsLoader = skills.NewSkillsLoader(d.workspace, globalSkillsDir, builtinSkillsDir)

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	installerFn := func() (*skills.SkillInstaller, error) {
		if d.installer == nil {
			return nil, fmt.Errorf("skills installer is not initialized")
		}
		return d.installer, nil
	}

	loaderFn := func() (*skills.SkillsLoader, error) {
		if d.skillsLoader == nil {
			return nil, fmt.Errorf("skills loader is not initialized")
		}
		return d.skillsLoader, nil
	}

	workspaceFn := func() (string, error) {
		if d.workspace == "" {
			return "", fmt.Errorf("workspace is not initialized")
		}
		return d.workspace, nil
	}

	cmd.AddCommand(
		newListCommand(loaderFn),
		newInstallCommand(installerFn),
		newInstallBuiltinCommand(workspaceFn),
		newListBuiltinCommand(),
		newRemoveCommand(installerFn),
		newSearchCommand(),
		newShowCommand(loaderFn),
	)

	return cmd
}
