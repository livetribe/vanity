/*
 * Copyright (c) 2026 the original author or authors.
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
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"l7e.io/vanity"
)

const (
	metricNamespace = "vanity"
	metricSubsystem = "api"
)

// durationHistogram records the duration of each Backend.Get call in seconds.
var durationHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Namespace: metricNamespace,
	Subsystem: metricSubsystem,
	Name:      "duration_seconds",
	Help:      "The Backend duration in seconds",
})

type prometheusWrapper struct {
	backend vanity.Backend
}

var _ vanity.Backend = (*prometheusWrapper)(nil)

// WrapWithPrometheus returns a vanity.Backend that sends all calls to backend.
// The returned Backend also records the duration of each Get call in the
// vanity_api_duration_seconds histogram.
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

	return w.backend.Get(ctx, importPath)
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
