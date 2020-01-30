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
 *
 */

package server

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/glog"
	"l7e.io/vanity"
	"l7e.io/vanity/cmd/vanity/cli/log"
)

// newHandlerCheck creates a handler instance that can be used for a healthz checkpoint.
func newHandlerCheck(backend vanity.Backend, kind string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx, _ := context.WithDeadline(r.Context(), time.Now().Add(time.Second)) // nolint
			err := backend.Healthz(ctx)
			if err != nil {
				glog.Error(err)
				http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			} else {
				glog.V(log.Trace).Infof("%s-check...", kind)
				w.WriteHeader(http.StatusOK)
			}
		})
}
