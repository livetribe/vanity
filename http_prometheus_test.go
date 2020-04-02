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

package vanity_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"

	"l7e.io/vanity"
)

var snapshots = make(map[string]float64)

func prometheusReset() {
	prometheusSnapshot(vanity.APICalls)
	prometheusSnapshot(vanity.APIErrors)
	prometheusSnapshot(vanity.APINotFound)
	prometheusSnapshot(vanity.APIDocRedirects)
	prometheusSnapshot(vanity.APIErrTemplates)
}

func prometheusSnapshot(c prometheus.Counter) {
	m := &dto.Metric{}
	_ = c.Write(m)
	snapshots[c.Desc().String()] = m.Counter.GetValue()
}

func prometheusCheck(t *testing.T, calls, errors, notFound, docRedirects, errTemplates int) {
	prometheusCheckMetric(t, vanity.APICalls, calls)
	prometheusCheckMetric(t, vanity.APIErrors, errors)
	prometheusCheckMetric(t, vanity.APINotFound, notFound)
	prometheusCheckMetric(t, vanity.APIDocRedirects, docRedirects)
	prometheusCheckMetric(t, vanity.APIErrTemplates, errTemplates)
}

func prometheusCheckMetric(t *testing.T, c prometheus.Counter, v int) {
	m := &dto.Metric{}
	err := c.Write(m)
	assert.NoError(t, err)
	start := snapshots[c.Desc().String()]
	assert.Equal(t, float64(v), m.Counter.GetValue()-start, c.Desc().String())
}
