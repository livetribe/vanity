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

package log

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
	}{
		{name: "trace", level: Trace},
		{name: "debug", level: slog.LevelDebug},
		{name: "info", level: slog.LevelInfo},
		{name: "WARN", level: slog.LevelWarn},
		{name: "Error", level: slog.LevelError},
	}

	t.Cleanup(func() {
		level.Set(slog.LevelInfo)
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.NoError(t, SetLevel(test.name))
			assert.Equal(t, test.level, level.Level())
		})
	}
}

func TestSetLevel_unknown(t *testing.T) {
	err := SetLevel("bogus")

	require.ErrorIs(t, err, ErrUnknownLevel)
	assert.Equal(t, slog.LevelInfo, level.Level())
}

func TestInit(t *testing.T) {
	var buf bytes.Buffer

	previous := slog.Default()

	Init(&buf)

	require.NoError(t, SetLevel("trace"))

	t.Cleanup(func() {
		slog.SetDefault(previous)
		level.Set(slog.LevelInfo)
	})

	slog.Log(t.Context(), Trace, "how now brown cow")

	assert.Contains(t, buf.String(), `"message":"how now brown cow"`)
	assert.Contains(t, buf.String(), `"severity":"DEFAULT"`)
}
