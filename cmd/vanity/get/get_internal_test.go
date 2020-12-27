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

package get

import (
	"testing"

	"github.com/kami-zh/go-capturer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"l7e.io/vanity/apitest"
	"l7e.io/vanity/cmd/vanity/cli/backends"
	"l7e.io/vanity/cmd/vanity/cmdtest"
)

func TestList(t *testing.T) {
	backends.Set(&apitest.MockBackend{Urls: map[string][]string{"a.com/b": {"vcs", "vcsPath"}}})
	cmd := cmdtest.NewCommand(func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		assert.NoError(t, err)
	})

	out := capturer.CaptureOutput(func() {

		getCmd(cmd, []string{"a.com/b"})
	})

	assert.Equal(t, "a.com/b,vcs,vcsPath\n", out)
}
