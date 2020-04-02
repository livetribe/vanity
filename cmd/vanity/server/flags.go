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

package server

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"l7e.io/vanity"
	"l7e.io/vanity/cmd/vanity/cli/backends/helpers"
	"l7e.io/vanity/cmd/vanity/server/interceptors"
)

func init() { //nolint:gochecknoinits
	helpers.AddCommand(func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "server",
			Short: "Serve port URLs",
			Long:  "Serve port URLs using [command] for a backend port store",
			Args:  cobra.NoArgs,
			Run:   serverCmd,
		}

		initFlags(cmd)

		return cmd
	})
}

const (
	bind    = "bind"
	port    = "port"
	healthz = "healthz"
	readyz  = "readyz"
	metrics = "prometheus"
)

func initFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.Int16P(port, "p", 8080, "port on which the server will listen")
	flags.StringP(bind, "", "127.0.0.1", "interface to which the server will bind")
	flags.Int16P(healthz, "", 8081, "port on which application health checks will listen")
	flags.Int16P(readyz, "", 8082, "port on which application ready checks will listen")
	flags.Int16P(metrics, "", 9100, "port on which the Prometheus will listen")
}

type helper struct {
	*pflag.FlagSet
}

func newHelper(cmd *cobra.Command) *helper {
	return &helper{cmd.Flags()}
}

// getHTTPServer returns an http.Server configured by the helper.
func (h *helper) getHTTPServer(api vanity.Backend) *http.Server {
	port := viper.GetInt(port)
	nic := viper.GetString(bind)

	addr := fmt.Sprintf("%s:%d", nic, port)

	glog.Infof("port configured to listen to %s", addr)

	mux := http.NewServeMux()
	mux.Handle("/", interceptors.WrapHandler(vanity.NewVanityHandler(api)))

	return &http.Server{Addr: addr, Handler: mux}
}

// getHealthz returns an http.Server for healthz configured by the helper.
func (h *helper) getHealthz(handler http.Handler) *http.Server {
	port := viper.GetInt(healthz)
	nic := viper.GetString(bind)

	addr := fmt.Sprintf("%s:%d", nic, port)

	glog.Infof("healthz configured to listen to %s/healthz", addr)

	mux := http.NewServeMux()
	mux.Handle("/healthz", handler)

	return &http.Server{Addr: addr, Handler: mux}
}

// getReadyz returns an http.Server for readyz configured by the helper.
func (h *helper) getReadyz(handler http.Handler) *http.Server {
	port := viper.GetInt(readyz)
	nic := viper.GetString(bind)

	addr := fmt.Sprintf("%s:%d", nic, port)

	glog.Infof("readyz configured to listen to %s/readyz", addr)

	mux := http.NewServeMux()
	mux.Handle("/readyz", handler)

	return &http.Server{Addr: addr, Handler: mux}
}

// getMetrics returns an http.Server for readyz configured by the helper.
func (h *helper) getMetrics() *http.Server {
	port := viper.GetInt(metrics)
	nic := viper.GetString(bind)

	addr := fmt.Sprintf("%s:%d", nic, port)

	glog.Infof("metrics configured to listen to %s/metrics", addr)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{Addr: addr, Handler: mux}
}
