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

package spanner // import "l7e.io/vanity/pkg/gcp/spanner"

import (
	"cloud.google.com/go/spanner"
	"google.golang.org/api/option"
)

const (

	// DefaultTable is the default Spanner table name.
	DefaultTable = "urls"

	// DefaultNumChannels is the default number of channels in the Spanner client.
	DefaultNumChannels = 10
)

type backendSettings struct {
	table   string
	config  *spanner.ClientConfig
	options []option.ClientOption
}

// A BackendOption is an option for a Spanner-based backend.
type BackendOption interface {
	Apply(*backendSettings)
}

// WithTable configures the Spanner table; default is "urls".
func WithTable(t string) BackendOption {
	return withTable{t}
}

type withTable struct{ t string }

func (w withTable) Apply(o *backendSettings) {
	o.table = w.t
}

// WithClientConfig returns a BackendOption that specifies configurations for
// the Spanner client.
func WithClientConfig(c spanner.ClientConfig) BackendOption {
	cc := c

	return withClientConfig{&cc}
}

type withClientConfig struct{ c *spanner.ClientConfig }

func (w withClientConfig) Apply(o *backendSettings) {
	o.config = w.c
}

// WithClientOptions returns a BackendOption that specifies Google API client
// configurations for the Spanner client.
func WithClientOptions(o []option.ClientOption) BackendOption {
	return withClientOptions{o}
}

type withClientOptions struct{ o []option.ClientOption }

func (w withClientOptions) Apply(o *backendSettings) {
	o.options = w.o
}

func collectSettings(opts ...BackendOption) *backendSettings {
	bs := &backendSettings{
		table: DefaultTable,
		config: &spanner.ClientConfig{
			NumChannels: DefaultNumChannels,
		},
	}

	for _, o := range opts {
		o.Apply(bs)
	}

	return bs
}
