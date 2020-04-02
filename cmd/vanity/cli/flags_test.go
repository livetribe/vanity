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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"l7e.io/vanity/cmd/vanity/cmdtest"
)

const (
	apiKey       = "api-key"
	keyValue     = "ABCD"
	envVarKey    = "VANITY_TEST_GOOGLE_API_KEY"
	cfgFileValue = "EFGH"
)

func initFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringP(apiKey, "", "", "API key to be used as the basis for authentication (optional)")

	viper.SetEnvPrefix("VANITY_TEST")
	viper.RegisterAlias(apiKey, "google-api.key")
}

func setupViper(dir string) error {
	viper.AutomaticEnv()

	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return nil
	}

	return err
}

func setupConfigFile(dir string) error {
	f, err := os.Create(path.Join(dir, "config.toml"))
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(fmt.Sprintf(`[google-api]
key = "%s"
`, cfgFileValue))
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func TestFlags_viaCmdLine(t *testing.T) {
	expected := keyValue
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		f := Flags(cmd)
		v, ok := f.GetValue(apiKey)
		assert.True(t, ok)
		assert.Equal(t, expected, v)
	})
	initFlags(cmd)

	_, err := cmdtest.ExecuteCommand(cmd, "--"+apiKey, expected)
	assert.NoError(t, err)
}

func TestFlags_viaEnvVar(t *testing.T) {
	expected := keyValue
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		f := Flags(cmd)
		v, ok := f.GetValue(apiKey)
		assert.True(t, ok)
		assert.Equal(t, expected, v)
	})
	initFlags(cmd)

	err := os.Setenv(envVarKey, expected)
	assert.NoError(t, err)
	defer os.Unsetenv(envVarKey)

	err = setupViper("/")
	assert.NoError(t, err)

	_, err = cmdtest.ExecuteCommand(cmd)
	assert.NoError(t, err)
}

func TestFlags_viaConfigFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "prefix")
	assert.NoError(t, err)
	defer os.Remove(dir)

	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		f := Flags(cmd)
		v, ok := f.GetValue(apiKey)
		assert.True(t, ok)
		assert.Equal(t, cfgFileValue, v)
	})
	initFlags(cmd)

	err = setupConfigFile(dir)
	assert.NoError(t, err)

	err = setupViper(dir)
	assert.NoError(t, err)

	_, err = cmdtest.ExecuteCommand(cmd)
	assert.NoError(t, err)
}

func TestFlags_envVarTakesPrecedence(t *testing.T) {
	expected := keyValue
	dir, err := ioutil.TempDir("", "prefix")
	assert.NoError(t, err)
	defer os.Remove(dir)

	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		f := Flags(cmd)
		v, ok := f.GetValue(apiKey)
		assert.True(t, ok)
		assert.Equal(t, expected, v)
	})
	initFlags(cmd)

	err = os.Setenv(envVarKey, expected)
	assert.NoError(t, err)
	defer os.Unsetenv(envVarKey)

	err = setupConfigFile(dir)
	assert.NoError(t, err)

	err = setupViper(dir)
	assert.NoError(t, err)

	_, err = cmdtest.ExecuteCommand(cmd)
	assert.NoError(t, err)
}

func TestFlags_cmdLineTakesPrecedence(t *testing.T) {
	expected := keyValue
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)

		f := Flags(cmd)
		v, ok := f.GetValue(apiKey)
		assert.True(t, ok)
		assert.Equal(t, expected, v)
	})
	initFlags(cmd)

	err := os.Setenv(envVarKey, "FOOBAR")
	assert.NoError(t, err)
	defer os.Unsetenv(envVarKey)

	err = setupViper("/")
	assert.NoError(t, err)

	_, err = cmdtest.ExecuteCommand(cmd, "--"+apiKey, expected)
	assert.NoError(t, err)
}
