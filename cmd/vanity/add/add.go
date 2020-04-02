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

// Package add contains the add sub-command to add vanity URLs.
package add

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
			Use:   "add <importPath> <vcs> <vcsPath>",
			Short: "Add vanity URL",
			Long:  "Add vanity URL",
			Args:  cobra.ExactArgs(3), // nolint
			Run:   addCmd,
		}
	})
}

func addCmd(_ *cobra.Command, args []string) {
	importPath := args[0]
	vcs := args[1]
	vcsPath := args[2]

	glog.V(log.Debug).Infof("Adding %s %s %s...", importPath, vcs, vcsPath)

	err := backends.Backend.Add(context.Background(), importPath, vcs, vcsPath)
	if err != nil {
		glog.Exitf("Unable to add %s %s %s: %s", importPath, vcs, vcsPath, err)
	}

	glog.V(log.Debug).Info("Added")
}
