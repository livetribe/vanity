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

package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestWithClientOptions_Apply(t *testing.T) {
	c := option.WithCredentialsFile("")
	o := WithClientOptions([]option.ClientOption{c})

	settings := &APISettings{}
	o.Apply(settings)

	assert.Equal(t, 1, len(settings.options))
	assert.Equal(t, c, settings.options[0])
}
