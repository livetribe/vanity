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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricNamespace = "vanity"
	metricSubsystem = "api"
)

// docRedirectCounter counts each request that the vanity handler answers with a
// redirect to the documentation site.
var docRedirectCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: metricNamespace,
	Subsystem: metricSubsystem,
	Name:      "doc_total",
	Help:      "The total vanity Backend doc redirects",
})

// WithPrometheus returns a http.Handler that sends all requests to next. The
// returned handler counts each documentation redirect that next writes in the
// vanity_api_doc_total counter.
func WithPrometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		if wrappedWriter.statusCode == http.StatusTemporaryRedirect {
			docRedirectCounter.Inc()
		}
	})
}
