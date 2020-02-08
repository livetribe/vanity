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

package vanity // import "l7e.io/vanity"

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"l7e.io/vanity/cmd/vanity/cli/log"
)

const (
	gceIngressUserAgent = "GoogleHC/1.0"
	xForwardedHost      = "X-Forwarded-Host"
	xForwardedProto     = "X-Forwarded-Proto"

	// DefaultDocURL is the default Go doc URL.
	// It can be replaced with https://https://godoc.org/.
	DefaultDocURL = "https://pkg.go.dev/"
)

type data struct {
	ImportRoot string
	VCS        string
	VCSRoot    string
}

var tmpl = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
  <meta name="go-import" content="{{.ImportRoot}} {{.VCS}} {{.VCSRoot}}">
  <meta name="go-source" content="{{.ImportRoot}} {{.VCSRoot}} {{.VCSRoot}}/tree/master{/dir} {{.VCSRoot}}/blob/master{/dir}/{file}#L{line}">
</head>
</html>
`))

var (
	apiCalls = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vanity_api_calls_total",
		Help: "The total vanity Backend calls",
	})

	apiErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vanity_api_errors_total",
		Help: "The total vanity Backend errors",
	})

	apiNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vanity_api_not_found_total",
		Help: "The total vanity Backend not found calls",
	})

	apiDocRedirects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vanity_api_doc_total",
		Help: "The total vanity Backend doc redirects",
	})

	apiErrTemplates = promauto.NewCounter(prometheus.CounterOpts{
		Name: "vanity_api_error_templates_total",
		Help: "The total templating errors",
	})

	summaryVec = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "vanity_api_duration_seconds",
		Help: "The Backend duration in seconds",
	})
)

// Handler is a http.Handler that services vanity URLs using api
// as a backend service.
type Handler struct {
	api Backend

	// DocURL is the base URL browsers are redirected to - default is https://pkg.go.dev/
	DocURL string

	// Duration is the timeout duration for the calls to the backend implementation.
	// Default is five seconds.
	Duration time.Duration
}

// NewVanityHandler creates a new http.Handler that services vanity URLs using api
// as a backend service.
func NewVanityHandler(api Backend) http.Handler {
	return &Handler{
		api:      api,
		DocURL:   DefaultDocURL,
		Duration: 5 * time.Second, // nolint
	}
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.V(log.Trace).Infof("User-Agent: %s", r.UserAgent())

	// tell Google LB everything is fine
	if r.UserAgent() == gceIngressUserAgent {
		w.WriteHeader(http.StatusOK)

		return
	}

	if !isHTTPS(r) {
		url := *r.URL
		url.Scheme = "https"
		url.Host = host(r)
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)

		glog.V(log.Trace).Info("Redirecting to https")

		return
	}

	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)

		glog.V(log.Trace).Infof("%s HTTP method not allowed", r.Method)

		return
	}

	apiCalls.Inc()

	ctx, _ := context.WithTimeout(r.Context(), s.Duration) // nolint

	importPath := host(r) + r.URL.Path

	glog.V(log.Trace).Infof("importPath: %s", importPath)

	vcs, repoRoot, err := s.timedGet(ctx, importPath)
	if err != nil {
		if err == ErrNotFound {
			glog.Warningf("%s not found", importPath)
			apiNotFound.Inc()
			http.NotFound(w, r)
		} else {
			logger.Printf("Unable to get %s: %s", importPath, err)
			apiErrors.Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	vcsRoot := repoRoot + r.URL.Path

	glog.V(log.Trace).Infof("vcsRoot: %s", vcsRoot)

	if r.FormValue("go-get") != "1" {
		apiDocRedirects.Inc()
		url := "https://pkg.go.dev/" + importPath
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

		return
	}

	body, err := templatize(importPath, vcs, vcsRoot)
	if err != nil {
		logger.Printf("Unable to templatize %s: %s", importPath, err)
		apiErrTemplates.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Cache-Control", "public, max-age=300")

	_, err = w.Write(body)
	if err != nil {
		logger.Printf("Error writing body for %s: %s", importPath, err)
	}
}

func (s *Handler) timedGet(ctx context.Context, importPath string) (vcs, vcsPath string, err error) {
	start := time.Now()
	defer summaryVec.Observe(time.Since(start).Seconds())

	vcs, vcsPath, err = s.api.Get(ctx, importPath)

	return
}

func isHTTPS(r *http.Request) bool {
	if r.URL.Scheme == "https" {
		return true
	}
	if strings.HasPrefix(r.Proto, "HTTPS") {
		return true
	}
	if r.Header.Get(xForwardedProto) == "https" {
		return true
	}
	if r.TLS != nil {
		return true
	}

	return false
}

func host(r *http.Request) string {
	if r.Header.Get(xForwardedHost) == "" {
		return r.Host
	}

	return r.Header.Get(xForwardedHost)
}

func templatize(importRoot, vcs, vcsRoot string) (body []byte, err error) {
	logger.Printf("%s %s %s", importRoot, vcs, vcsRoot)
	d := &data{
		ImportRoot: importRoot,
		VCS:        vcs,
		VCSRoot:    vcsRoot,
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, d)
	if err != nil {
		return
	}

	body = buf.Bytes()
	return
}
