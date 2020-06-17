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

// Package remove contains the remove sub-command to remove a vanity URL.
package remove

import (
	"context"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"l7e.io/vanity/cmd/vanity/cli/backends"
	"l7e.io/vanity/cmd/vanity/cli/backends/helpers"
	"l7e.io/vanity/cmd/vanity/cli/log"
)

func init() { //nolint:gochecknoinits
	helpers.AddCommand(func() *cobra.Command {
		return &cobra.Command{
			Use:   "remove <importPath>",
			Short: "Remove vanity URL",
			Long:  "Remove vanity URL",
			Args:  cobra.ExactArgs(1),
			Run:   removeCmd,
		}
	})
}

func removeCmd(_ *cobra.Command, args []string) {
	importPath := args[0]

	glog.V(log.Debug).Infof("Removing %s...", importPath)

	err := backends.Backend.Remove(context.Background(), importPath)
	if err != nil {
		glog.Exitf("Unable to remove %s: %s", importPath, err)
	}

	glog.V(log.Debug).Info("Removed")
}
