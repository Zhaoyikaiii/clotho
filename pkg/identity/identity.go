// Package identity provides unified user identity utilities for Clotho.
// It introduces a canonical "platform:id" format and matching logic
// that is backward-compatible with all legacy allow-list formats.
package identity

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
	"strings"

	"github.com/Zhaoyikaiii/clotho/pkg/bus"
)

// BuildCanonicalID constructs a canonical "platform:id" identifier.
// Both platform and platformID are lowercased and trimmed.
func BuildCanonicalID(platform, platformID string) string {
	p := strings.ToLower(strings.TrimSpace(platform))
	id := strings.TrimSpace(platformID)
	if p == "" || id == "" {
		return ""
	}
	return p + ":" + id
}

// ParseCanonicalID splits a canonical ID ("platform:id") into its parts.
// Returns ok=false if the input does not contain a colon separator.
func ParseCanonicalID(canonical string) (platform, id string, ok bool) {
	canonical = strings.TrimSpace(canonical)
	idx := strings.Index(canonical, ":")
	if idx <= 0 || idx == len(canonical)-1 {
		return "", "", false
	}
	return canonical[:idx], canonical[idx+1:], true
}

// MatchAllowed checks whether the given sender matches a single allow-list entry.
// It is backward-compatible with all legacy formats:
//
//   - "123456"              → matches sender.PlatformID
//   - "@alice"              → matches sender.Username
//   - "123456|alice"        → matches PlatformID or Username
//   - "telegram:123456"     → exact match on sender.CanonicalID
func MatchAllowed(sender bus.SenderInfo, allowed string) bool {
	allowed = strings.TrimSpace(allowed)
	if allowed == "" {
		return false
	}

	// Try canonical match first: "platform:id" format
	if platform, id, ok := ParseCanonicalID(allowed); ok {
		// Only treat as canonical if the platform portion looks like a known platform name
		// (not a pure-numeric string, which could be a compound ID)
		if !isNumeric(platform) {
			candidate := BuildCanonicalID(platform, id)
			if candidate != "" && sender.CanonicalID != "" {
				return strings.EqualFold(sender.CanonicalID, candidate)
			}
			// If sender has no canonical ID, try matching platform + platformID
			return strings.EqualFold(platform, sender.Platform) &&
				sender.PlatformID == id
		}
	}

	// Strip leading "@" for username matching
	trimmed := strings.TrimPrefix(allowed, "@")

	// Split compound "id|username" format
	allowedID := trimmed
	allowedUser := ""
	if idx := strings.Index(trimmed, "|"); idx > 0 {
		allowedID = trimmed[:idx]
		allowedUser = trimmed[idx+1:]
	}

	// Match against PlatformID
	if sender.PlatformID != "" && sender.PlatformID == allowedID {
		return true
	}

	// Match against Username
	if sender.Username != "" {
		if sender.Username == trimmed || sender.Username == allowedUser {
			return true
		}
	}

	// Match compound sender format against allowed parts
	if allowedUser != "" && sender.PlatformID != "" && sender.PlatformID == allowedID {
		return true
	}
	if allowedUser != "" && sender.Username != "" && sender.Username == allowedUser {
		return true
	}

	return false
}

// isNumeric returns true if s consists entirely of digits.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
