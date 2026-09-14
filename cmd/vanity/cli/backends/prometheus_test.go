/*
 * Copyright (c) 2026 the original author or authors.
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

package backends

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"l7e.io/vanity"
)

var errBackend = errors.New("backend failed")

// recorder is a vanity.Backend that records each call. Every method of
// recorder returns errBackend.
type recorder struct {
	called     string
	importPath string
	vcs        string
	vcsPath    string
	consumer   vanity.Consumer
}

func (r *recorder) Close() error {
	r.called = "Close"

	return errBackend
}

func (r *recorder) Get(_ context.Context, importPath string) (vcs, vcsPath string, err error) {
	r.called = "Get"
	r.importPath = importPath

	return "git", "example.com/repo", errBackend
}

func (r *recorder) Add(_ context.Context, importPath, vcs, vcsPath string) error {
	r.called = "Add"
	r.importPath = importPath
	r.vcs = vcs
	r.vcsPath = vcsPath

	return errBackend
}

func (r *recorder) Remove(_ context.Context, importPath string) error {
	r.called = "Remove"
	r.importPath = importPath

	return errBackend
}

func (r *recorder) List(_ context.Context, consumer vanity.Consumer) error {
	r.called = "List"
	r.consumer = consumer

	return errBackend
}

func (r *recorder) Healthz(_ context.Context) error {
	r.called = "Healthz"

	return errBackend
}

func TestWrapWithPrometheus_Close(t *testing.T) {
	be := &recorder{}
	w := WrapWithPrometheus(be)

	err := w.Close()

	require.ErrorIs(t, err, errBackend)
	assert.Equal(t, "Close", be.called)
}

func TestWrapWithPrometheus_Get(t *testing.T) {
	be := &recorder{}
	w := WrapWithPrometheus(be)

	vcs, vcsPath, err := w.Get(t.Context(), "example.com/repo")

	require.ErrorIs(t, err, errBackend)
	assert.Equal(t, "Get", be.called)
	assert.Equal(t, "example.com/repo", be.importPath)
	assert.Equal(t, "git", vcs)
	assert.Equal(t, "example.com/repo", vcsPath)
}

func TestWrapWithPrometheus_Add(t *testing.T) {
	be := &recorder{}
	w := WrapWithPrometheus(be)

	err := w.Add(t.Context(), "example.com/repo", "git", "https://github.com/a/b")

	require.ErrorIs(t, err, errBackend)
	assert.Equal(t, "Add", be.called)
	assert.Equal(t, "example.com/repo", be.importPath)
	assert.Equal(t, "git", be.vcs)
	assert.Equal(t, "https://github.com/a/b", be.vcsPath)
}

func TestWrapWithPrometheus_Remove(t *testing.T) {
	be := &recorder{}
	w := WrapWithPrometheus(be)

	err := w.Remove(t.Context(), "example.com/repo")

	require.ErrorIs(t, err, errBackend)
	assert.Equal(t, "Remove", be.called)
	assert.Equal(t, "example.com/repo", be.importPath)
}

func TestWrapWithPrometheus_List(t *testing.T) {
	be := &recorder{}
	w := WrapWithPrometheus(be)

	consumer := vanity.ConsumerFunc(func(_ context.Context, _, _, _ string) {})

	err := w.List(t.Context(), consumer)

	require.ErrorIs(t, err, errBackend)
	assert.Equal(t, "List", be.called)
	assert.NotNil(t, be.consumer)
}

func TestWrapWithPrometheus_Healthz(t *testing.T) {
	be := &recorder{}
	w := WrapWithPrometheus(be)

	err := w.Healthz(t.Context())

	require.ErrorIs(t, err, errBackend)
	assert.Equal(t, "Healthz", be.called)
}
