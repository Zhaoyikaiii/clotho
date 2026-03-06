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

package auth

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func LoginPasteToken(provider string, r io.Reader) (*AuthCredential, error) {
	fmt.Printf("Paste your API key or session token from %s:\n", providerDisplayName(provider))
	fmt.Print("> ")

	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("reading token: %w", err)
		}
		return nil, fmt.Errorf("no input received")
	}

	token := strings.TrimSpace(scanner.Text())
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	return &AuthCredential{
		AccessToken: token,
		Provider:    provider,
		AuthMethod:  "token",
	}, nil
}

func providerDisplayName(provider string) string {
	switch provider {
	case "anthropic":
		return "console.anthropic.com"
	case "openai":
		return "platform.openai.com"
	default:
		return provider
	}
}
