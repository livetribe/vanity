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

package server

import (
	"context"

	"l7e.io/vanity"
)

type be struct {
	healthy error
}

func (b *be) Close() error {
	return nil
}

// Get vanity URL configuration for a given import path.
func (b *be) Get(ctx context.Context, importPath string) (vcs, vcsPath string, err error) {
	return "", "", nil
}

// Add a vanity URL configuration.
func (b *be) Add(ctx context.Context, importPath, vcs, vcsPath string) error {
	return nil
}

// Remove a vanity URL configuration by it's key, the import path.
func (b *be) Remove(ctx context.Context, importPath string) error {
	return nil
}

// List all registered URL configurations, delivering them to the consumer callback.
func (b *be) List(ctx context.Context, consumer vanity.Consumer) error {
	return nil
}

// Healthz is a health check point for Kubernetes.
func (b *be) Healthz(ctx context.Context) error {
	return b.healthy
}
