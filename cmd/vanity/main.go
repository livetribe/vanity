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

package main

import (
	"flag"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"l7e.io/vanity"

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

func init() { //nolint:gochecknoinits
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// set vanity logger to glog.Errorf()
	vanity.SetLogger(vanity.LoggerFunc(func(format string, v ...interface{}) {
		glog.Errorf(format, v...)
	}))
}

func main() {
	_ = flag.CommandLine.Parse([]string{}) // used to turn off noisy glog warnings

	if err := setupViper(); err != nil {
		glog.Exit(err)
	}

	if err := cli.RootCmd.Execute(); err != nil {
		glog.Exit(err)
	}
}

func setupViper() error {
	viper.SetEnvPrefix("vanity")
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath("/etc/vanity/")
	viper.AddConfigPath("$HOME/.config/vanity")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if cnf, ok := err.(viper.ConfigFileNotFoundError); ok {
		glog.V(log.Debug).Infof(cnf.Error())

		return nil
	}

	return err
}
