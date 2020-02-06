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

package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"l7e.io/vanity"
)

func TestInMemoryAPI(t *testing.T) {

	Convey("Ensure get obtains entry from AddEntry", t, func() {
		be := NewInMemoryAPI()
		be.AddEntry("l7e.io/vanity", "git", "https://github.com/livetribe/vanity")

		vcs, vcsPath, err := be.Get(context.Background(), "l7e.io/vanity")
		So(err, ShouldBeNil)
		So(vcs, ShouldEqual, "git")
		So(vcsPath, ShouldEqual, "https://github.com/livetribe/vanity")
	})

	Convey("Ensure get returns error for unknown entry", t, func() {
		be := NewInMemoryAPI()

		vcs, vcsPath, err := be.Get(context.Background(), "foo")
		So(vcs, ShouldBeEmpty)
		So(vcsPath, ShouldBeEmpty)
		So(err, ShouldBeError, vanity.ErrNotFound)
	})

	Convey("Ensure get obtains entry from Add", t, func() {
		be := NewInMemoryAPI()

		err := be.Add(context.Background(), "l7e.io/vanity", "git", "https://github.com/livetribe/vanity")
		So(err, ShouldBeNil)

		vcs, vcsPath, err := be.Get(context.Background(), "l7e.io/vanity")
		So(err, ShouldBeNil)
		So(vcs, ShouldEqual, "git")
		So(vcsPath, ShouldEqual, "https://github.com/livetribe/vanity")
	})

	Convey("Ensure Remove works properly", t, func() {
		be := NewInMemoryAPI()
		be.AddEntry("l7e.io/vanity", "git", "https://github.com/livetribe/vanity")

		err := be.Remove(context.Background(), "l7e.io/vanity")
		So(err, ShouldBeNil)

		vcs, vcsPath, err := be.Get(context.Background(), "l7e.io/vanity")
		So(vcs, ShouldBeEmpty)
		So(vcsPath, ShouldBeEmpty)
		So(err, ShouldBeError, vanity.ErrNotFound)

		err = be.Remove(context.Background(), "l7e.io/vanity")
		So(err, ShouldBeError, vanity.ErrNotFound)
	})

	Convey("Ensure Remove non-existing entry returns an error", t, func() {
		be := NewInMemoryAPI()

		err := be.Remove(context.Background(), "foo")
		So(err, ShouldBeError, vanity.ErrNotFound)
	})

	Convey("Ensure close always returns nil", t, func() {
		be := NewInMemoryAPI()
		err := be.Close()

		So(err, ShouldBeNil)

		err = be.Close()

		So(err, ShouldBeNil)
	})

	Convey("Ensure always healthy", t, func() {
		be := NewInMemoryAPI()
		err := be.Healthz(context.Background())

		So(err, ShouldBeNil)
	})

	Convey("Ensure List obtains all entries from AddEntry", t, func() {
		be := NewInMemoryAPI()
		be.AddEntry("l7e.io/vanity", "git", "https://github.com/livetribe/vanity")
		be.AddEntry("m4o.io/pbf", "git", "https://github.com/magurl/pbf")

		entries := make(map[string]*entry)
		err := be.List(context.Background(),
			vanity.ConsumerFunc(func(_ context.Context, importPath, vcs, vcsPath string) {
				entries[importPath] = &entry{vcs: vcs, vcsPath: vcsPath}
			}))
		So(err, ShouldBeNil)
		So(*entries["l7e.io/vanity"], ShouldResemble, entry{"git", "https://github.com/livetribe/vanity"})
		So(*entries["m4o.io/pbf"], ShouldResemble, entry{"git", "https://github.com/magurl/pbf"})
	})

	Convey("Ensure List can be canceled", t, func() {
		be := NewInMemoryAPI()
		be.AddEntry("a", "va", "vaPath")
		be.AddEntry("b", "vb", "vbPath")
		be.AddEntry("c", "vc", "vcPath")

		start := &sync.WaitGroup{}
		start.Add(1) //nolint
		end := &sync.WaitGroup{}
		end.Add(1) //nolint

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
				So(vcs, ShouldEqual, fmt.Sprintf("v%s", importPath))
				So(vcsPath, ShouldEqual, fmt.Sprintf("v%sPath", importPath))

				end.Wait()
			}))
		So(err, ShouldEqual, context.Canceled)
	})
}

func TestInMemory_AddEntry(t *testing.T) {
	be := NewInMemoryAPI()

	be.AddEntry("l7e.io/vanity", "git", "https://github.com/livetribe/vanity")
	vcs, vcsPath, err := be.Get(context.Background(), "l7e.io/vanity")
	assert.NoError(t, err)
	assert.Equal(t, "git", vcs)
	assert.Equal(t, "https://github.com/livetribe/vanity", vcsPath)
}

func TestInMemory_Close(t *testing.T) {
	be := NewInMemoryAPI()
	assert.NoError(t, be.Close())
	assert.NoError(t, be.Close())

	err := be.Add(context.Background(), "l7e.io/vanity", "git", "https://github.com/livetribe/vanity")
	assert.Equal(t, err, vanity.ErrAlreadyClosed)

	_, _, err = be.Get(context.Background(), "l7e.io/vanity")
	assert.Equal(t, err, vanity.ErrAlreadyClosed)

	err = be.Remove(context.Background(), "l7e.io/vanity")
	assert.Equal(t, err, vanity.ErrAlreadyClosed)

	err = be.List(context.Background(), vanity.ConsumerFunc(func(_ context.Context, _, _, _ string) {}))
	assert.Equal(t, err, vanity.ErrAlreadyClosed)
}

func TestInMemory_List(t *testing.T) {
	be := NewInMemoryAPI()
	be.AddEntry("l7e.io/vanity", "git", "https://github.com/livetribe/vanity")
	be.AddEntry("m4o.io/pbf", "git", "https://github.com/magurl/pbf")

	entries := make(map[string]*entry)
	err := be.List(context.Background(),
		vanity.ConsumerFunc(func(_ context.Context, importPath, vcs, vcsPath string) {
			entries[importPath] = &entry{vcs: vcs, vcsPath: vcsPath}
		}))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entries))
	assert.Equal(t, &entry{"git", "https://github.com/livetribe/vanity"}, entries["l7e.io/vanity"])
	assert.Equal(t, &entry{"git", "https://github.com/magurl/pbf"}, entries["m4o.io/pbf"])
}

func TestInMemory_List_Timeout(t *testing.T) {
	be := NewInMemoryAPI()
	be.AddEntry("l7e.io/vanity", "git", "https://github.com/livetribe/vanity")
	be.AddEntry("m4o.io/pbf", "git", "https://github.com/magurl/pbf")

	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10) // nolint

	err := be.List(ctx,
		vanity.ConsumerFunc(func(_ context.Context, _, _, _ string) {
			time.Sleep(time.Millisecond * 50) // nolint
		}))
	assert.Equal(t, err, context.DeadlineExceeded)
}

func TestInMemory_Healthz(t *testing.T) {
	be := NewInMemoryAPI()
	assert.NoError(t, be.Healthz(context.Background())) // always healthy
}
