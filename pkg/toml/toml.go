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

package toml

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml"

	"l7e.io/vanity"
)

type entry struct {
	ImportPath string `toml:"import_path"`
	Vcs        string
	VcsPath    string `toml:"vcs_path"`
}

type tomlBE struct {
	entries map[string]*entry
	closed  bool
}

type settings struct {
	Tables []string
	Path   string
	Reader io.Reader
	String string
	Bytes  []byte
}

// An Option is an option for a TOML-based Backend.
type Option interface {
	Apply(*settings)
}

var (
	errNoContentSpecified     = fmt.Errorf("no content specified")
	errTableDoesNotExist      = fmt.Errorf("table does not exist")
	errImportPathNotSpecified = fmt.Errorf("import_path not specified")
	errVcsNotSpecified        = fmt.Errorf("vcs not specified")
	errVcsPathNotSpecified    = fmt.Errorf("vcs_path not specified")
)

// InTable is used to specify the table the configuration can be found.
// Dotted table names have their tokens specified separately, in order.
func InTable(tables ...string) Option {
	return tablesOption{tables: tables}
}

type tablesOption struct{ tables []string }

func (t tablesOption) Apply(o *settings) {
	o.Tables = t.tables
}

// FromFile is used to specify the path of the TOML configuration file.
func FromFile(path string) Option {
	return fileOption{path: path}
}

type fileOption struct{ path string }

func (f fileOption) Apply(o *settings) {
	o.Path = f.path
}

// FromReader is used to specify a Reader that contains the TOML contents of the configuration.
func FromReader(r io.Reader) Option {
	return readerOption{reader: r}
}

type readerOption struct{ reader io.Reader }

func (r readerOption) Apply(o *settings) {
	o.Reader = r.reader
}

// FromString is used to specify the string TOML contents of the configuration.
func FromString(s string) Option {
	return stringOption{string: s}
}

type stringOption struct{ string string }

func (s stringOption) Apply(o *settings) {
	o.String = s.string
}

// FromBytes is used to specify the byte TOML contents of the configuration.
func FromBytes(b []byte) Option {
	return bytesOption{bytes: b}
}

type bytesOption struct{ bytes []byte }

func (b bytesOption) Apply(o *settings) {
	o.Bytes = b.bytes
}

// NewTOMLBackend creates a new TOML-backend using the specified options.
func NewTOMLBackend(options ...Option) (vanity.Backend, error) {
	s := settings{Tables: []string{}}
	for _, o := range options {
		o.Apply(&s)
	}

	b, err := readContent(&s)
	if err != nil {
		return nil, err
	}

	tree, err := toml.LoadBytes(b)
	if err != nil {
		return nil, err
	}

	paths := s.Tables
	tables := paths[:len(paths)-1]
	key := paths[len(paths)-1]

	tree, ok := tree.GetPath(tables).(*toml.Tree)
	if !ok {
		return nil, errTableDoesNotExist
	}

	array, ok := tree.Get(key).([]*toml.Tree)
	if !ok {
		return nil, errTableDoesNotExist
	}

	entries, err := readEntries(array)
	if err != nil {
		return nil, err
	}

	return &tomlBE{entries: entries}, nil
}

// readContent returns the TOML content that the settings select. The settings
// supply the content as bytes, as a string, as a file path, or as a reader.
func readContent(s *settings) ([]byte, error) {
	switch {
	case s.Bytes != nil:
		return s.Bytes, nil
	case s.String != "":
		return []byte(s.String), nil
	case s.Path != "":
		file, err := os.Open(s.Path)
		if err != nil {
			return nil, err
		}
		defer func() { _ = file.Close() }()

		return io.ReadAll(file)
	case s.Reader != nil:
		return io.ReadAll(s.Reader)
	default:
		return nil, errNoContentSpecified
	}
}

// readEntries converts the TOML trees to entries, keyed by import path. It
// returns an error when a tree omits a required field.
func readEntries(array []*toml.Tree) (map[string]*entry, error) {
	entries := make(map[string]*entry)

	for _, z := range array {
		e := &entry{}
		if err := z.Unmarshal(e); err != nil {
			return nil, err
		}

		switch {
		case e.ImportPath == "":
			return nil, errImportPathNotSpecified
		case e.Vcs == "":
			return nil, errVcsNotSpecified
		case e.VcsPath == "":
			return nil, errVcsPathNotSpecified
		}

		entries[e.ImportPath] = e
	}

	return entries, nil
}

func (s *tomlBE) Close() error {
	s.closed = true
	s.entries = nil

	return nil
}

func (s *tomlBE) Get(_ context.Context, importPath string) (vcs, vcsPath string, err error) {
	if err := s.check(); err != nil {
		return "", "", err
	}

	e, found := s.entries[importPath]
	if !found {
		return "", "", vanity.ErrNotFound
	}

	return e.Vcs, e.VcsPath, nil
}

func (s *tomlBE) Add(_ context.Context, importPath, vcs, vcsPath string) error {
	return vanity.ErrNotSupported
}

func (s *tomlBE) Remove(_ context.Context, importPath string) error {
	return vanity.ErrNotSupported
}

func (s *tomlBE) List(ctx context.Context, consumer vanity.Consumer) error {
	if err := s.check(); err != nil {
		return err
	}

	c := make(map[string]*entry)
	for k, v := range s.entries {
		c[k] = v
	}

	for importPath, e := range c {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		consumer.OnEntry(ctx, importPath, e.Vcs, e.VcsPath)
	}

	return nil
}

func (s *tomlBE) Healthz(_ context.Context) error {
	return nil
}

// check reports whether the caller closed this Backend.
func (s *tomlBE) check() error {
	if s.closed {
		return vanity.ErrAlreadyClosed
	}

	return nil
}
