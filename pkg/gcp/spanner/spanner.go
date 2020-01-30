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

// Package spanner contains the GCP Spanner Backend.
package spanner // import "l7e.io/vanity/pkg/gcp/spanner"

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"l7e.io/vanity"
)

const (
	importPathColumn = "import_path"
	vcsColumn        = "vcs"
	vcsPathColumn    = "vcs_path"
)

type spannerClient struct {
	table  string
	client *spanner.Client
	lock   sync.RWMutex
}

// NewClient creates a client to a database. A valid database name has the
// form projects/PROJECT_ID/instances/INSTANCE_ID/databases/DATABASE_ID.
func NewClient(ctx context.Context, database string, opts ...BackendOption) (api vanity.Backend, err error) {
	var dataClient *spanner.Client

	s := collectSettings(opts...)
	if s.config != nil {
		dataClient, err = spanner.NewClientWithConfig(ctx, database, *s.config, s.options...)
	} else {
		dataClient, err = spanner.NewClient(ctx, database, s.options...)
	}

	if err != nil {
		return nil, err
	}

	return &spannerClient{
		table:  s.table,
		client: dataClient,
	}, nil
}

func (s *spannerClient) checkClosed() error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.client == nil {
		return vanity.ErrAlreadyClosed
	}

	return nil
}

func (s *spannerClient) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.client == nil {
		return nil
	}

	c := s.client
	go c.Close()

	s.client = nil

	return nil
}

func (s *spannerClient) Healthz(ctx context.Context) error {
	return s.List(ctx, vanity.ConsumerFunc(func(ctx context.Context, importPath, vcs, vcsPath string) {}))
}

func (s *spannerClient) Get(ctx context.Context, importPath string) (vcs, vcsPath string, err error) {
	if err = s.checkClosed(); err != nil {
		return
	}

	row, err :=
		s.client.Single().ReadRow(ctx, s.table, spanner.Key{importPath}, []string{vcsColumn, vcsPathColumn})
	if err != nil {
		if spanner.ErrCode(err) == codes.NotFound {
			return "", "", vanity.ErrNotFound
		}
		return "", "", fmt.Errorf(unableToRetrieve, importPath, err)
	}

	if err = row.ColumnByName(vcsColumn, &vcs); err != nil {
		return "", "", fmt.Errorf(unableToExtractVcs, importPath, err)
	}

	if err = row.ColumnByName(vcsPathColumn, &vcsPath); err != nil {
		return "", "", fmt.Errorf(unableToExtractVcsPath, importPath, err)
	}

	return vcs, vcsPath, err
}

func (s *spannerClient) Add(ctx context.Context, importPath, vcs, vcsPath string) error {
	if err := s.checkClosed(); err != nil {
		return err
	}

	ms := []*spanner.Mutation{
		spanner.Insert(
			s.table,
			[]string{importPathColumn, vcsColumn, vcsPathColumn},
			[]interface{}{importPath, vcs, vcsPath}),
	}
	_, err := s.client.Apply(ctx, ms)

	return err
}

func (s *spannerClient) Remove(ctx context.Context, importPath string) error {
	if err := s.checkClosed(); err != nil {
		return err
	}

	ms := []*spanner.Mutation{
		spanner.Delete(s.table, spanner.Key{importPath}),
	}
	_, err := s.client.Apply(ctx, ms)

	return err
}

func (s *spannerClient) List(ctx context.Context, consumer vanity.Consumer) error {
	if err := s.checkClosed(); err != nil {
		return err
	}

	stmt := spanner.Statement{SQL: s.sql()}
	iter := s.client.Single().Query(ctx, stmt)

	defer iter.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		row, err := iter.Next()

		switch {
		case err == iterator.Done:
			return nil
		case err != nil:
			return err
		}

		var importPath, vcs, vcsPath string

		if err := row.Columns(&importPath, &vcs, &vcsPath); err != nil {
			return err
		}

		consumer.OnEntry(ctx, importPath, vcs, vcsPath)
	}
}

func (s *spannerClient) sql() string {
	return fmt.Sprintf("SELECT %s, %s, %s FROM %s", importPathColumn, vcsColumn, vcsPathColumn, s.table) // nolint:gosec
}
