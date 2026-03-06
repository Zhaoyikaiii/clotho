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

package channels

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorsIs(t *testing.T) {
	wrapped := fmt.Errorf("telegram API: %w", ErrRateLimit)
	if !errors.Is(wrapped, ErrRateLimit) {
		t.Error("wrapped ErrRateLimit should match")
	}
	if errors.Is(wrapped, ErrTemporary) {
		t.Error("wrapped ErrRateLimit should not match ErrTemporary")
	}
}

func TestErrorsIsAllTypes(t *testing.T) {
	sentinels := []error{ErrNotRunning, ErrRateLimit, ErrTemporary, ErrSendFailed}

	for _, sentinel := range sentinels {
		wrapped := fmt.Errorf("context: %w", sentinel)
		if !errors.Is(wrapped, sentinel) {
			t.Errorf("wrapped %v should match itself", sentinel)
		}

		// Verify it doesn't match other sentinel errors
		for _, other := range sentinels {
			if other == sentinel {
				continue
			}
			if errors.Is(wrapped, other) {
				t.Errorf("wrapped %v should not match %v", sentinel, other)
			}
		}
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{ErrNotRunning, "channel not running"},
		{ErrRateLimit, "rate limited"},
		{ErrTemporary, "temporary failure"},
		{ErrSendFailed, "send failed"},
	}

	for _, tt := range tests {
		if got := tt.err.Error(); got != tt.want {
			t.Errorf("error message = %q, want %q", got, tt.want)
		}
	}
}
