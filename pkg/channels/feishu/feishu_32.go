//go:build !amd64 && !arm64 && !riscv64 && !mips64 && !ppc64

package feishu

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
	"errors"

	"github.com/Zhaoyikaiii/clotho/pkg/bus"
	"github.com/Zhaoyikaiii/clotho/pkg/channels"
	"github.com/Zhaoyikaiii/clotho/pkg/config"
)

// FeishuChannel is a stub implementation for 32-bit architectures
type FeishuChannel struct {
	*channels.BaseChannel
}

var errUnsupported = errors.New("feishu channel is not supported on 32-bit architectures")

// NewFeishuChannel returns an error on 32-bit architectures where the Feishu SDK is not supported
func NewFeishuChannel(cfg config.FeishuConfig, bus *bus.MessageBus) (*FeishuChannel, error) {
	return nil, errors.New(
		"feishu channel is not supported on 32-bit architectures (armv7l, 386, etc.). Please use a 64-bit system or disable feishu in your config",
	)
}

// Start is a stub method to satisfy the Channel interface
func (c *FeishuChannel) Start(ctx context.Context) error {
	return errUnsupported
}

// Stop is a stub method to satisfy the Channel interface
func (c *FeishuChannel) Stop(ctx context.Context) error {
	return errUnsupported
}

// Send is a stub method to satisfy the Channel interface
func (c *FeishuChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	return errUnsupported
}

// EditMessage is a stub method to satisfy MessageEditor
func (c *FeishuChannel) EditMessage(ctx context.Context, chatID, messageID, content string) error {
	return errUnsupported
}

// SendPlaceholder is a stub method to satisfy PlaceholderCapable
func (c *FeishuChannel) SendPlaceholder(ctx context.Context, chatID string) (string, error) {
	return "", errUnsupported
}

// ReactToMessage is a stub method to satisfy ReactionCapable
func (c *FeishuChannel) ReactToMessage(ctx context.Context, chatID, messageID string) (func(), error) {
	return func() {}, errUnsupported
}

// SendMedia is a stub method to satisfy MediaSender
func (c *FeishuChannel) SendMedia(ctx context.Context, msg bus.OutboundMediaMessage) error {
	return errUnsupported
}
