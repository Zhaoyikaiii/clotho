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

package wecom

import (
	"context"
	"testing"

	"github.com/Zhaoyikaiii/clotho/pkg/bus"
	"github.com/Zhaoyikaiii/clotho/pkg/config"
)

func TestNewWeComAIBotChannel(t *testing.T) {
	t.Run("success with valid config", func(t *testing.T) {
		cfg := config.WeComAIBotConfig{
			Enabled:        true,
			Token:          "test_token",
			EncodingAESKey: "testkey1234567890123456789012345678901234567",
			WebhookPath:    "/webhook/test",
		}

		messageBus := bus.NewMessageBus()
		ch, err := NewWeComAIBotChannel(cfg, messageBus)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if ch == nil {
			t.Fatal("Expected channel to be created")
		}

		if ch.Name() != "wecom_aibot" {
			t.Errorf("Expected name 'wecom_aibot', got '%s'", ch.Name())
		}
	})

	t.Run("error with missing token", func(t *testing.T) {
		cfg := config.WeComAIBotConfig{
			Enabled:        true,
			EncodingAESKey: "testkey1234567890123456789012345678901234567",
		}

		messageBus := bus.NewMessageBus()
		_, err := NewWeComAIBotChannel(cfg, messageBus)

		if err == nil {
			t.Fatal("Expected error for missing token, got nil")
		}
	})

	t.Run("error with missing encoding key", func(t *testing.T) {
		cfg := config.WeComAIBotConfig{
			Enabled: true,
			Token:   "test_token",
		}

		messageBus := bus.NewMessageBus()
		_, err := NewWeComAIBotChannel(cfg, messageBus)

		if err == nil {
			t.Fatal("Expected error for missing encoding key, got nil")
		}
	})
}

func TestWeComAIBotChannelStartStop(t *testing.T) {
	cfg := config.WeComAIBotConfig{
		Enabled:        true,
		Token:          "test_token",
		EncodingAESKey: "testkey1234567890123456789012345678901234567",
	}

	messageBus := bus.NewMessageBus()
	ch, err := NewWeComAIBotChannel(cfg, messageBus)
	if err != nil {
		t.Fatalf("Failed to create channel: %v", err)
	}

	ctx := context.Background()

	// Test Start
	if err := ch.Start(ctx); err != nil {
		t.Fatalf("Failed to start channel: %v", err)
	}

	if !ch.IsRunning() {
		t.Error("Expected channel to be running")
	}

	// Test Stop
	if err := ch.Stop(ctx); err != nil {
		t.Fatalf("Failed to stop channel: %v", err)
	}

	if ch.IsRunning() {
		t.Error("Expected channel to be stopped")
	}
}

func TestWeComAIBotChannelWebhookPath(t *testing.T) {
	t.Run("default path", func(t *testing.T) {
		cfg := config.WeComAIBotConfig{
			Enabled:        true,
			Token:          "test_token",
			EncodingAESKey: "testkey1234567890123456789012345678901234567",
		}

		messageBus := bus.NewMessageBus()
		ch, _ := NewWeComAIBotChannel(cfg, messageBus)

		expectedPath := "/webhook/wecom-aibot"
		if ch.WebhookPath() != expectedPath {
			t.Errorf("Expected webhook path '%s', got '%s'", expectedPath, ch.WebhookPath())
		}
	})

	t.Run("custom path", func(t *testing.T) {
		customPath := "/custom/webhook"
		cfg := config.WeComAIBotConfig{
			Enabled:        true,
			Token:          "test_token",
			EncodingAESKey: "testkey1234567890123456789012345678901234567",
			WebhookPath:    customPath,
		}

		messageBus := bus.NewMessageBus()
		ch, _ := NewWeComAIBotChannel(cfg, messageBus)

		if ch.WebhookPath() != customPath {
			t.Errorf("Expected webhook path '%s', got '%s'", customPath, ch.WebhookPath())
		}
	})
}

func TestGenerateStreamID(t *testing.T) {
	cfg := config.WeComAIBotConfig{
		Enabled:        true,
		Token:          "test_token",
		EncodingAESKey: "testkey1234567890123456789012345678901234567",
	}

	messageBus := bus.NewMessageBus()
	ch, _ := NewWeComAIBotChannel(cfg, messageBus)

	// Generate multiple IDs and check they are unique
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := ch.generateStreamID()

		if len(id) != 10 {
			t.Errorf("Expected stream ID length 10, got %d", len(id))
		}

		if ids[id] {
			t.Errorf("Duplicate stream ID generated: %s", id)
		}
		ids[id] = true
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Use a valid 43-character base64 key (企业微信标准格式)
	cfg := config.WeComAIBotConfig{
		Enabled:        true,
		Token:          "test_token",
		EncodingAESKey: "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG", // 43 characters
	}

	messageBus := bus.NewMessageBus()
	ch, _ := NewWeComAIBotChannel(cfg, messageBus)

	plaintext := "Hello, World!"
	receiveid := ""

	// Encrypt
	encrypted, err := ch.encryptMessage(plaintext, receiveid)
	if err != nil {
		t.Fatalf("Failed to encrypt message: %v", err)
	}

	if encrypted == "" {
		t.Fatal("Encrypted message is empty")
	}

	// Decrypt
	decrypted, err := decryptMessageWithVerify(encrypted, cfg.EncodingAESKey, receiveid)
	if err != nil {
		t.Fatalf("Failed to decrypt message: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected decrypted message '%s', got '%s'", plaintext, decrypted)
	}
}

func TestGenerateSignature(t *testing.T) {
	token := "test_token"
	timestamp := "1234567890"
	nonce := "test_nonce"
	encrypt := "encrypted_msg"

	signature := computeSignature(token, timestamp, nonce, encrypt)

	if signature == "" {
		t.Error("Generated signature is empty")
	}

	// Verify signature using verifySignature function
	if !verifySignature(token, signature, timestamp, nonce, encrypt) {
		t.Error("Generated signature does not verify correctly")
	}
}
