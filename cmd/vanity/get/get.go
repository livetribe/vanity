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

/*
Package get contains the get sub-command to get vanity URLs.
*/
package get

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

const getDescription = "Get vanity URL"

func init() { //nolint:gochecknoinits
	helpers.AddCommand(func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "get <importPath>",
			Short: getDescription,
			Long:  getDescription,
			Args:  cobra.ExactArgs(1),
			RunE:  getCmd,
		}

		flags := cmd.Flags()
		flags.BoolVarP(&outputJSON, "json", "", false, "output in JSON")

		return cmd
	})
}

func getCmd(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("unable to bind viper to the command line flags: %w", err)
	}

	var c vanity.Consumer
	if outputJSON {
		c = cli.NewJSONConsumer()
	} else {
		c = cli.NewPlainConsumer()
	}

	ctx := cmd.Context()
	importPath := args[0]

	vcs, vcsPath, err := backends.Get().Get(ctx, importPath)
	if err != nil {
		return fmt.Errorf("unable to get %s: %w", importPath, err)
	}

	c.OnEntry(ctx, importPath, vcs, vcsPath)

	return nil
}
