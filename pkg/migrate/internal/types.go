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

package internal

type Options struct {
	DryRun        bool
	ConfigOnly    bool
	WorkspaceOnly bool
	Force         bool
	Refresh       bool
	Source        string
	SourceHome    string
	TargetHome    string
}

type Operation interface {
	GetSourceName() string
	GetSourceHome() (string, error)
	GetSourceWorkspace() (string, error)
	GetSourceConfigFile() (string, error)
	ExecuteConfigMigration(srcConfigPath, dstConfigPath string) error
	GetMigrateableFiles() []string
	GetMigrateableDirs() []string
}

type HandlerFactory func(opts Options) Operation

type ActionType int

const (
	ActionCopy ActionType = iota
	ActionSkip
	ActionBackup
	ActionConvertConfig
	ActionCreateDir
	ActionMergeConfig
)

type Action struct {
	Type        ActionType
	Source      string
	Target      string
	Description string
}

type Result struct {
	FilesCopied    int
	FilesSkipped   int
	BackupsCreated int
	ConfigMigrated bool
	DirsCreated    int
	Warnings       []string
	Errors         []error
}
