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

package cli

import (
	"crypto/sha1" // nolint
	"encoding/hex"
	"fmt"

	json "github.com/gibson042/canonicaljson-go"
)

// SHA1FromJSON generates a SHA1 hash from a canonical JSON object.
func SHA1FromJSON(doc string) (string, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(doc), &v); err != nil {
		return "", fmt.Errorf("oi")
	}
	bytes, _ := json.Marshal(v)
	sum := sha1.Sum(bytes) // nolint
	return hex.EncodeToString(sum[0:]), nil
}

// SHA1FromString generates a SHA1 hash from a string.
func SHA1FromString(s string) string {
	sum := sha1.Sum([]byte(s)) // nolint
	return hex.EncodeToString(sum[0:])
}
