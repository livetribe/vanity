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
	_ = viper.BindPFlag(database, flags.Lookup(database))
	_ = viper.BindEnv(database)

	gcp.InitFlags(Command)

	viper.RegisterAlias(database, "spanner.database")
}

// Command is the vanity sub-command for a GCP Spanner backend.
var Command = &cobra.Command{
	Use:   "spanner",
	Short: "Use GCP Spanner backend for a vanity store",
	Long:  "Use GCP Spanner backend for a vanity store",
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

		err = beHelp.gh.AddInterceptors()
		if err != nil {
			return err
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
	*cli.FlagSet

	gh *gcp.Helper
}

// newHelper wraps the Cobra command's flags with a utility wrapper to assist in
// the creation of a Spanner-based backend.
func newHelper(cmd *cobra.Command) *helper {
	return &helper{FlagSet: cli.Flags(cmd), gh: gcp.NewHelper(cmd)}
}

// getBackend returns a Spanner-based api.Backend instance, configured by the helper.
func (h *helper) getBackend() (vanity.Backend, error) {
	db, ok := h.GetValue(database)
	if !ok {
		return nil, fmt.Errorf("unable to get database")
	}
	glog.V(log.Debug).Infof("database: %s", db)

	options, err := h.getClientOptions()
	if err != nil {
		return nil, err
	}

	return be.NewClient(context.Background(), db, options...)
}

func (h *helper) getClientOptions() ([]be.BackendOption, error) {
	co, err := h.gh.GetClientOptions()
	if err != nil {
		return nil, err
	}

	return []be.BackendOption{be.WithClientOptions(co)}, nil
}
