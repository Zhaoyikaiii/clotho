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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstallSubcommand(t *testing.T) {
	cmd := newInstallCommand(nil)

	require.NotNil(t, cmd)

	assert.Equal(t, "install", cmd.Use)
	assert.Equal(t, "Install skill from GitHub", cmd.Short)

	assert.Nil(t, cmd.Run)
	assert.NotNil(t, cmd.RunE)

	assert.True(t, cmd.HasExample())
	assert.False(t, cmd.HasSubCommands())

	assert.True(t, cmd.HasFlags())
	assert.NotNil(t, cmd.Flags().Lookup("registry"))

	assert.Len(t, cmd.Aliases, 0)
}

func TestInstallCommandArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		registry    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no registry, one arg",
			args:        []string{"sipeed/clotho-skills/weather"},
			registry:    "",
			expectError: false,
		},
		{
			name:        "no registry, no args",
			args:        []string{},
			registry:    "",
			expectError: true,
			errorMsg:    "exactly 1 argument is required: <github>",
		},
		{
			name:        "no registry, too many args",
			args:        []string{"arg1", "arg2"},
			registry:    "",
			expectError: true,
			errorMsg:    "exactly 1 argument is required: <github>",
		},
		{
			name:        "with registry, one arg",
			args:        []string{"weather-skill"},
			registry:    "clawhub",
			expectError: false,
		},
		{
			name:        "with registry, no args",
			args:        []string{},
			registry:    "clawhub",
			expectError: true,
			errorMsg:    "when --registry is set, exactly 1 argument is required: <slug>",
		},
		{
			name:        "with registry, too many args",
			args:        []string{"arg1", "arg2"},
			registry:    "clawhub",
			expectError: true,
			errorMsg:    "when --registry is set, exactly 1 argument is required: <slug>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newInstallCommand(nil)

			if tt.registry != "" {
				require.NoError(t, cmd.Flags().Set("registry", tt.registry))
			}

			err := cmd.Args(cmd, tt.args)
			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
