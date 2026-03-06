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

import "context"

// TypingCapable — channels that can show a typing/thinking indicator.
// StartTyping begins the indicator and returns a stop function.
// The stop function MUST be idempotent and safe to call multiple times.
type TypingCapable interface {
	StartTyping(ctx context.Context, chatID string) (stop func(), err error)
}

// MessageEditor — channels that can edit an existing message.
// messageID is always string; channels convert platform-specific types internally.
type MessageEditor interface {
	EditMessage(ctx context.Context, chatID string, messageID string, content string) error
}

// ReactionCapable — channels that can add a reaction (e.g. 👀) to an inbound message.
// ReactToMessage adds a reaction and returns an undo function to remove it.
// The undo function MUST be idempotent and safe to call multiple times.
type ReactionCapable interface {
	ReactToMessage(ctx context.Context, chatID, messageID string) (undo func(), err error)
}

// PlaceholderCapable — channels that can send a placeholder message
// (e.g. "Thinking... 💭") that will later be edited to the actual response.
// The channel MUST also implement MessageEditor for the placeholder to be useful.
// SendPlaceholder returns the platform message ID of the placeholder so that
// Manager.preSend can later edit it via MessageEditor.EditMessage.
type PlaceholderCapable interface {
	SendPlaceholder(ctx context.Context, chatID string) (messageID string, err error)
}

// PlaceholderRecorder is injected into channels by Manager.
// Channels call these methods on inbound to register typing/placeholder state.
// Manager uses the registered state on outbound to stop typing and edit placeholders.
type PlaceholderRecorder interface {
	RecordPlaceholder(channel, chatID, placeholderID string)
	RecordTypingStop(channel, chatID string, stop func())
	RecordReactionUndo(channel, chatID string, undo func())
}
