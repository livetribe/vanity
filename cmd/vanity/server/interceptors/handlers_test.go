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

package interceptors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"l7e.io/vanity/cmd/vanity/server/interceptors"
)

func TestWrapHandler(t *testing.T) {
	var values []string

	interceptors.RegisterInterceptor(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			values = append(values, "c")
			h.ServeHTTP(w, r)
		})
	})
	interceptors.RegisterInterceptor(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			values = append(values, "b")
			h.ServeHTTP(w, r)
		})
	})
	interceptors.RegisterInterceptor(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			values = append(values, "a")
			h.ServeHTTP(w, r)
		})
	})

	h := interceptors.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values = append(values, "d")
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://a.com", nil)
	h.ServeHTTP(w, r)

	assert.Equal(t, []string{"a", "b", "c", "d"}, values)
}
