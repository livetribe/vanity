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

package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"l7e.io/vanity/cmd/vanity/cli"
)

func TestSHA1FromString(t *testing.T) {
	sha1 := cli.SHA1FromString("how now brown cow")
	assert.Equal(t, "0de4bd7dfbc0f048319be1dad049a6cd7bede689", sha1)
}

func TestSHA1FromJSON(t *testing.T) {
	expected := "ae141a09b1c9d2a7cb26c7741fcb5e762c51b80f"

	sha1, err := cli.SHA1FromJSON("{\"how\": \"now\", \"brown\": \"cow\"}")
	assert.NoError(t, err)
	assert.Equal(t, expected, sha1)

	sha1, err = cli.SHA1FromJSON("{\"brown\": \"cow\", \"how\": \"now\"}")
	assert.NoError(t, err)
	assert.Equal(t, expected, sha1)
}

func TestSHA1FromJSON_badJSON(t *testing.T) {
	_, err := cli.SHA1FromJSON("{\"how\": \"now\"")
	assert.Error(t, err)
}
