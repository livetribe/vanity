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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"l7e.io/vanity"
)

func TestConsumerFunc(t *testing.T) {
	var ip, v, vp string
	c := vanity.ConsumerFunc(func(context context.Context, importPath, vcs, vcsPath string) {
		ip = importPath
		v = vcs
		vp = vcsPath
	})

	c.OnEntry(context.Background(), "a", "b", "c")

	assert.Equal(t, "a", ip)
	assert.Equal(t, "b", v)
	assert.Equal(t, "c", vp)
}
