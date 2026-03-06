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

package utils

import (
	"fmt"
	"strings"
)

// ValidateSkillIdentifier validates that the given skill identifier (slug or registry name) is non-empty
// and does not contain path separators ("/", "\\") or ".." for security.
func ValidateSkillIdentifier(identifier string) error {
	trimmed := strings.TrimSpace(identifier)
	if trimmed == "" {
		return fmt.Errorf("identifier is required and must be a non-empty string")
	}
	if strings.ContainsAny(trimmed, "/\\") || strings.Contains(trimmed, "..") {
		return fmt.Errorf("identifier must not contain path separators or '..' to prevent directory traversal")
	}
	return nil
}
