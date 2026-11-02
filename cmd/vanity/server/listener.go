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
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

const (
	// readHeaderTimeout bounds the time that a client can take to send the
	// request headers.
	readHeaderTimeout = 10 * time.Second

	// drainTimeout bounds the time that a listener takes to finish the requests
	// in flight. Stop closes the remaining connections after this time.
	drainTimeout = 2 * time.Second
)

// Fail stops the application. A listener calls it when the listener stops
// serving with an error.
type Fail func()

// listener serves one HTTP endpoint as a lifecycle component. Start binds the
// address and serves it. Quiesce drains the requests in flight. Stop closes the
// connections that the drain did not finish.
type listener struct {
	name   string
	path   string
	server *http.Server
	fail   Fail
}

// newListener makes a listener called name that serves handler at addr under
// path. The listener calls fail if it stops serving with an error.
func newListener(name, addr, path string, handler http.Handler, fail Fail) *listener {
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	server := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: readHeaderTimeout}

	return &listener{name: name, path: path, server: server, fail: fail}
}

// Start binds the address of the listener and serves it in a new goroutine.
// Start returns an error if the bind fails. Start also logs that error, because
// the lifecycle reports only that a start failed.
func (l *listener) Start(ctx context.Context) error {
	var config net.ListenConfig

	ln, err := config.Listen(ctx, "tcp", l.server.Addr)
	if err != nil {
		slog.ErrorContext(ctx, "Unable to bind a listener", slog.String("name", l.name),
			slog.String("address", l.server.Addr), slog.String("path", l.path), slog.Any("error", err))

		return fmt.Errorf("unable to bind the %s listener to %s: %w", l.name, l.server.Addr, err)
	}

	addr := ln.Addr().String()
	slog.InfoContext(ctx, "Started a listener", slog.String("name", l.name),
		slog.String("address", addr), slog.String("path", l.path))

	go l.serve(ln)

	return nil
}

// Quiesce stops the intake of new requests and waits for the requests in
// flight. Quiesce gives the requests in flight the drain timeout to finish. The
// drain gets that time even if ctx is already done.
func (l *listener) Quiesce(ctx context.Context) {
	drainCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), drainTimeout)
	defer cancel()

	if err := l.server.Shutdown(drainCtx); err != nil {
		slog.WarnContext(ctx, "Unable to drain a listener",
			slog.String("name", l.name), slog.Any("error", err))
	}
}

// Stop closes the connections that the drain did not finish.
func (l *listener) Stop(ctx context.Context) {
	if err := l.server.Close(); err != nil {
		slog.WarnContext(ctx, "Unable to close a listener",
			slog.String("name", l.name), slog.Any("error", err))
	}
}

// serve serves ln until the server stops. serve stops the application if the
// server stops with an error.
func (l *listener) serve(ln net.Listener) {
	err := l.server.Serve(ln)
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return
	}

	slog.Error("Unable to serve", slog.String("name", l.name), slog.Any("error", err))
	l.fail()
}
