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

import "net/http"

// WebhookHandler is an optional interface for channels that receive messages
// via HTTP webhooks. Manager discovers channels implementing this interface
// and registers them on the shared HTTP server.
type WebhookHandler interface {
	// WebhookPath returns the path to mount this handler on the shared server.
	// Examples: "/webhook/line", "/webhook/wecom"
	WebhookPath() string
	http.Handler // ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// HealthChecker is an optional interface for channels that expose
// a health check endpoint on the shared HTTP server.
type HealthChecker interface {
	HealthPath() string
	HealthHandler(w http.ResponseWriter, r *http.Request)
}
