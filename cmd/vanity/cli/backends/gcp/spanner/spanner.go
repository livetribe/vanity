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

// Package spanner contains the spanner sub-command.
package spanner

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"l7e.io/vanity"
	"l7e.io/vanity/cmd/vanity/cli"
	"l7e.io/vanity/cmd/vanity/cli/backends"
	"l7e.io/vanity/cmd/vanity/cli/backends/gcp"
	"l7e.io/vanity/cmd/vanity/cli/log"
	be "l7e.io/vanity/pkg/gcp/spanner"
)

const (
	database = "database"
)

func init() { //nolint:gochecknoinits
	cli.RootCmd.AddCommand(Command)

	flags := Command.PersistentFlags()
	flags.StringP(database, "", "", "Spanner database connection string")

	gcp.InitFlags(Command)

	viper.RegisterAlias(database, "spanner.database")
}

// Command is the vanity sub-command for a GCP Spanner backend.
var Command = &cobra.Command{
	Use:   "spanner",
	Short: "Use GCP Spanner for a backend vanity store",
	Long:  "Use GCP Spanner for a backend vanity store",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		glog.V(log.Debug).Infoln("Set backend w/ spanner")
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf(unableToBind, err)
		}

		beHelp := newHelper(cmd)
		backend, err := beHelp.getBackend()
		if err != nil {
			return fmt.Errorf(unableToInstantiate, err)
		}

		backends.Backend = backend

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		glog.V(log.Debug).Infoln("Clean backend of spanner")
		defer func() {
			backends.Backend = nil
		}()
		return backends.Backend.Close()
	},
}

type helper struct {
	*pflag.FlagSet

	gh *gcp.Helper
}

// newHelper wraps the Cobra command's flags with a utility wrapper to assist in
// the creation of a Spanner-based backend.
func newHelper(cmd *cobra.Command) *helper {
	return &helper{FlagSet: cmd.Flags(), gh: gcp.NewHelper(cmd)}
}

// getDatabase obtains the Spanner database connection string; the project,
// instance and optionally database to use. The format is
// projects/{project}/instances/{instance}/databases/{database}, with the
// databases/{database} part being optional.
func (h *helper) getDatabase() (db string, found bool) {
	db = viper.GetString(database)
	return db, db != ""
}

// getBackend returns a Spanner-based api.Backend instance, configured by the helper.
func (h *helper) getBackend() (vanity.Backend, error) {
	db, found := h.getDatabase()
	if !found {
		return nil, fmt.Errorf("unable to obtain database connection string")
	}

	glog.V(log.Debug).Infof("database: %s", db)

	options := h.getClientOptions()

	return be.NewClient(context.Background(), db, options...)
}

func (h *helper) getClientOptions() []be.BackendOption {
	co := h.gh.GetClientOptions()

	return []be.BackendOption{be.WithClientOptions(co)}
}
