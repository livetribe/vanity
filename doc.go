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
Package vanity contains an HTTP Handler that provides go vanity URL support.

See https://cloud.google.com/spanner/docs/getting-started/go/ for an
introduction to Cloud Spanner and additional help on using this API.

See https://godoc.org/cloud.google.com/go for authentication, timeouts,
connection pooling and similar aspects of this package.


Creating a Client

To start working with this package, create a client that refers to the database
of interest:

    ctx := context.Background()
    client, err := spanner.NewClient(ctx, "projects/P/instances/I/databases/D")
    if err != nil {
        // TODO: Handle error.
    }
    defer client.Close()

Remember to close the client after use to free up the sessions in the session
pool.
*/
package vanity  // import "l7e.io/vanity"
