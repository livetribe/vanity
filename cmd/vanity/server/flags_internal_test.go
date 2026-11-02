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

package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"l7e.io/vanity/cmd/vanity/cmdtest"
)

// noFail is a Fail that records nothing. The tests here build listeners but
// never serve them.
func noFail() {}

// withFlags runs check against the configuration that args produce.
func withFlags(t *testing.T, check func(cfg Config), args ...string) {
	t.Helper()

	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, _ []string) {
		err := viper.BindPFlags(cmd.Flags())
		require.NoError(t, err)

		check(NewConfig())
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, args...)
	require.NoError(t, err)
}

func TestNewVanity_addr_default(t *testing.T) {
	withFlags(t, func(cfg Config) {
		v := NewVanity(cfg, &be{}, noFail)
		assert.Equal(t, "127.0.1.2:8080", v.server.Addr)
	}, "--bind", "127.0.1.2")
}

func TestNewVanity_addr_port(t *testing.T) {
	withFlags(t, func(cfg Config) {
		v := NewVanity(cfg, &be{}, noFail)
		assert.Equal(t, "127.0.1.2:1234", v.server.Addr)
	}, "--bind", "127.0.1.2", "--port", "1234")
}

func TestNewHealthz_addr_default(t *testing.T) {
	withFlags(t, func(cfg Config) {
		h := NewHealthz(cfg, &be{}, noFail)
		assert.Equal(t, "127.0.1.2:8081", h.server.Addr)
	}, "--bind", "127.0.1.2")
}

func TestNewHealthz_addr_port(t *testing.T) {
	withFlags(t, func(cfg Config) {
		h := NewHealthz(cfg, &be{}, noFail)
		assert.Equal(t, "127.0.1.2:1234", h.server.Addr)
	}, "--bind", "127.0.1.2", "--healthz", "1234")
}

func TestNewReadyz_addr_default(t *testing.T) {
	withFlags(t, func(cfg Config) {
		r := NewReadyz(cfg, &be{}, noFail)
		assert.Equal(t, "127.0.1.2:8082", r.server.Addr)
	}, "--bind", "127.0.1.2")
}

func TestNewReadyz_addr_port(t *testing.T) {
	withFlags(t, func(cfg Config) {
		r := NewReadyz(cfg, &be{}, noFail)
		assert.Equal(t, "127.0.1.2:1234", r.server.Addr)
	}, "--bind", "127.0.1.2", "--readyz", "1234")
}

func TestNewMetrics_addr_default(t *testing.T) {
	withFlags(t, func(cfg Config) {
		m := NewMetrics(cfg, noFail)
		assert.Equal(t, "127.0.1.2:9100", m.server.Addr)
	}, "--bind", "127.0.1.2")
}

func TestNewMetrics_addr_port(t *testing.T) {
	withFlags(t, func(cfg Config) {
		m := NewMetrics(cfg, noFail)
		assert.Equal(t, "127.0.1.2:1234", m.server.Addr)
	}, "--bind", "127.0.1.2", "--prometheus", "1234")
}
