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

package spanner

import (
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

func TestCollectSettingsDefaults(t *testing.T) {
	s := collectSettings()

	assert.Equal(t, DefaultTable, s.table)
	assert.Nil(t, s.config)
	assert.Equal(t, []option.ClientOption{option.WithGRPCConnectionPool(DefaultNumChannels)}, s.options)
}

func TestCollectSettingsWithTable(t *testing.T) {
	s := collectSettings(WithTable("other"))

	assert.Equal(t, "other", s.table)
}

func TestCollectSettingsPrependsConnectionPool(t *testing.T) {
	agent := option.WithUserAgent("test")

	s := collectSettings(WithClientOptions([]option.ClientOption{agent}))

	require.Len(t, s.options, 2)
	assert.Equal(t, option.WithGRPCConnectionPool(DefaultNumChannels), s.options[0])
	assert.Equal(t, agent, s.options[1])
}

func TestCollectSettingsClientConfigOwnsThePool(t *testing.T) {
	config := spanner.ClientConfig{NumChannels: 2} //nolint:staticcheck

	s := collectSettings(WithClientConfig(config))

	require.NotNil(t, s.config)
	assert.Empty(t, s.options)
}
