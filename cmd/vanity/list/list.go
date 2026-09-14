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

// Package list contains the list sub-command to list vanity URLs.
package list

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"l7e.io/vanity"
	"l7e.io/vanity/cmd/vanity/cli"
	"l7e.io/vanity/cmd/vanity/cli/backends"
	"l7e.io/vanity/cmd/vanity/cli/backends/helpers"
)

var outputJSON bool

const listDescription = "List vanity URLs"

func init() { //nolint:gochecknoinits
	helpers.AddCommand(func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "list",
			Short: listDescription,
			Long:  listDescription,
			Args:  cobra.NoArgs,
			RunE:  listCmd,
		}

		flags := cmd.Flags()
		flags.BoolVarP(&outputJSON, "json", "", false, "output in JSON")

		return cmd
	})
}

func listCmd(cmd *cobra.Command, _ []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("unable to bind viper to the command line flags: %w", err)
	}

	var c vanity.Consumer
	if outputJSON {
		c = cli.NewJSONConsumer()
		fmt.Println("[")
	} else {
		c = cli.NewPlainConsumer()
	}

	ctx := cmd.Context()

	if err := backends.Get().List(ctx, c); err != nil {
		return fmt.Errorf("unable to list the vanity URLs: %w", err)
	}

	if outputJSON {
		fmt.Println("\n]")
	}

	return nil
}
