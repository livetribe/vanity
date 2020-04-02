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

// Package gcp declares an InterceptorFactory that wraps a Handler with a user
// agent check for a GCE ingress agent.
package gcp // import "l7e.io/vanity/pkg/gcp"

import "net/http"

const gceIngressUserAgent = "GoogleHC/1.0"

// GLB wraps a Handler with a user agent check for a GCE ingress agent.  It
// will always return an HTTP status OK.
func GLB(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tell Google LB everything is fine
		if r.UserAgent() == gceIngressUserAgent {
			w.WriteHeader(http.StatusOK)

			return
		}

		h.ServeHTTP(w, r)
	})
}
