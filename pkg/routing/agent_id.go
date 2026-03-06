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
	"regexp"
	"strings"
)

const (
	DefaultAgentID   = "main"
	DefaultMainKey   = "main"
	DefaultAccountID = "default"
	MaxAgentIDLength = 64
)

var (
	validIDRe      = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)
	invalidCharsRe = regexp.MustCompile(`[^a-z0-9_-]+`)
	leadingDashRe  = regexp.MustCompile(`^-+`)
	trailingDashRe = regexp.MustCompile(`-+$`)
)

// NormalizeAgentID sanitizes an agent ID to [a-z0-9][a-z0-9_-]{0,63}.
// Invalid characters are collapsed to "-". Leading/trailing dashes stripped.
// Empty input returns DefaultAgentID ("main").
func NormalizeAgentID(id string) string {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return DefaultAgentID
	}
	lower := strings.ToLower(trimmed)
	if validIDRe.MatchString(lower) {
		return lower
	}
	result := invalidCharsRe.ReplaceAllString(lower, "-")
	result = leadingDashRe.ReplaceAllString(result, "")
	result = trailingDashRe.ReplaceAllString(result, "")
	if len(result) > MaxAgentIDLength {
		result = result[:MaxAgentIDLength]
	}
	if result == "" {
		return DefaultAgentID
	}
	return result
}

// NormalizeAccountID sanitizes an account ID. Empty returns DefaultAccountID.
func NormalizeAccountID(id string) string {
	trimmed := strings.TrimSpace(id)
	if trimmed == "" {
		return DefaultAccountID
	}
	lower := strings.ToLower(trimmed)
	if validIDRe.MatchString(lower) {
		return lower
	}
	result := invalidCharsRe.ReplaceAllString(lower, "-")
	result = leadingDashRe.ReplaceAllString(result, "")
	result = trailingDashRe.ReplaceAllString(result, "")
	if len(result) > MaxAgentIDLength {
		result = result[:MaxAgentIDLength]
	}
	if result == "" {
		return DefaultAccountID
	}
	return result
}
