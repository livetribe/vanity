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

package server

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"l7e.io/vanity/cmd/vanity/cmdtest"
)

func TestGetHTTPServer_addr_default(t *testing.T) {
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		h := newHelper(cmd)
		server := h.getHTTPServer(&be{})
		assert.Equal(t, "127.0.1.2:8080", server.Addr)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--bind", "127.0.1.2")
	assert.NoError(t, err)
}

func TestGetHTTPServer_addr_port(t *testing.T) {
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		h := newHelper(cmd)
		server := h.getHTTPServer(&be{})
		assert.Equal(t, "127.0.1.2:1234", server.Addr)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--bind", "127.0.1.2", "--port", "1234")
	assert.NoError(t, err)
}

func TestGetHealthz_addr_default(t *testing.T) {
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		h := newHelper(cmd)
		server := h.getHealthz(newHandlerCheck(&be{healthy: nil}, ""))
		assert.Equal(t, "127.0.1.2:8081", server.Addr)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--bind", "127.0.1.2")
	assert.NoError(t, err)
}

func TestGetHealthz_addr_port(t *testing.T) {
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		h := newHelper(cmd)
		server := h.getHealthz(newHandlerCheck(&be{healthy: nil}, ""))
		assert.Equal(t, "127.0.1.2:1234", server.Addr)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--bind", "127.0.1.2", "--healthz", "1234")
	assert.NoError(t, err)
}

func TestGetReadyz_addr_default(t *testing.T) {
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		h := newHelper(cmd)
		server := h.getReadyz(newHandlerCheck(&be{healthy: nil}, ""))
		assert.Equal(t, "127.0.1.2:8082", server.Addr)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--bind", "127.0.1.2")
	assert.NoError(t, err)
}

func TestGetReadyz_addr_port(t *testing.T) {
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		h := newHelper(cmd)
		server := h.getReadyz(newHandlerCheck(&be{healthy: nil}, ""))
		assert.Equal(t, "127.0.1.2:1234", server.Addr)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--bind", "127.0.1.2", "--readyz", "1234")
	assert.NoError(t, err)
}

func TestGetMetrics_addr_default(t *testing.T) {
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		h := newHelper(cmd)
		server := h.getMetrics()
		assert.Equal(t, "127.0.1.2:9100", server.Addr)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--bind", "127.0.1.2")
	assert.NoError(t, err)
}

func TestGetMetric_addr_port(t *testing.T) {
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		h := newHelper(cmd)
		server := h.getMetrics()
		assert.Equal(t, "127.0.1.2:1234", server.Addr)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--bind", "127.0.1.2", "--prometheus", "1234")
	assert.NoError(t, err)
}
