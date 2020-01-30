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

package vanity_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"l7e.io/vanity"
	"l7e.io/vanity/pkg/memory"
)

func TestInMemoryAPI(t *testing.T) {
	Convey("Ensure get obtains entry from AddEntry", t, func() {
		api := memory.NewInMemoryAPI()
		api.AddEntry("a", "va", "vaPath")

		vcs, vcsPath, err := api.Get(context.Background(), "a")
		So(err, ShouldBeNil)
		So(vcs, ShouldEqual, "va")
		So(vcsPath, ShouldEqual, "vaPath")
	})

	Convey("Ensure get returns error for unknown entry", t, func() {
		api := memory.NewInMemoryAPI()

		vcs, vcsPath, err := api.Get(context.Background(), "foo")
		So(vcs, ShouldBeEmpty)
		So(vcsPath, ShouldBeEmpty)
		So(err, ShouldEqual, vanity.ErrNotFound)
	})

	Convey("Ensure get obtains entry from Add", t, func() {
		api := memory.NewInMemoryAPI()

		err := api.Add(context.Background(), "a", "va", "vaPath")
		So(err, ShouldBeNil)

		vcs, vcsPath, err := api.Get(context.Background(), "a")
		So(err, ShouldBeNil)
		So(vcs, ShouldEqual, "va")
		So(vcsPath, ShouldEqual, "vaPath")
	})

	Convey("Ensure Remove works properly", t, func() {
		api := memory.NewInMemoryAPI()
		api.AddEntry("a", "va", "vaPath")

		err := api.Remove(context.Background(), "a")
		So(err, ShouldBeNil)

		vcs, vcsPath, err := api.Get(context.Background(), "a")
		So(vcs, ShouldBeEmpty)
		So(vcsPath, ShouldBeEmpty)
		So(err, ShouldEqual, vanity.ErrNotFound)
	})

	Convey("Ensure Remove non-existing entry returns an error", t, func() {
		api := memory.NewInMemoryAPI()

		err := api.Remove(context.Background(), "foo")
		So(err, ShouldEqual, vanity.ErrNotFound)
	})

	Convey("Ensure close always returns nil", t, func() {
		api := memory.NewInMemoryAPI()
		err := api.Close()

		So(err, ShouldBeNil)

		err = api.Close()

		So(err, ShouldBeNil)
	})

	Convey("Ensure always healthy", t, func() {
		api := memory.NewInMemoryAPI()
		err := api.Healthz(context.Background())

		So(err, ShouldBeNil)
	})

	Convey("Ensure List obtains all entries from AddEntry", t, func() {
		api := memory.NewInMemoryAPI()
		api.AddEntry("a", "va", "vaPath")
		api.AddEntry("b", "vb", "vbPath")
		api.AddEntry("c", "vc", "vcPath")

		count := 0
		err := api.List(context.Background(),
			vanity.ConsumerFunc(func(_ context.Context, importPath, vcs, vcsPath string) {
				count++
				So(importPath, ShouldNotBeEmpty)
				So(vcs, ShouldEqual, fmt.Sprintf("v%s", importPath))
				So(vcsPath, ShouldEqual, fmt.Sprintf("v%sPath", importPath))
			}))
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 3)
	})

	Convey("Ensure List can be canceled", t, func() {
		api := memory.NewInMemoryAPI()
		api.AddEntry("a", "va", "vaPath")
		api.AddEntry("b", "vb", "vbPath")
		api.AddEntry("c", "vc", "vcPath")

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

		err := api.List(ctx,
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
