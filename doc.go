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
 *
 */

/*
Package vanity contains an HTTP handler that serves Go vanity import paths.

The handler makes an import path from the host of the request and the first
element of the request path. It gets the version control system and the
repository root for that import path from a Backend.

A request with the query parameter go-get=1 receives an HTML document that
contains go-import and go-source meta tags. Every other request receives a
temporary redirect to the Go documentation site.

# Creating a handler

Give NewVanityHandler a Backend:

	backend := memory.NewInMemoryAPI()
	defer backend.Close()

	http.Handle("/", vanity.NewVanityHandler(backend))

The packages under pkg supply Backend implementations for an in-memory store,
a TOML file, Google Cloud Datastore, and Google Cloud Spanner.

# Metrics and logs

This package exports Prometheus collectors. They count the calls to the
Backend, the errors, the not-found results, the documentation redirects, and
the template errors. SummaryVec records the duration of each Backend call.

The handler writes its messages to a log/slog logger from the request context.
The handler uses slog.Default() when the request context has no logger.
*/
package vanity // import "l7e.io/vanity"
