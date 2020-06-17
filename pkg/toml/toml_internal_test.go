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
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"l7e.io/vanity"
)

func TestTomlConstructionOptions(t *testing.T) {
	Convey("Test construction options", t, func() {
		Convey("Ensure InTable correctly fills settings", func() {
			o := InTable("one", "two", "three")
			So(o, ShouldNotBeNil)

			var s settings
			o.Apply(&s)

			So(s.Tables, ShouldResemble, []string{"one", "two", "three"})
		})

		Convey("Ensure FromFile correctly fills settings", func() {
			o := FromFile(".config/vanity/entries.toml")
			So(o, ShouldNotBeNil)

			var s settings
			o.Apply(&s)

			So(s.Path, ShouldEqual, ".config/vanity/entries.toml")
		})

		Convey("Ensure FromReader correctly fills settings", func() {
			r := strings.NewReader(`
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
`)

			o := FromReader(r)
			So(o, ShouldNotBeNil)

			var s settings
			o.Apply(&s)

			So(s.Reader, ShouldEqual, r)
		})

		Convey("Ensure FromString correctly fills settings", func() {
			str := `
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
`

			o := FromString(str)
			So(o, ShouldNotBeNil)

			var s settings
			o.Apply(&s)

			So(s.String, ShouldEqual, str)
		})

		Convey("Ensure FromBytes correctly fills settings", func() {
			b := []byte(`
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
`)

			o := FromBytes(b)
			So(o, ShouldNotBeNil)

			var s settings
			o.Apply(&s)

			So(s.Bytes, ShouldResemble, b)
		})
	})
}

func TestTomlConstruction(t *testing.T) {
	Convey("Test TOML backend construction", t, func() {
		Convey("Construction with a file", func() {
			content := []byte(`
[[obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
[[obj]]
import_path = "l7e.io/two"
vcs = "git"
vcs_path = "https://github.com/livetribe/two"
[[obj]]
import_path = "l7e.io/three"
vcs = "git"
vcs_path = "https://github.com/livetribe/three"
`)
			tmpfile, err := ioutil.TempFile("", "toml")
			So(err, ShouldBeNil)

			defer func() {
				_ = tmpfile.Close()
				_ = os.Remove(tmpfile.Name())
			}()

			_, err = tmpfile.Write(content)
			So(err, ShouldBeNil)

			be, err := NewTOMLBackend(InTable("obj"), FromFile(tmpfile.Name()))
			So(err, ShouldBeNil)
			So(be, ShouldNotBeNil)

			verify(be)
		})

		Convey("Construction with Reader", func() {
			r := strings.NewReader(`
[[obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
[[obj]]
import_path = "l7e.io/two"
vcs = "git"
vcs_path = "https://github.com/livetribe/two"
[[obj]]
import_path = "l7e.io/three"
vcs = "git"
vcs_path = "https://github.com/livetribe/three"
`)

			be, err := NewTOMLBackend(InTable("obj"), FromReader(r))
			So(err, ShouldBeNil)
			So(be, ShouldNotBeNil)

			verify(be)
		})

		Convey("Construction with bytes", func() {
			be, err := NewTOMLBackend(InTable("obj"), FromBytes([]byte(`
[[obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
[[obj]]
import_path = "l7e.io/two"
vcs = "git"
vcs_path = "https://github.com/livetribe/two"
[[obj]]
import_path = "l7e.io/three"
vcs = "git"
vcs_path = "https://github.com/livetribe/three"
`)))
			So(err, ShouldBeNil)
			So(be, ShouldNotBeNil)

			verify(be)
		})

		Convey("Construction with no nested path", func() {
			be, err := NewTOMLBackend(InTable("obj"), FromString(`
[[obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
[[obj]]
import_path = "l7e.io/two"
vcs = "git"
vcs_path = "https://github.com/livetribe/two"
[[obj]]
import_path = "l7e.io/three"
vcs = "git"
vcs_path = "https://github.com/livetribe/three"
`))
			So(err, ShouldBeNil)
			So(be, ShouldNotBeNil)

			verify(be)
		})

		Convey("Construction with nested path", func() {
			be, err := NewTOMLBackend(InTable("a", "b", "obj"), FromString(`
[[a.b.obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
[[a.b.obj]]
import_path = "l7e.io/two"
vcs = "git"
vcs_path = "https://github.com/livetribe/two"
[[a.b.obj]]
import_path = "l7e.io/three"
vcs = "git"
vcs_path = "https://github.com/livetribe/three"
`))
			So(err, ShouldBeNil)
			So(be, ShouldNotBeNil)

			verify(be)
		})

		Convey("Construction with bad TOML", func() {
			be, err := NewTOMLBackend(InTable("a", "b", "obj"), FromString(`
How now brown cow
`))
			So(err, ShouldNotBeNil)
			So(be, ShouldBeNil)
		})

		Convey("Construction with bad entry", func() {
			Convey("missing import_path", func() {
				be, err := NewTOMLBackend(InTable("obj"), FromString(`
[[obj]]
# import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
`))
				So(err, ShouldNotBeNil)
				So(be, ShouldBeNil)
			})

			Convey("missing vcs", func() {
				be, err := NewTOMLBackend(InTable("obj"), FromString(`
[[obj]]
import_path = "l7e.io/one"
# vcs = "git"
vcs_path = "https://github.com/livetribe/one"
`))
				So(err, ShouldNotBeNil)
				So(be, ShouldBeNil)
			})

			Convey("missing vcs_path", func() {
				be, err := NewTOMLBackend(InTable("obj"), FromString(`
[[obj]]
import_path = "l7e.io/one"
vcs = "git"
# vcs_path = "https://github.com/livetribe/one"
`))
				So(err, ShouldNotBeNil)
				So(be, ShouldBeNil)
			})
		})

		Convey("Construction with bad nested path", func() {
			Convey("bad path", func() {
				be, err := NewTOMLBackend(InTable("x", "y", "obj"), FromString(`
[[a.b.obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
`))
				So(err, ShouldNotBeNil)
				So(be, ShouldBeNil)
			})

			Convey("missing key", func() {
				be, err := NewTOMLBackend(InTable("a", "b", "wrong"), FromString(`
[[a.b.obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
`))
				So(err, ShouldNotBeNil)
				So(be, ShouldBeNil)
			})
		})

		Convey("Construction with no content", func() {
			be, err := NewTOMLBackend(InTable("x", "y", "obj"))
			So(err, ShouldNotBeNil)
			So(be, ShouldBeNil)
		})
	})
}

func TestTomlBackend(t *testing.T) {
	var expected = map[string]entry{
		"l7e.io/one":   {"l7e.io/one", "git", "https://github.com/livetribe/one"},
		"l7e.io/two":   {"l7e.io/two", "git", "https://github.com/livetribe/two"},
		"l7e.io/three": {"l7e.io/three", "git", "https://github.com/livetribe/three"},
	}

	Convey("Test TOML backend methods", t, func() {
		be, err := NewTOMLBackend(InTable("obj"), FromString(`
[[obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
[[obj]]
import_path = "l7e.io/two"
vcs = "git"
vcs_path = "https://github.com/livetribe/two"
[[obj]]
import_path = "l7e.io/three"
vcs = "git"
vcs_path = "https://github.com/livetribe/three"
`))
		So(err, ShouldBeNil)
		So(be, ShouldNotBeNil)

		Convey("Add should panic", func() {
			err = be.Add(context.Background(), "l7e.io/one", "git", "https://github.com/livetribe/one")
			So(err, ShouldNotBeNil)
		})

		Convey("Remove should panic", func() {
			err = be.Remove(context.Background(), "l7e.io/one")
			So(err, ShouldNotBeNil)
		})

		Convey("Test List", func() {
			entries := make(map[string]entry)
			err := be.List(context.Background(),
				vanity.ConsumerFunc(func(_ context.Context, importPath, vcs, vcsPath string) {
					entries[importPath] = entry{importPath, vcs, vcsPath}
				}))
			So(err, ShouldBeNil)
			So(entries, ShouldResemble, expected)
		})

		Convey("Ensure List can be canceled", func() {
			start := &sync.WaitGroup{}
			start.Add(1)
			end := &sync.WaitGroup{}
			end.Add(1)

			ctx, cancel := context.WithCancel(context.Background())

			go func() {
				start.Wait()
				cancel()
				end.Done()
			}()

			err := be.List(ctx,
				vanity.ConsumerFunc(func(_ context.Context, importPath, vcs, vcsPath string) {
					start.Done()

					So(importPath, ShouldNotBeEmpty)
					So(vcs, ShouldEqual, "git")

					end.Wait()
				}))
			So(err, ShouldEqual, context.Canceled)
		})

		Convey("Ensure always healthy", func() {
			err := be.Healthz(context.Background())

			So(err, ShouldBeNil)
		})
	})
}

func TestTomlBE_Close(t *testing.T) {
	Convey("Test TOML Close methods", t, func() {
		be, err := NewTOMLBackend(InTable("obj"), FromString(`
[[obj]]
import_path = "l7e.io/one"
vcs = "git"
vcs_path = "https://github.com/livetribe/one"
[[obj]]
import_path = "l7e.io/two"
vcs = "git"
vcs_path = "https://github.com/livetribe/two"
[[obj]]
import_path = "l7e.io/three"
vcs = "git"
vcs_path = "https://github.com/livetribe/three"
`))
		So(err, ShouldBeNil)
		So(be, ShouldNotBeNil)

		Convey("Can close multiple times without error", func() {
			err = be.Close()
			So(err, ShouldBeNil)
			err = be.Close()
			So(err, ShouldBeNil)
		})

		Convey("Get should return an error", func() {
			err = be.Close()
			So(err, ShouldBeNil)

			_, _, err = be.Get(context.Background(), "l7e.io/three")
			So(err, ShouldNotBeNil)
		})

		Convey("List should return an error", func() {
			err = be.Close()
			So(err, ShouldBeNil)

			err = be.List(context.Background(), vanity.ConsumerFunc(func(_ context.Context, _, _, _ string) {}))
			So(err, ShouldNotBeNil)
		})
	})
}

func verify(be vanity.Backend) {
	vcs, vcsPath, err := be.Get(context.Background(), "l7e.io/one")
	So(err, ShouldBeNil)
	So(vcs, ShouldResemble, "git")
	So(vcsPath, ShouldResemble, "https://github.com/livetribe/one")

	vcs, vcsPath, err = be.Get(context.Background(), "l7e.io/two")
	So(err, ShouldBeNil)
	So(vcs, ShouldResemble, "git")
	So(vcsPath, ShouldResemble, "https://github.com/livetribe/two")

	vcs, vcsPath, err = be.Get(context.Background(), "l7e.io/three")
	So(err, ShouldBeNil)
	So(vcs, ShouldResemble, "git")
	So(vcsPath, ShouldResemble, "https://github.com/livetribe/three")
}
