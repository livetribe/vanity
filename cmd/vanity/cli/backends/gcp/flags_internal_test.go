/*
 * Copyright (c) 2026 the original author or authors.
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
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

func TestCredentialsType(t *testing.T) {
	tests := []struct {
		name string
		cj   string
		want option.CredentialsType
	}{
		{"service account", `{"type":"service_account"}`, option.ServiceAccount},
		{"authorized user", `{"type":"authorized_user"}`, option.AuthorizedUser},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			credType, err := credentialsType([]byte(test.cj))

			require.NoError(t, err)
			assert.Equal(t, test.want, credType)
		})
	}
}

func TestCredentialsTypeUnsupported(t *testing.T) {
	tests := []struct {
		name string
		cj   string
	}{
		{"external account", `{"type":"external_account"}`},
		{"impersonated service account", `{"type":"impersonated_service_account"}`},
		{"external account authorized user", `{"type":"external_account_authorized_user"}`},
		{"gdc service account", `{"type":"gdc_service_account"}`},
		{"unknown type", `{"type":"not_a_real_type"}`},
		{"absent type", `{}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := credentialsType([]byte(test.cj))

			assert.ErrorIs(t, err, errUnsupportedCredentials)
		})
	}
}

func TestCredentialsTypeInvalidJSON(t *testing.T) {
	_, err := credentialsType([]byte("not json"))

	require.Error(t, err)
	assert.NotErrorIs(t, err, errUnsupportedCredentials)
}
