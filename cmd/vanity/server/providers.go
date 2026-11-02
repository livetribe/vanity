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
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"l7e.io/vanity"
	"l7e.io/vanity/cmd/vanity/cli/backends"
	"l7e.io/vanity/cmd/vanity/server/interceptors"
	"l7e.io/vanity/internal/mw"
)

// GraphSet builds the component graph of the server sub-command.
var GraphSet = wire.NewSet(
	NewBackend,
	NewVanity,
	NewHealthz,
	NewReadyz,
	NewMetrics,
	wire.Struct(new(App), "*"))

// App holds the listeners of the server sub-command.
type App struct {
	Vanity  *Vanity
	Healthz *Healthz
	Readyz  *Readyz
	Metrics *Metrics
}

// Vanity is the listener that serves vanity URLs.
type Vanity struct {
	*listener
}

// NewVanity provides the listener that serves vanity URLs.
func NewVanity(cfg Config, be vanity.Backend, fail Fail) *Vanity {
	vanityHandler := vanity.NewVanityHandler(be)
	meteredHandler := mw.WithPrometheus(vanityHandler)
	loggedHandler := mw.WithLogger(meteredHandler)
	handler := interceptors.WrapHandler(loggedHandler)

	addr := cfg.addr(cfg.Port)

	return &Vanity{newListener(vanityName, addr, "/", handler, fail)}
}

// Healthz is the listener that serves the application health checks.
type Healthz struct {
	*listener
}

// NewHealthz provides the listener that serves the application health checks.
func NewHealthz(cfg Config, be vanity.Backend, fail Fail) *Healthz {
	handler := newHandlerCheck(be, healthzName)

	addr := cfg.addr(cfg.HealthzPort)

	return &Healthz{newListener(healthzName, addr, "/healthz", handler, fail)}
}

// Readyz is the listener that serves the application ready checks.
type Readyz struct {
	*listener
}

// NewReadyz provides the listener that serves the application ready checks.
func NewReadyz(cfg Config, be vanity.Backend, fail Fail) *Readyz {
	handler := newHandlerCheck(be, readyzName)

	addr := cfg.addr(cfg.ReadyzPort)

	return &Readyz{newListener(readyzName, addr, "/readyz", handler, fail)}
}

// Metrics is the listener that serves the Prometheus metrics.
type Metrics struct {
	*listener
}

// NewMetrics provides the listener that serves the Prometheus metrics.
func NewMetrics(cfg Config, fail Fail) *Metrics {
	handler := promhttp.Handler()

	addr := cfg.addr(cfg.MetricsPort)

	return &Metrics{newListener(metricsName, addr, "/metrics", handler, fail)}
}

// NewBackend provides the vanity URL store that the sub-command installed. The
// store counts the calls that it receives.
//
//yama:closer
func NewBackend() vanity.Backend {
	be := backends.Get()

	return backends.WrapWithPrometheus(be)
}
