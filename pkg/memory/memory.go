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

// Package memory provides an in-memory implementation of Backend.
package memory // import "l7e.io/vanity/pkg/memory"

import (
	"context"
	"sync"

	"l7e.io/vanity"
)

// ConvenientBackend is the interface that wraps Backend with the AddEntry method.
type ConvenientBackend interface {
	vanity.Backend

	// AddEntry is a convenience method for adding a vanity URL configuration
	// without having to deal with errors.
	AddEntry(importPath, vcs, vcsPath string)
}
type inMemory struct {
	lock    sync.RWMutex
	entries map[string]*entry
	closed  bool
}

type entry struct {
	vcs, vcsPath string
}

// NewInMemoryAPI creates an in-memory Backend instance.
func NewInMemoryAPI() ConvenientBackend {
	return &inMemory{entries: make(map[string]*entry)}
}

func (s *inMemory) AddEntry(importPath, vcs, vcsPath string) {
	s.entries[importPath] = &entry{vcs: vcs, vcsPath: vcsPath}
}

func (s *inMemory) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.closed = true

	return nil
}

func (s *inMemory) check() error {
	if s.closed {
		return vanity.ErrAlreadyClosed
	}
	return nil
}

func (s *inMemory) Get(_ context.Context, importPath string) (string, string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if err := s.check(); err != nil {
		return "", "", err
	}

	e, found := s.entries[importPath]
	if !found {
		return "", "", vanity.ErrNotFound
	}
	return e.vcs, e.vcsPath, nil
}

func (s *inMemory) Add(_ context.Context, importPath, vcs, vcsPath string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.check(); err != nil {
		return err
	}

	s.AddEntry(importPath, vcs, vcsPath)

	return nil
}

func (s *inMemory) Remove(_ context.Context, importPath string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.check(); err != nil {
		return err
	}

	_, found := s.entries[importPath]
	if !found {
		return vanity.ErrNotFound
	}

	delete(s.entries, importPath)

	return nil
}

func (s *inMemory) List(ctx context.Context, consumer vanity.Consumer) error {
	s.lock.RLock()

	if err := s.check(); err != nil {
		return err
	}

	c := make(map[string]*entry)
	for k, v := range s.entries {
		c[k] = v
	}

	s.lock.RUnlock()

	for importPath, e := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		consumer.OnEntry(ctx, importPath, e.vcs, e.vcsPath)
	}

	return nil
}

func (s *inMemory) Healthz(_ context.Context) error {
	return nil
}
