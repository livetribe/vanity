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

// Package cli contains the RootCmd for main and utilities for printing.
package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"l7e.io/vanity/cmd/vanity/cli/log"
)

const rootDescription = "A go vanity server"

const levelUsage = "log level: trace, debug, info, warn, or error"

// RootCmd is the root cobra command for vanity. The command writes its own
// errors to the log, so cobra does not write them again.
var RootCmd = &cobra.Command{
	Use:           "vanity",
	Short:         rootDescription,
	Long:          rootDescription,
	Version:       getVersion(),
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() { //nolint:gochecknoinits
	flags := RootCmd.PersistentFlags()
	flags.StringP(log.LevelFlag, "", log.DefaultLevel, levelUsage)

	_ = viper.BindPFlag(log.LevelFlag, flags.Lookup(log.LevelFlag))

	cobra.OnInitialize(setLevel)
}

// setLevel gives the configured log level to the default logger. setLevel
// keeps the default level when the configuration holds a name that is not a
// log level.
func setLevel() {
	name := viper.GetString(log.LevelFlag)

	if err := log.SetLevel(name); err != nil {
		slog.Warn("Unable to set the log level", slog.String("level", name), slog.Any("error", err))
	}
}
