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

package backends

import (
	"context"
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"l7e.io/vanity"
)

const (
	metricNamespace = "vanity"
	metricSubsystem = "api"
)

var (
	// durationHistogram records the duration of each Backend.Get call in seconds.
	durationHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "duration_seconds",
		Help:      "The Backend duration in seconds",
	})

	// callCounter counts each Backend.Get call.
	callCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "calls_total",
		Help:      "The total vanity Backend calls",
	})

	// errorCounter counts each Backend.Get call that fails for a reason other
	// than vanity.ErrNotFound.
	errorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "errors_total",
		Help:      "The total vanity Backend errors",
	})

	// notFoundCounter counts each Backend.Get call that fails with
	// vanity.ErrNotFound.
	notFoundCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "not_found_total",
		Help:      "The total vanity Backend not found calls",
	})
)

type prometheusWrapper struct {
	backend vanity.Backend
}

var _ vanity.Backend = (*prometheusWrapper)(nil)

// WrapWithPrometheus returns a vanity.Backend that sends all calls to backend.
// The returned Backend also records the duration of each Get call in the
// vanity_api_duration_seconds histogram, and counts each Get call in the
// vanity_api_calls_total, vanity_api_errors_total and
// vanity_api_not_found_total counters.
func WrapWithPrometheus(backend vanity.Backend) vanity.Backend {
	return &prometheusWrapper{backend: backend}
}

func (w *prometheusWrapper) Close() error {
	return w.backend.Close()
}

func (w *prometheusWrapper) Get(ctx context.Context, importPath string) (vcs, vcsPath string, err error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		durationHistogram.Observe(elapsed.Seconds())
	}()

	callCounter.Inc()

	vcs, vcsPath, err = w.backend.Get(ctx, importPath)
	if err != nil {
		if errors.Is(err, vanity.ErrNotFound) {
			notFoundCounter.Inc()
		} else {
			errorCounter.Inc()
		}
	}

	return vcs, vcsPath, err
}

func (w *prometheusWrapper) Add(ctx context.Context, importPath, vcs, vcsPath string) error {
	return w.backend.Add(ctx, importPath, vcs, vcsPath)
}

func (w *prometheusWrapper) Remove(ctx context.Context, importPath string) error {
	return w.backend.Remove(ctx, importPath)
}

func (w *prometheusWrapper) List(ctx context.Context, consumer vanity.Consumer) error {
	return w.backend.List(ctx, consumer)
}

func (w *prometheusWrapper) Healthz(ctx context.Context) error {
	return w.backend.Healthz(ctx)
}
