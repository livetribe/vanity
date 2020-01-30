/*
 * Copyright (c) 2019 the original author or authors.
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

// Package datastore contains the GCP Datastore Backend.
package datastore

import (
	"context"
	"sync"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
	"l7e.io/vanity"
)

const (
	kind = "GolangVanityEntry"
	key  = "ImportPath"
)

type datastoreClient struct {
	client *datastore.Client
	lock   sync.RWMutex
}

// NewClient creates a new Client for a given dataset.  If the project ID is
// empty, it is derived from the DATASTORE_PROJECT_ID environment variable.
// If the DATASTORE_EMULATOR_HOST environment variable is set, client will use
// its value to connect to a locally-running datastore emulator.
// DetectProjectID can be passed as the projectID argument to instruct
// NewClient to detect the project ID from the credentials.
func NewClient(projectID string, opts ...option.ClientOption) (vanity.Backend, error) {
	client, err := datastore.NewClient(context.Background(), projectID, opts...)
	if err != nil {
		return nil, err
	}

	return &datastoreClient{
		client: client,
	}, nil
}

func (d *datastoreClient) checkClosed() error {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if d.client == nil {
		return vanity.ErrAlreadyClosed
	}

	return nil
}

func (d *datastoreClient) Healthz(ctx context.Context) error {
	return d.List(ctx, nil)
}

func (d *datastoreClient) Close() error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.client == nil {
		return nil
	}

	err := d.client.Close()
	d.client = nil

	return err
}

func (d *datastoreClient) Get(ctx context.Context, importPath string) (vcs, vcsPath string, err error) {
	if err = d.checkClosed(); err != nil {
		return
	}

	key := datastore.NameKey(kind, importPath, nil)

	var e Entry

	if err := d.client.Get(ctx, key, &e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return "", "", vanity.ErrNotFound
		}

		return "", "", err
	}

	return e.Vcs, e.VcsRoot, nil
}

func (d *datastoreClient) Add(ctx context.Context, importPath, vcs, vcsPath string) error {
	if err := d.checkClosed(); err != nil {
		return err
	}

	var add = &Entry{ImportPath: importPath, Vcs: vcs, VcsRoot: vcsPath}

	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.NameKey(kind, importPath, nil)

		if _, err := tx.Put(key, add); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (d *datastoreClient) Remove(ctx context.Context, importPath string) error {
	if err := d.checkClosed(); err != nil {
		return err
	}

	_, err := d.client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.NameKey(kind, importPath, nil)

		if e := tx.Delete(key); e != nil {
			return e
		}

		return nil
	})

	return err
}

func (d *datastoreClient) List(ctx context.Context, consumer vanity.Consumer) error {
	if err := d.checkClosed(); err != nil {
		return err
	}

	query := datastore.NewQuery(kind).Order(key)

	var all []*Entry

	if _, err := d.client.GetAll(ctx, query, &all); err != nil {
		return err
	}

	if consumer == nil {
		return nil
	}

	for _, e := range all {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		consumer.OnEntry(ctx, e.ImportPath, e.Vcs, e.VcsRoot)
	}

	return nil
}
