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
package gcp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGLB_normal(t *testing.T) {
	called := false
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNotFound)
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://a.com", nil)

	glb := GLB(h)
	glb.ServeHTTP(w, r)

	resp := w.Result()
	assert.True(t, called)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestGLB_behind_glb(t *testing.T) {
	called := false
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://a.com", nil)
	r.Header.Add("User-Agent", gceIngressUserAgent)

	glb := GLB(h)
	glb.ServeHTTP(w, r)

	resp := w.Result()
	assert.False(t, called)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
