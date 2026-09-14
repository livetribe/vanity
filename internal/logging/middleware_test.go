/*
 * Copyright (c) 2026 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logging

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureLog makes the default logger write JSON to a buffer. captureLog
// returns the buffer. At the end of the test, captureLog sets the default
// logger that the test had before.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer

	handler := slog.NewJSONHandler(&buf, nil)
	previous := slog.Default()

	slog.SetDefault(slog.New(handler))

	t.Cleanup(func() {
		slog.SetDefault(previous)
	})

	return &buf
}

func TestMiddleware_oneHTTPGroup(t *testing.T) {
	buf := captureLog(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), "GET", "https://a.com/b", http.NoBody)
	Middleware(next).ServeHTTP(w, r)

	line := buf.String()

	require.Contains(t, line, "request processed")
	assert.Equal(t, 1, strings.Count(line, `"http":`))
	assert.Contains(t, line, `"status":418`)
	assert.Contains(t, line, `"latency":"`)
}

func TestMiddleware_requestLoggerInContext(t *testing.T) {
	buf := captureLog(t)

	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		FromContext(r.Context()).Info("handler ran")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), "GET", "https://a.com/b", http.NoBody)
	Middleware(next).ServeHTTP(w, r)

	assert.Contains(t, buf.String(), `"msg":"handler ran"`)
	assert.Contains(t, buf.String(), `"path":"/b"`)
}

func TestFromContext_fallback(t *testing.T) {
	assert.Same(t, slog.Default(), FromContext(t.Context()))
}
