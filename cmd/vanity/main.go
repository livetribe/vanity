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

package main

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"

	_ "l7e.io/vanity/cmd/vanity/add"
	"l7e.io/vanity/cmd/vanity/cli"
	_ "l7e.io/vanity/cmd/vanity/cli/backends/gcp/datastore"
	_ "l7e.io/vanity/cmd/vanity/cli/backends/gcp/spanner"
	"l7e.io/vanity/cmd/vanity/cli/log"
	_ "l7e.io/vanity/cmd/vanity/get"
	_ "l7e.io/vanity/cmd/vanity/list"
	_ "l7e.io/vanity/cmd/vanity/remove"
	_ "l7e.io/vanity/cmd/vanity/server"
)

// exitCode is the status that main gives to the operating system when the
// command fails.
const exitCode = 1

func init() { //nolint:gochecknoinits
	log.Init(os.Stderr)

	viper.SetEnvPrefix("vanity")
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.AddConfigPath("/etc/vanity/")
	viper.AddConfigPath("$HOME/.config/vanity")
}

func main() {
	if err := run(); err != nil {
		slog.Error("Unable to run the command", slog.Any("error", err))
		os.Exit(exitCode)
	}
}

// run reads the configuration. Then run runs the command that the arguments
// name.
func run() error {
	if err := setupViper(); err != nil {
		return err
	}

	return cli.RootCmd.Execute()
}

func setupViper() error {
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if cnf, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
		slog.Debug("No configuration file", slog.Any("error", cnf))

		return nil
	}

	return err
}
