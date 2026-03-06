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

package server

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogBuffer_Basic(t *testing.T) {
	buf := NewLogBuffer(5)

	// Empty buffer
	lines, total, runID := buf.LinesSince(0)
	assert.Nil(t, lines)
	assert.Equal(t, 0, total)
	assert.Equal(t, 0, runID)

	// Append some lines
	buf.Append("line1")
	buf.Append("line2")
	buf.Append("line3")

	lines, total, runID = buf.LinesSince(0)
	assert.Equal(t, []string{"line1", "line2", "line3"}, lines)
	assert.Equal(t, 3, total)
	assert.Equal(t, 0, runID)

	// Incremental read
	lines, total, _ = buf.LinesSince(2)
	assert.Equal(t, []string{"line3"}, lines)
	assert.Equal(t, 3, total)

	// No new lines
	lines, total, _ = buf.LinesSince(3)
	assert.Nil(t, lines)
	assert.Equal(t, 3, total)
}

func TestLogBuffer_Wrap(t *testing.T) {
	buf := NewLogBuffer(3)

	buf.Append("a")
	buf.Append("b")
	buf.Append("c")
	buf.Append("d") // evicts "a"
	buf.Append("e") // evicts "b"

	lines, total, _ := buf.LinesSince(0)
	assert.Equal(t, []string{"c", "d", "e"}, lines)
	assert.Equal(t, 5, total)

	// Incremental after wrap
	lines, total, _ = buf.LinesSince(3)
	assert.Equal(t, []string{"d", "e"}, lines)
	assert.Equal(t, 5, total)

	// Offset too old (before buffer start), get all buffered
	lines, total, _ = buf.LinesSince(1)
	assert.Equal(t, []string{"c", "d", "e"}, lines)
	assert.Equal(t, 5, total)
}

func TestLogBuffer_Reset(t *testing.T) {
	buf := NewLogBuffer(5)

	buf.Append("before")
	assert.Equal(t, 0, buf.RunID())

	buf.Reset()
	assert.Equal(t, 1, buf.RunID())
	assert.Equal(t, 0, buf.Total())

	lines, total, runID := buf.LinesSince(0)
	assert.Nil(t, lines)
	assert.Equal(t, 0, total)
	assert.Equal(t, 1, runID)

	buf.Append("after")
	lines, total, runID = buf.LinesSince(0)
	assert.Equal(t, []string{"after"}, lines)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, runID)
}

func TestLogBuffer_Concurrent(t *testing.T) {
	buf := NewLogBuffer(100)
	var wg sync.WaitGroup

	// 10 writers
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range 50 {
				buf.Append(fmt.Sprintf("writer-%d-line-%d", id, j))
			}
		}(i)
	}

	// 5 readers
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				buf.LinesSince(0)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, 500, buf.Total())
}
