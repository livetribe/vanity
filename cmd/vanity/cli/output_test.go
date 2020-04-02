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

package cli

import (
	"context"
	"testing"

	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
)

func TestPlainConsumer(t *testing.T) {
	c := NewPlainConsumer()

	out := capturer.CaptureOutput(func() {
		c.OnEntry(context.Background(), "a", "b", "c")
		c.OnEntry(context.Background(), "d", "e", "f")
	})

	assert.Equal(t, "a,b,c\nd,e,f\n", out)
}

func TestJSONConsumer(t *testing.T) {
	c := NewJSONConsumer()

	out := capturer.CaptureOutput(func() {
		c.OnEntry(context.Background(), "a", "b", "c")
		c.OnEntry(context.Background(), "d", "e", "f")
	})

	assert.Equal(t, "{\"importPath\": \"a\", \"vcs\": \"b\", \"vcsPath\": \"c\"},\n{\"importPath\": \"d\", \"vcs\": \"e\", \"vcsPath\": \"f\"}", out)
}
