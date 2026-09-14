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

// Package server contains the server sub-command to serve vanity URLs.
package server

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"l7e.io/yama"

	"l7e.io/vanity/cmd/vanity/cli/backends"
)

// closeTimeout is the time that the watcher gives the closers to finish.
const closeTimeout = 2 * time.Second

func serverCmd(cmd *cobra.Command, _ []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("unable to bind viper to the command line flags: %w", err)
	}

	svrHelp := newHelper(cmd)

	be := backends.Get()
	be = backends.WrapWithPrometheus(be)

	vanity := svrHelp.getHTTPServer(be)

	healthz := svrHelp.getHealthz(newHandlerCheck(be, "healthz"))
	readyz := svrHelp.getReadyz(newHandlerCheck(be, "readyz"))

	metrics := svrHelp.getMetrics()

	watcher := yama.NewWatcher(
		yama.WatchingSignals(syscall.SIGINT, syscall.SIGTERM),
		yama.WithTimeout(closeTimeout),
		yama.WithClosers(be, vanity, healthz, readyz, metrics))

	go serve(metrics, watcher)
	go serve(healthz, watcher)
	go serve(readyz, watcher)

	serve(vanity, watcher)

	if err := watcher.Wait(); err != nil {
		slog.Warn("Unable to shut down", slog.Any("error", err))
	}

	slog.Info("Vanity exited")

	return nil
}

// serve runs the server until the server stops. It closes the watcher when the
// server stops with an error.
func serve(server *http.Server, watcher io.Closer) {
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Unable to serve", slog.Any("error", err))
		_ = watcher.Close()
	}
}
