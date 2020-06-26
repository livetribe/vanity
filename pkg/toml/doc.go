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

/*
Package toml provides an R/O, in-memory, TOML-initialized, implementation of Backend.

This in-memory backend is an R/O implementation of vanity.Backend which is initialized
using TOML, https://github.com/toml-lang/toml.  The vanity URL configurations are
stored in-memory.


TOML Configuration

Vanity entries are specified via an array of tables, https://github.com/toml-lang/toml#array-of-tables.
Each table has the following keys:

* importPath - the vanity URL, e.g. "l7e.io/vanity"

* vcs - the version control system, e.g. "git"

* vcsPath - the actual root of the version control system pointed to by the import path

The vanity entry for this project could be

	[[entry]]
	importPath=l7e.io/vanity
	vcs=git
	vcsPath=https://github.com/livetribe/vanity

Here, the table name is "entry" and is specified using the InTable option.

The TOML configuration can be passed in during construction by passing in the appropriate
option to the NewTOMLBackend method to specify file, Reader, string, or byte array content.
Specifying multiple content, e.g. both a file and string results in undefined behavior.

Creating a TOML-based Backend

To create a backend instance:

    be := toml.NewTOMLBackend(InTable("obj"), FromString(`
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
    defer be.Close()

Remember to close the backend instance after use.

*/
package toml // import "l7e.io/vanity/pkg/toml"
