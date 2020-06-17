/*
 * Copyright (c) 2020 the original author or authors.
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

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errUnhealthy = fmt.Errorf("unhealthy")

func TestNewHandlerCheck_ok(t *testing.T) {
	c := newHandlerCheck(&be{healthy: nil}, "")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://a.com", nil)

	c.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNewHandlerCheck_error(t *testing.T) {
	c := newHandlerCheck(&be{healthy: errUnhealthy}, "")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://a.com", nil)

	c.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}
