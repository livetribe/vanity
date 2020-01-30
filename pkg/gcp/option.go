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

// Package gcp contains common configuration options for GCP backends.
package gcp

import (
	"google.golang.org/api/option"
)

// APISettings holds a collection of Google API client options.
type APISettings struct {
	options []option.ClientOption
}

// An BackendOption is an option for a GCP API based client.
type BackendOption interface {
	Apply(*APISettings)
}

// WithClientOptions returns a BackendOption that specifies Google API client
// configurations for the Spanner client.
func WithClientOptions(o []option.ClientOption) BackendOption {
	return withClientOptions{o}
}

type withClientOptions struct{ o []option.ClientOption }

func (w withClientOptions) Apply(a *APISettings) {
	a.options = w.o
}
