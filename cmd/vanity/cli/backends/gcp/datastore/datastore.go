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

// Package datastore contains the datastore sub-command.
package datastore // import "l7e.io/vanity/cmd/vanity/backends/gcp/datastore"

import (
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
	be "l7e.io/vanity/pkg/gcp/datastore"
)

const (
	projectID = "project-id"
)

func init() { //nolint:gochecknoinits
	cli.RootCmd.AddCommand(Command)

	flags := Command.PersistentFlags()
	flags.StringP(projectID, "", "", "GCP project hosting datastore")

	gcp.InitFlags(Command)

	viper.RegisterAlias(projectID, "datastore.project-id")
}

// Command is the vanity sub-command for a GCP Datastore backend.
var Command = &cobra.Command{
	Use:   "datastore",
	Short: "Use GCP Datastore for a backend vanity store",
	Long:  "Use GCP Datastore for a backend vanity store",
	Args:  cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		glog.V(log.Debug).Infoln("Set backend w/ datastore")

		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("unable to bind viper to Command line flags: %w", err)
		}

		beHelp := newHelper(cmd)
		backend, err := beHelp.getBackend()
		if err != nil {
			return fmt.Errorf("unable to instantiate Datastore backend: %w", err)
		}

		backends.Backend = backend
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		glog.V(log.Debug).Infoln("Clean backend of datastore")
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
// the creation of a Datastore-based backend.
func newHelper(cmd *cobra.Command) *helper {
	return &helper{FlagSet: cmd.Flags(), gh: gcp.NewHelper(cmd)}
}

// getProjectID obtains the GCP project hosting datastore.
func (h *helper) getProjectID() (id string, found bool) {
	id = viper.GetString(projectID)
	return id, id != ""
}

// getBackend returns a Datastore-based api.Backend instance, configured by the
// helper.
func (h *helper) getBackend() (vanity.Backend, error) {
	id, found := h.getProjectID()
	if !found {
		glog.Info("Unable to obtain project id - relying on DATASTORE_PROJECT_ID")
	} else {
		glog.Infof("project id: %s", id)
	}

	return be.NewClient(id, h.gh.GetClientOptions()...)
}
