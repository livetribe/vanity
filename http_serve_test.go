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

package vanity_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"l7e.io/vanity"
	"l7e.io/vanity/apitest"
)

var errNotHealthy = fmt.Errorf("not healthy")

func TestHandler_ServeHTTP_put(t *testing.T) {
	prometheusReset()

	h := vanity.NewVanityHandler(&apitest.MockBackend{})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "https://a.com", nil)
	h.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	prometheusCheck(t, 0, 0, 0, 0, 0)
}

func TestHandler_ServeHTTP_get_no_go_get(t *testing.T) {
	prometheusReset()

	h := vanity.NewVanityHandler(&apitest.MockBackend{Urls: map[string][]string{"a.com/b": {"vcs", "vcsPath"}}})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://a.com/b", nil)
	h.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	match, err := regexp.Match(`https://pkg\.go\.dev/a\.com/b`, body)
	assert.NoError(t, err)
	assert.True(t, match, string(body))

	prometheusCheck(t, 1, 0, 0, 1, 0)
}

func TestHandler_ServeHTTP_get_not_found(t *testing.T) {
	prometheusReset()

	h := vanity.NewVanityHandler(&apitest.MockBackend{Urls: map[string][]string{"a.com/b": {"vcs", "vcsPath"}}})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://a.com/z", nil)
	h.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	prometheusCheck(t, 1, 0, 1, 0, 0)
}

func TestHandler_ServeHTTP_not_healthy(t *testing.T) {
	prometheusReset()

	h := vanity.NewVanityHandler(&apitest.MockBackend{Healthy: errNotHealthy})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://a.com/z", nil)
	h.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	prometheusCheck(t, 1, 1, 0, 0, 0)
}

func TestHandler_ServeHTTP_get(t *testing.T) {
	prometheusReset()

	expected := `<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
  <meta name="go-import" content="a.com/b vcs vcsPath/b">
  <meta name="go-source" content="a.com/b vcsPath/b vcsPath/b/tree/master{/dir} vcsPath/b/blob/master{/dir}/{file}#L{line}">
</head>
</html>
`

	h := vanity.NewVanityHandler(&apitest.MockBackend{Urls: map[string][]string{"a.com/b": {"vcs", "vcsPath"}}})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://a.com/b?go-get=1", nil)
	h.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(body))

	prometheusCheck(t, 1, 0, 0, 0, 0)
}

func TestHandler_ServeHTTP_get_extendedPath(t *testing.T) {
	prometheusReset()

	expected := `<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
  <meta name="go-import" content="a.com/b/v1 vcs vcsPath/b">
  <meta name="go-source" content="a.com/b/v1 vcsPath/b vcsPath/b/tree/master{/dir} vcsPath/b/blob/master{/dir}/{file}#L{line}">
</head>
</html>
`

	h := vanity.NewVanityHandler(&apitest.MockBackend{Urls: map[string][]string{"a.com/b": {"vcs", "vcsPath"}}})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://a.com/b/v1?go-get=1", nil)
	h.ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(body))

	prometheusCheck(t, 1, 0, 0, 0, 0)
}
