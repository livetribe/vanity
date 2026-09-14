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
	"bytes"
	"context"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	xForwardedHost = "X-Forwarded-Host"

	// DefaultDocURL is the documentation site that a Handler uses when the
	// caller does not set Handler.DocURL.
	DefaultDocURL = "https://pkg.go.dev/"

	// DefaultDuration is the timeout that a Handler uses for a Backend call
	// when the caller does not set Handler.Duration.
	DefaultDuration = 5 * time.Second

	metricNamespace = "vanity"
	metricSubsystem = "api"

	contentTypeHeader = "Content-Type"
	htmlContentType   = "text/html; charset=utf-8"
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
  <meta name="go-source" content="{{.ImportRoot}} {{.VCSRoot}} ` +
	`{{.VCSRoot}}/tree/master{/dir} {{.VCSRoot}}/blob/master{/dir}/{file}#L{line}">
</head>
</html>
`))

var (
	// APICalls is a Prometheus counter that tracks the total vanity Backend calls.
	APICalls = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "calls_total",
		Help:      "The total vanity Backend calls",
	})

	// APIErrors is a Prometheus counter that tracks the total vanity Backend errors.
	APIErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "errors_total",
		Help:      "The total vanity Backend errors",
	})

	// APINotFound is a Prometheus counter that tracks the total vanity Backend not found calls.
	APINotFound = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "not_found_total",
		Help:      "The total vanity Backend not found calls",
	})

	// APIDocRedirects is a Prometheus counter that tracks the total vanity Backend doc redirects.
	APIDocRedirects = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "doc_total",
		Help:      "The total vanity Backend doc redirects",
	})

	// APIErrTemplates is a Prometheus counter that tracks the total templating errors.
	APIErrTemplates = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "error_templates_total",
		Help:      "The total templating errors",
	})

	// SummaryVec is a Prometheus histogram to track the Backend duration in seconds.
	SummaryVec = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "duration_seconds",
		Help:      "The Backend duration in seconds",
	})
)

// Handler is a http.Handler that services vanity URLs using api
// as a backend service.
type Handler struct {
	api Backend

	// DocURL is the base URL that the Handler redirects a browser to.
	// The Handler uses DefaultDocURL when DocURL is empty.
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
		Duration: DefaultDuration,
	}
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(status), status)

		return
	}

	APICalls.Inc()

	ctx, cancel := context.WithTimeout(r.Context(), s.Duration)
	defer cancel()

	root := r.URL.Path
	paths := strings.FieldsFunc(root, func(c rune) bool { return c == '/' })
	if len(paths) > 0 {
		root = "/" + paths[0]
	}

	importPath := host(r) + root

	vcs, repoRoot, err := s.timedGet(ctx, importPath)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			APINotFound.Inc()
			http.NotFound(w, r)
		} else {
			logger.Printf("Unable to get %s: %s", importPath, err)
			APIErrors.Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	vcsRoot := repoRoot + root

	if r.FormValue("go-get") != "1" {
		APIDocRedirects.Inc()
		url := s.docURL() + importPath
		http.Redirect(w, r, url, http.StatusTemporaryRedirect) //nolint:gosec

		return
	}

	body, err := templatize(host(r)+r.URL.Path, vcs, vcsRoot)
	if err != nil {
		logger.Printf("Unable to templatize %s: %s", importPath, err)
		APIErrTemplates.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(contentTypeHeader, htmlContentType)
	w.Header().Set("Cache-Control", "public, max-age=300")

	_, err = w.Write(body) //nolint:gosec
	if err != nil {
		logger.Printf("Error writing body for %s: %s", importPath, err)
	}
}

func (s *Handler) docURL() string {
	if s.DocURL == "" {
		return DefaultDocURL
	}

	return s.DocURL
}

func (s *Handler) timedGet(ctx context.Context, importPath string) (vcs, vcsPath string, err error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		SummaryVec.Observe(elapsed.Seconds())
	}()

	vcs, vcsPath, err = s.api.Get(ctx, importPath)

	return
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
