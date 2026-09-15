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

package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"l7e.io/vanity"
	"l7e.io/vanity/cmd/vanity/cli/backends/helpers"
	"l7e.io/vanity/cmd/vanity/server/interceptors"
	"l7e.io/vanity/internal/mw"
)

func init() { //nolint:gochecknoinits
	helpers.AddCommand(func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "server",
			Short: "Serve port URLs",
			Long:  "Serve port URLs using [command] for a backend port store",
			Args:  cobra.NoArgs,
			RunE:  serverCmd,
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

	defaultPort        = 8080
	defaultHealthzPort = 8081
	defaultReadyzPort  = 8082
	defaultMetricsPort = 9100

	// readHeaderTimeout bounds the time that a client can take to send the
	// request headers.
	readHeaderTimeout = 10 * time.Second
)

func initFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.Int16P(port, "p", defaultPort, "port on which the server will listen")
	flags.StringP(bind, "", "127.0.0.1", "interface to which the server will bind")
	flags.Int16P(healthz, "", defaultHealthzPort, "port on which application health checks will listen")
	flags.Int16P(readyz, "", defaultReadyzPort, "port on which application ready checks will listen")
	flags.Int16P(metrics, "", defaultMetricsPort, "port on which the Prometheus will listen")
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

	slog.Info("Configured the vanity port", slog.String("address", addr))

	vanityHandler := vanity.NewVanityHandler(api)
	meteredHandler := mw.WithPrometheus(vanityHandler)
	handler := mw.WithLogger(meteredHandler)

	mux := http.NewServeMux()
	mux.Handle("/", interceptors.WrapHandler(handler))

	return &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: readHeaderTimeout}
}

// getHealthz returns an http.Server for healthz configured by the helper.
func (h *helper) getHealthz(handler http.Handler) *http.Server {
	port := viper.GetInt(healthz)
	nic := viper.GetString(bind)

	addr := fmt.Sprintf("%s:%d", nic, port)

	slog.Info("Configured the healthz port", slog.String("address", addr+"/healthz"))

	mux := http.NewServeMux()
	mux.Handle("/healthz", handler)

	return &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: readHeaderTimeout}
}

// getReadyz returns an http.Server for readyz configured by the helper.
func (h *helper) getReadyz(handler http.Handler) *http.Server {
	port := viper.GetInt(readyz)
	nic := viper.GetString(bind)

	addr := fmt.Sprintf("%s:%d", nic, port)

	slog.Info("Configured the readyz port", slog.String("address", addr+"/readyz"))

	mux := http.NewServeMux()
	mux.Handle("/readyz", handler)

	return &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: readHeaderTimeout}
}

// getMetrics returns an http.Server for readyz configured by the helper.
func (h *helper) getMetrics() *http.Server {
	port := viper.GetInt(metrics)
	nic := viper.GetString(bind)

	addr := fmt.Sprintf("%s:%d", nic, port)

	slog.Info("Configured the metrics port", slog.String("address", addr+"/metrics"))

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: readHeaderTimeout}
}
