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

// Package interceptors contains the server sub-command to serve vanity URLs.
package interceptors

import "net/http"

var factories []InterceptorFactory

// InterceptorFactory can be used to wrap a Handler with interceptor code.
// Such factories are intended be used with the RegisterInterceptor() function
// where Backend implementations can register custom HTTP interceptors.
type InterceptorFactory func(h http.Handler) http.Handler

// RegisterInterceptor is used to register InterceptorFactory instances which
// will be invoked to wrap a Handler passed to WrapHandler().
func RegisterInterceptor(f InterceptorFactory) {
	factories = append(factories, f)
}

// WrapHandler wraps the Handler with the set of interceptors created by
// registered InterceptorFactory instances.
func WrapHandler(h http.Handler) http.Handler {
	wrapped := h
	for _, f := range factories {
		wrapped = f(wrapped)
	}
	return wrapped
}
