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
	"net/http"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"l7e.io/yama"

	"l7e.io/vanity/cmd/vanity/cli/backends"
)

func serverCmd(cmd *cobra.Command, _ []string) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		glog.Exitf("Unable to bind viper to command line flags: %s", err)
	}

	svrHelp := newHelper(cmd)

	vanity := svrHelp.getHTTPServer(backends.Backend)
	if err != nil {
		glog.Exitf("Unable create vanity server: %s", err)
	}

	healthz := svrHelp.getHealthz(newHandlerCheck(backends.Backend, "healthz"))
	if err != nil {
		glog.Exitf("Unable create healthz server: %s", err)
	}

	readyz := svrHelp.getReadyz(newHandlerCheck(backends.Backend, "readyz"))
	if err != nil {
		glog.Exitf("Unable create readyz server: %s", err)
	}

	metrics := svrHelp.getMetrics()
	if err != nil {
		glog.Exitf("Unable create metrics server: %s", err)
	}

	watcher := yama.NewWatcher(
		yama.WatchingSignals(syscall.SIGINT, syscall.SIGTERM),
		yama.WithTimeout(2*time.Second), // nolint
		yama.WithClosers(backends.Backend, vanity, healthz, readyz, metrics))

	go func() {
		if err = metrics.ListenAndServe(); err != http.ErrServerClosed {
			glog.Error(err)
			_ = watcher.Close()
		}
	}()

	go func() {
		if err = healthz.ListenAndServe(); err != http.ErrServerClosed {
			glog.Error(err)
			_ = watcher.Close()
		}
	}()

	go func() {
		if err = readyz.ListenAndServe(); err != http.ErrServerClosed {
			glog.Error(err)
			_ = watcher.Close()
		}
	}()

	if err = vanity.ListenAndServe(); err != http.ErrServerClosed {
		glog.Error(err)
		_ = watcher.Close()
	}

	if err = watcher.Wait(); err != nil {
		glog.Warningf("Shutdown error: %s", err)
	}
	glog.Info("Vanity exited")
}
