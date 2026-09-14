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

// Package add contains the add sub-command to add vanity URLs.
package add

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"l7e.io/vanity/cmd/vanity/cli/backends"
	"l7e.io/vanity/cmd/vanity/cli/backends/helpers"
)

const addDescription = "Add vanity URL"

func init() { //nolint:gochecknoinits
	helpers.AddCommand(func() *cobra.Command {
		return &cobra.Command{
			Use:   "add <importPath> <vcs> <vcsPath>",
			Short: addDescription,
			Long:  addDescription,
			Args:  cobra.ExactArgs(3), //nolint:mnd
			RunE:  addCmd,
		}
	})
}

func addCmd(cmd *cobra.Command, args []string) error {
	importPath := args[0]
	vcs := args[1]
	vcsPath := args[2]

	slog.Debug("Adding the vanity URL",
		slog.String("importPath", importPath), slog.String("vcs", vcs), slog.String("vcsPath", vcsPath))

	ctx := cmd.Context()

	if err := backends.Get().Add(ctx, importPath, vcs, vcsPath); err != nil {
		return fmt.Errorf("unable to add %s: %w", importPath, err)
	}

	slog.Debug("Added the vanity URL")

	return nil
}
