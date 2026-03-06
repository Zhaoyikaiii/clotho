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

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfigPath(t *testing.T) {
	t.Setenv("HOME", "/tmp/home")

	got := GetConfigPath()
	want := filepath.Join("/tmp/home", ".clotho", "config.json")

	assert.Equal(t, want, got)
}

func TestFormatVersion_NoGitCommit(t *testing.T) {
	oldVersion, oldGit := version, gitCommit
	t.Cleanup(func() { version, gitCommit = oldVersion, oldGit })

	version = "1.2.3"
	gitCommit = ""

	assert.Equal(t, "1.2.3", FormatVersion())
}

func TestFormatVersion_WithGitCommit(t *testing.T) {
	oldVersion, oldGit := version, gitCommit
	t.Cleanup(func() { version, gitCommit = oldVersion, oldGit })

	version = "1.2.3"
	gitCommit = "abc123"

	assert.Equal(t, "1.2.3 (git: abc123)", FormatVersion())
}

func TestFormatBuildInfo_UsesBuildTimeAndGoVersion_WhenSet(t *testing.T) {
	oldBuildTime, oldGoVersion := buildTime, goVersion
	t.Cleanup(func() { buildTime, goVersion = oldBuildTime, oldGoVersion })

	buildTime = "2026-02-20T00:00:00Z"
	goVersion = "go1.23.0"

	build, goVer := FormatBuildInfo()

	assert.Equal(t, buildTime, build)
	assert.Equal(t, goVersion, goVer)
}

func TestFormatBuildInfo_EmptyBuildTime_ReturnsEmptyBuild(t *testing.T) {
	oldBuildTime, oldGoVersion := buildTime, goVersion
	t.Cleanup(func() { buildTime, goVersion = oldBuildTime, oldGoVersion })

	buildTime = ""
	goVersion = "go1.23.0"

	build, goVer := FormatBuildInfo()

	assert.Empty(t, build)
	assert.Equal(t, goVersion, goVer)
}

func TestFormatBuildInfo_EmptyGoVersion_FallsBackToRuntimeVersion(t *testing.T) {
	oldBuildTime, oldGoVersion := buildTime, goVersion
	t.Cleanup(func() { buildTime, goVersion = oldBuildTime, oldGoVersion })

	buildTime = "x"
	goVersion = ""

	build, goVer := FormatBuildInfo()

	assert.Equal(t, "x", build)
	assert.Equal(t, runtime.Version(), goVer)
}

func TestGetConfigPath_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-specific HOME behavior varies; run on windows")
	}

	testUserProfilePath := `C:\Users\Test`
	t.Setenv("USERPROFILE", testUserProfilePath)

	got := GetConfigPath()
	want := filepath.Join(testUserProfilePath, ".clotho", "config.json")

	require.True(t, strings.EqualFold(got, want), "GetConfigPath() = %q, want %q", got, want)
}

func TestGetVersion(t *testing.T) {
	assert.Equal(t, "dev", GetVersion())
}

func TestGetConfigPath_WithEnv(t *testing.T) {
	t.Setenv("CLOTHO_CONFIG", "/tmp/custom/config.json")
	t.Setenv("HOME", "/tmp/home") // Also set home to ensure env is preferred

	got := GetConfigPath()
	want := "/tmp/custom/config.json"

	assert.Equal(t, want, got)
}
