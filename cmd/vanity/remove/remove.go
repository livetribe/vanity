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

// Package remove contains the remove sub-command to remove a vanity URL.
package remove

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"l7e.io/vanity/cmd/vanity/cli/backends"
	"l7e.io/vanity/cmd/vanity/cli/backends/helpers"
)

const removeDescription = "Remove vanity URL"

func init() { //nolint:gochecknoinits
	helpers.AddCommand(func() *cobra.Command {
		return &cobra.Command{
			Use:   "remove <importPath>",
			Short: removeDescription,
			Long:  removeDescription,
			Args:  cobra.ExactArgs(1),
			RunE:  removeCmd,
		}
	})
}

func removeCmd(cmd *cobra.Command, args []string) error {
	importPath := args[0]

	slog.Debug("Removing the vanity URL", slog.String("importPath", importPath))

	ctx := cmd.Context()

	if err := backends.Get().Remove(ctx, importPath); err != nil {
		return fmt.Errorf("unable to remove %s: %w", importPath, err)
	}

	slog.Debug("Removed the vanity URL")

	return nil
}
