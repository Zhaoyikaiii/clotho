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

package routing

import (
	"strings"
	"testing"
)

func TestNormalizeAgentID_Empty(t *testing.T) {
	if got := NormalizeAgentID(""); got != DefaultAgentID {
		t.Errorf("NormalizeAgentID('') = %q, want %q", got, DefaultAgentID)
	}
}

func TestNormalizeAgentID_Whitespace(t *testing.T) {
	if got := NormalizeAgentID("  "); got != DefaultAgentID {
		t.Errorf("NormalizeAgentID('  ') = %q, want %q", got, DefaultAgentID)
	}
}

func TestNormalizeAgentID_Valid(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"main", "main"},
		{"Main", "main"},
		{"SALES", "sales"},
		{"support-bot", "support-bot"},
		{"agent_1", "agent_1"},
		{"a", "a"},
		{"0test", "0test"},
	}
	for _, tt := range tests {
		if got := NormalizeAgentID(tt.input); got != tt.want {
			t.Errorf("NormalizeAgentID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeAgentID_InvalidChars(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Hello World", "hello-world"},
		{"agent@123", "agent-123"},
		{"foo.bar.baz", "foo-bar-baz"},
		{"--leading", "leading"},
		{"--both--", "both"},
	}
	for _, tt := range tests {
		if got := NormalizeAgentID(tt.input); got != tt.want {
			t.Errorf("NormalizeAgentID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeAgentID_AllInvalid(t *testing.T) {
	if got := NormalizeAgentID("@@@"); got != DefaultAgentID {
		t.Errorf("NormalizeAgentID('@@@') = %q, want %q", got, DefaultAgentID)
	}
}

func TestNormalizeAgentID_TruncatesAt64(t *testing.T) {
	var long strings.Builder
	for range 100 {
		long.WriteString("a")
	}
	got := NormalizeAgentID(long.String())
	if len(got) > MaxAgentIDLength {
		t.Errorf("length = %d, want <= %d", len(got), MaxAgentIDLength)
	}
}

func TestNormalizeAccountID_Empty(t *testing.T) {
	if got := NormalizeAccountID(""); got != DefaultAccountID {
		t.Errorf("NormalizeAccountID('') = %q, want %q", got, DefaultAccountID)
	}
}

func TestNormalizeAccountID_Valid(t *testing.T) {
	if got := NormalizeAccountID("MyBot"); got != "mybot" {
		t.Errorf("NormalizeAccountID('MyBot') = %q, want 'mybot'", got)
	}
}

func TestNormalizeAccountID_InvalidChars(t *testing.T) {
	if got := NormalizeAccountID("bot@home"); got != "bot-home" {
		t.Errorf("NormalizeAccountID('bot@home') = %q, want 'bot-home'", got)
	}
}
