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

/*
Package memory provides an in-memory implementation of Backend.

This in-memory backend is an implementation of vanity.Backend.  Vanity URL
configurations are stored in-memory.  The implementation provides a convenience
method for adding configurations without having to check for errors.


Creating an In-Memory Backend

To create an in-memory backend instance:

    be := memory.NewInMemoryAPI()
    defer be.Close()

Remember to close the backend instance after use.


Using the AddEntry Convenience Method

The implementation provides a convenience method for adding configurations
without having to pass a context.Context instance and check for errors.

Calling AddEntry on a closed instance has no effect; the call is virtually
ignored.

    be := memory.NewInMemoryAPI()
    defer be.Close()

    be.AddEntry("l7e.io/vanity", "git", "https://github.com/livetribe/vanity")

*/
package memory // import "l7e.io/vanity/pkg/memory"
