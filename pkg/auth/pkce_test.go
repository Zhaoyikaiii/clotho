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
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestGeneratePKCE(t *testing.T) {
	codes, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE() error: %v", err)
	}

	if codes.CodeVerifier == "" {
		t.Fatal("CodeVerifier is empty")
	}
	if codes.CodeChallenge == "" {
		t.Fatal("CodeChallenge is empty")
	}

	verifierBytes, err := base64.RawURLEncoding.DecodeString(codes.CodeVerifier)
	if err != nil {
		t.Fatalf("CodeVerifier is not valid base64url: %v", err)
	}
	if len(verifierBytes) != 64 {
		t.Errorf("CodeVerifier decoded length = %d, want 64", len(verifierBytes))
	}

	hash := sha256.Sum256([]byte(codes.CodeVerifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	if codes.CodeChallenge != expectedChallenge {
		t.Errorf("CodeChallenge = %q, want SHA256 of verifier = %q", codes.CodeChallenge, expectedChallenge)
	}
}

func TestGeneratePKCEUniqueness(t *testing.T) {
	codes1, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE() error: %v", err)
	}

	codes2, err := GeneratePKCE()
	if err != nil {
		t.Fatalf("GeneratePKCE() error: %v", err)
	}

	if codes1.CodeVerifier == codes2.CodeVerifier {
		t.Error("two GeneratePKCE() calls produced identical verifiers")
	}
}
