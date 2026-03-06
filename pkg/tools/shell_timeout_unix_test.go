//go:build !windows

package tools

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

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func processExists(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil || err == syscall.EPERM
}

func TestShellTool_TimeoutKillsChildProcess(t *testing.T) {
	tool, err := NewExecTool(t.TempDir(), false)
	if err != nil {
		t.Errorf("unable to configure exec tool: %s", err)
	}

	tool.SetTimeout(500 * time.Millisecond)

	args := map[string]any{
		// Spawn a child process that would outlive the shell unless process-group kill is used.
		"command": "sleep 60 & echo $! > child.pid; wait",
	}

	result := tool.Execute(context.Background(), args)
	if !result.IsError {
		t.Fatalf("expected timeout error, got success: %s", result.ForLLM)
	}
	if !strings.Contains(result.ForLLM, "timed out") {
		t.Fatalf("expected timeout message, got: %s", result.ForLLM)
	}

	childPIDPath := filepath.Join(tool.workingDir, "child.pid")
	data, err := os.ReadFile(childPIDPath)
	if err != nil {
		t.Fatalf("failed to read child pid file: %v", err)
	}

	childPID, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("failed to parse child pid: %v", err)
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if !processExists(childPID) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	t.Fatalf("child process %d is still running after timeout", childPID)
}
