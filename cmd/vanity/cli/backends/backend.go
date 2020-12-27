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

// Package backends contains the shared implementation of vanity.Backend.
package backends

import (
	"context"

	"l7e.io/vanity"
)

/*
Backend is the shared implementation of vanity.Backend, installed by one of
the backend sub-commands of the vanity root command, cli.RootCmd.

Sub-commands of backend sub-commands, e.g. add and list, can use this to perform
their functionality.
*/
var backend vanity.Backend

func init() {
	backend = &doNothing{}
}

func Set(be vanity.Backend) {
	if be != nil {
		backend = be
	}
}

func Get() vanity.Backend {
	return backend
}

type doNothing struct {
}

func (s *doNothing) Close() error {
	return nil
}

func (s *doNothing) Get(_ context.Context, _ string) (string, string, error) {
	return "", "", vanity.ErrNotFound
}

func (s *doNothing) Add(_ context.Context, _, _, _ string) error {
	return nil
}

func (s *doNothing) Remove(_ context.Context, _ string) error {
	return vanity.ErrNotFound
}

func (s *doNothing) List(_ context.Context, _ vanity.Consumer) error {
	return nil
}

func (s *doNothing) Healthz(_ context.Context) error {
	return nil
}
