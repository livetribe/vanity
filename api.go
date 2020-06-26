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

package vanity

import (
	"context"
	"fmt"
	"io"
)

var (
	// ErrAlreadyClosed is returned if a Backend implementation is already closed.
	ErrAlreadyClosed = fmt.Errorf("already closed")

	// ErrNotFound is returned if the import path cannot MockBackend found.
	ErrNotFound = fmt.Errorf("not found")

	// ErrNotSupported is returned if the Backend method is not supported by the implementation.
	ErrNotSupported = fmt.Errorf("not supported")
)

// Backend implementations provide access to a vanity URL store.
//
// Sub-commands, such as add, list, server, use this interface to perform their
// functionality.
type Backend interface {
	io.Closer

	// Get vanity URL configuration for a given import path
	Get(ctx context.Context, importPath string) (vcs, vcsPath string, err error)

	// Add a vanity URL configuration
	Add(ctx context.Context, importPath, vcs, vcsPath string) error

	// Remove a vanity URL configuration by it's key, the import path
	Remove(ctx context.Context, importPath string) error

	// List all registered URL configurations, delivering them to the consumer callback
	List(ctx context.Context, consumer Consumer) error

	// Healthz is a health check point for Kubernetes
	Healthz(ctx context.Context) error
}

// Consumer is the interface whose implementations are provided to the
// Backend.List() method which calls their OnEntry method with the vanity
// entries found.
type Consumer interface {
	OnEntry(context context.Context, importPath, vcs, vcsPath string)
}

// The ConsumerFunc type is an adapter to allow the use of
// ordinary functions as consumers. If f is a function
// with the appropriate signature, ConsumerFunc(f) is a
// Consumer that calls f.
type ConsumerFunc func(context context.Context, importPath, vcs, vcsPath string)

// OnEntry calls f(w, r).
func (f ConsumerFunc) OnEntry(context context.Context, importPath, vcs, vcsPath string) {
	f(context, importPath, vcs, vcsPath)
}

var (
	// log is error log.
	logger Logger
)

func init() {
	logger = LoggerFunc(func(format string, v ...interface{}) {})
}

// Logger describes functions available for logging purposes.
type Logger interface {
	Printf(format string, v ...interface{})
}

// SetLogger sets the logger used by vanity package's error log.
func SetLogger(l Logger) {
	logger = l
}

// The LoggerFunc type is an adapter to allow the use of
// ordinary functions as Loggers. If f is a function
// with the appropriate signature, LoggerFunc(f) is a
// Logger that calls f.
type LoggerFunc func(string, ...interface{})

// Printf calls f(w, r).
func (f LoggerFunc) Printf(format string, v ...interface{}) {
	f(format, v...)
}
