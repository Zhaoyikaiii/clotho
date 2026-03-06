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
	"sync"
	"testing"
)

func TestMessageDeduplicator_DuplicateDetection(t *testing.T) {
	d := NewMessageDeduplicator(wecomMaxProcessedMessages)

	if ok := d.MarkMessageProcessed("msg-1"); !ok {
		t.Fatalf("first message should be accepted")
	}

	if ok := d.MarkMessageProcessed("msg-1"); ok {
		t.Fatalf("duplicate message should be rejected")
	}
}

func TestMessageDeduplicator_ConcurrentSameMessage(t *testing.T) {
	d := NewMessageDeduplicator(wecomMaxProcessedMessages)

	const goroutines = 64
	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make(chan bool, goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			results <- d.MarkMessageProcessed("msg-concurrent")
		}()
	}

	wg.Wait()
	close(results)

	successes := 0
	for ok := range results {
		if ok {
			successes++
		}
	}

	if successes != 1 {
		t.Fatalf("expected exactly 1 successful mark, got %d", successes)
	}
}

func TestMessageDeduplicator_CircularQueueEviction(t *testing.T) {
	// Create a deduplicator with a very small capacity to test eviction easily.
	capacity := 3
	d := NewMessageDeduplicator(capacity)

	// Fill the queue.
	d.MarkMessageProcessed("msg-1")
	d.MarkMessageProcessed("msg-2")
	d.MarkMessageProcessed("msg-3")

	// At this point, the queue is full. msg-1 is the oldest.
	if len(d.msgs) != 3 {
		t.Fatalf("expected map size to be 3, got %d", len(d.msgs))
	}

	// This should evict msg-1 and add msg-4.
	if ok := d.MarkMessageProcessed("msg-4"); !ok {
		t.Fatalf("msg-4 should be accepted")
	}

	if len(d.msgs) != 3 {
		t.Fatalf("expected map size to remain at max capacity (3), got %d", len(d.msgs))
	}

	// msg-1 should now be forgotten (evicted).
	if ok := d.MarkMessageProcessed("msg-1"); !ok {
		t.Fatalf("msg-1 should be accepted again because it was evicted")
	}

	// msg-2 should have been evicted when we added msg-1 back.
	if ok := d.MarkMessageProcessed("msg-2"); !ok {
		t.Fatalf("msg-2 should be accepted again because it was evicted")
	}
}
