/*
 * Copyright the original author or authors.
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

package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// docRedirects returns the current value of docRedirectCounter.
func docRedirects(t *testing.T) float64 {
	t.Helper()

	m := &dto.Metric{}
	err := docRedirectCounter.Write(m)
	require.NoError(t, err)

	return m.Counter.GetValue()
}

// serveStatus sends one GET request through WithPrometheus to a handler that
// writes code. serveStatus returns the increase of docRedirectCounter.
func serveStatus(t *testing.T, code int) float64 {
	t.Helper()

	before := docRedirects(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(code)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), "GET", "https://a.com/b", http.NoBody)
	WithPrometheus(next).ServeHTTP(w, r)

	assert.Equal(t, code, w.Result().StatusCode)

	return docRedirects(t) - before
}

func TestWithPrometheus_docRedirect(t *testing.T) {
	assert.Equal(t, float64(1), serveStatus(t, http.StatusTemporaryRedirect))
}

func TestWithPrometheus_notFound(t *testing.T) {
	assert.Equal(t, float64(0), serveStatus(t, http.StatusNotFound))
}

func TestWithPrometheus_ok(t *testing.T) {
	assert.Equal(t, float64(0), serveStatus(t, http.StatusOK))
}
