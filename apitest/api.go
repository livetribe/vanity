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

// Package apitest contains a helpful mock implementation of Backend.
package apitest

import (
	"context"

	"l7e.io/vanity"
)

// MockBackend is a simple test mock of Backend.
type MockBackend struct {
	Healthy error
	Urls    map[string][]string
}

// Close implements the io.Closer interface.
func (b *MockBackend) Close() error {
	return nil
}

// Get vanity URL configuration for a given import path
func (b *MockBackend) Get(ctx context.Context, importPath string) (vcs, vcsPath string, err error) {
	if b.Healthy != nil {
		return "", "", b.Healthy
	}
	v, ok := b.Urls[importPath]
	if !ok {
		return "", "", vanity.ErrNotFound
	}
	return v[0], v[1], nil
}

// Add a vanity URL configuration
func (b *MockBackend) Add(ctx context.Context, importPath, vcs, vcsPath string) error {
	if b.Healthy != nil {
		return b.Healthy
	}
	b.Urls[importPath] = []string{vcs, vcsPath}
	return nil
}

// Remove a vanity URL configuration by it's key, the import path
func (b *MockBackend) Remove(ctx context.Context, importPath string) error {
	if b.Healthy != nil {
		return b.Healthy
	}
	_, ok := b.Urls[importPath]
	if !ok {
		return vanity.ErrNotFound
	}
	delete(b.Urls, importPath)
	return nil
}

// List all registered URL configurations, delivering them to the consumer callback
func (b *MockBackend) List(ctx context.Context, consumer vanity.Consumer) error {
	if b.Healthy != nil {
		return b.Healthy
	}
	for k, v := range b.Urls {
		consumer.OnEntry(ctx, k, v[0], v[1])
	}
	return nil
}

// Healthz is a health check point for Kubernetes
func (b *MockBackend) Healthz(ctx context.Context) error {
	return b.Healthy
}
