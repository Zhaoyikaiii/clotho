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

package openclaw

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zhaoyikaiii/clotho/pkg/config"
	"github.com/Zhaoyikaiii/clotho/pkg/migrate/internal"
)

var providerMapping = map[string]string{
	"anthropic":  "anthropic",
	"claude":     "anthropic",
	"openai":     "openai",
	"gpt":        "openai",
	"groq":       "groq",
	"ollama":     "ollama",
	"openrouter": "openrouter",
	"deepseek":   "deepseek",
	"together":   "together",
	"mistral":    "mistral",
	"fireworks":  "fireworks",
	"google":     "google",
	"gemini":     "google",
	"xai":        "xai",
	"grok":       "xai",
	"cerebras":   "cerebras",
	"sambanova":  "sambanova",
}

type OpenclawHandler struct {
	opts             Options
	sourceConfigFile string
	sourceWorkspace  string
}

type (
	Options   = internal.Options
	Action    = internal.Action
	Result    = internal.Result
	Operation = internal.Operation
)

func NewOpenclawHandler(opts Options) (Operation, error) {
	home, err := resolveSourceHome(opts.SourceHome)
	if err != nil {
		return nil, err
	}
	opts.SourceHome = home

	configFile, err := findSourceConfig(home)
	if err != nil {
		return nil, err
	}
	return &OpenclawHandler{
		opts:             opts,
		sourceWorkspace:  filepath.Join(opts.SourceHome, "workspace"),
		sourceConfigFile: configFile,
	}, nil
}

func (o *OpenclawHandler) GetSourceName() string {
	return "openclaw"
}

func (o *OpenclawHandler) GetSourceHome() (string, error) {
	return o.opts.SourceHome, nil
}

func (o *OpenclawHandler) GetSourceWorkspace() (string, error) {
	return o.sourceWorkspace, nil
}

func (o *OpenclawHandler) GetSourceConfigFile() (string, error) {
	return o.sourceConfigFile, nil
}

func (o *OpenclawHandler) GetMigrateableFiles() []string {
	return migrateableFiles
}

func (o *OpenclawHandler) GetMigrateableDirs() []string {
	return migrateableDirs
}

func (o *OpenclawHandler) ExecuteConfigMigration(srcConfigPath, dstConfigPath string) error {
	openclawCfg, err := LoadOpenClawConfig(srcConfigPath)
	if err != nil {
		return err
	}

	Cfg, warnings, err := openclawCfg.ConvertToClotho(o.opts.SourceHome)
	if err != nil {
		return err
	}

	for _, w := range warnings {
		fmt.Printf("  Warning: %s\n", w)
	}

	incoming := Cfg.ToStandardConfig()
	if err := os.MkdirAll(filepath.Dir(dstConfigPath), 0o755); err != nil {
		return err
	}

	return config.SaveConfig(dstConfigPath, incoming)
}

func resolveSourceHome(override string) (string, error) {
	if override != "" {
		return internal.ExpandHome(override), nil
	}
	if envHome := os.Getenv("OPENCLAW_HOME"); envHome != "" {
		return internal.ExpandHome(envHome), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	return filepath.Join(home, ".openclaw"), nil
}

func findSourceConfig(sourceHome string) (string, error) {
	candidates := []string{
		filepath.Join(sourceHome, "openclaw.json"),
		filepath.Join(sourceHome, "config.json"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("no config file found in %s (tried openclaw.json, config.json)", sourceHome)
}

func rewriteWorkspacePath(path string) string {
	path = strings.Replace(path, ".openclaw", ".clotho", 1)
	return path
}

func mapProvider(provider string) string {
	if mapped, ok := providerMapping[strings.ToLower(provider)]; ok {
		return mapped
	}
	return strings.ToLower(provider)
}
