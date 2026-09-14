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

// Package log configures the structured logger of the vanity command.
package log

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"m4o.io/gslog/core"
	"m4o.io/gslog/stdout"
)

const (
	// LevelFlag is the name of the command line flag that sets the log level.
	LevelFlag = "log-level"

	// DefaultLevel is the log level that the command uses when the caller sets
	// no other level.
	DefaultLevel = "info"

	// Trace is the log level below Debug. The handler gives the entries of
	// this level the severity DEFAULT.
	Trace = slog.Level(-8)
)

// ErrUnknownLevel is the error that SetLevel returns for a name that is not a
// log level.
var ErrUnknownLevel = errors.New("unknown log level")

// level holds the log level of the default logger.
var level = new(slog.LevelVar)

// levels maps the name of a log level to its slog.Level.
var levels = map[string]slog.Level{
	"trace": Trace,
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// Init makes a Cloud Logging handler that writes to w. Init sets this handler
// as the default slog handler. The handler writes each entry as one line of
// JSON in the structured logging format that a Google logging agent reads.
func Init(w io.Writer) {
	handler := stdout.NewHandler(w, core.WithLogLeveler(level))

	slog.SetDefault(slog.New(handler))
}

// SetLevel sets the log level of the default logger. SetLevel accepts the
// names trace, debug, info, warn, and error. The name is not case-sensitive.
func SetLevel(name string) error {
	lowercase := strings.ToLower(name)

	l, found := levels[lowercase]
	if !found {
		return fmt.Errorf("%w: %s", ErrUnknownLevel, name)
	}
	level.Set(l)

	return nil
}
