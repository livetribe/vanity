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

package cli

import (
	"context"
	"fmt"

	"l7e.io/vanity"
)

// NewPlainConsumer creates a Consumer that prints vanity URL configurations
// to standard output as plainConsumer comma-delimited text.
func NewPlainConsumer() vanity.Consumer {
	return plainConsumer{}
}

// NewJSONConsumer creates a Consumer that prints vanity URL configurations
// to standard output as a JSON object.
func NewJSONConsumer() vanity.Consumer {
	return &jsonConsumer{first: true}
}

type plainConsumer struct{}

func (p plainConsumer) OnEntry(_ context.Context, importPath, vcs, vcsPath string) {
	fmt.Printf("%s,%s,%s\n", importPath, vcs, vcsPath)
}

type jsonConsumer struct {
	first bool
}

func (j *jsonConsumer) OnEntry(_ context.Context, importPath, vcs, vcsPath string) {
	if j.first {
		j.first = false
	} else {
		fmt.Printf(",\n")
	}
	fmt.Printf("{\"importPath\": \"%s\", \"vcs\": \"%s\", \"vcsPath\": \"%s\"}", importPath, vcs, vcsPath)
}
