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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"l7e.io/vanity/cmd/vanity/cli/backends/helpers"
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

// The names of the command line flags of the server sub-command.
const (
	bind    = "bind"
	port    = "port"
	healthz = "healthz"
	readyz  = "readyz"
	metrics = "prometheus"
)

// The ports that the server sub-command uses when the caller sets no others.
const (
	defaultPort        = 8080
	defaultHealthzPort = 8081
	defaultReadyzPort  = 8082
	defaultMetricsPort = 9100
)

// The names that identify each listener in the log records of the application.
const (
	vanityName  = "vanity"
	healthzName = "healthz"
	readyzName  = "readyz"
	metricsName = "metrics"
)

// Config holds the network configuration of the listeners of the server
// sub-command.
type Config struct {
	// Bind is the interface that every listener binds.
	Bind string

	// Port is the port that serves vanity URLs.
	Port int

	// HealthzPort is the port that serves the application health checks.
	HealthzPort int

	// ReadyzPort is the port that serves the application ready checks.
	ReadyzPort int

	// MetricsPort is the port that serves the Prometheus metrics.
	MetricsPort int
}

// NewConfig reads the network configuration of the listeners from the command
// line flags and the other sources that viper holds.
func NewConfig() Config {
	return Config{
		Bind:        viper.GetString(bind),
		Port:        viper.GetInt(port),
		HealthzPort: viper.GetInt(healthz),
		ReadyzPort:  viper.GetInt(readyz),
		MetricsPort: viper.GetInt(metrics),
	}
}

// addr returns the address that a listener binds for the given port.
func (c Config) addr(p int) string {
	return fmt.Sprintf("%s:%d", c.Bind, p)
}

// initFlags adds the command line flags of the server sub-command to cmd.
func initFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.Int16P(port, "p", defaultPort, "port on which the server will listen")
	flags.StringP(bind, "", "127.0.0.1", "interface to which the server will bind")
	flags.Int16P(healthz, "", defaultHealthzPort, "port on which application health checks will listen")
	flags.Int16P(readyz, "", defaultReadyzPort, "port on which application ready checks will listen")
	flags.Int16P(metrics, "", defaultMetricsPort, "port on which the Prometheus will listen")
}
