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

// Package cli contains the RootCmd for main and utilities for printing.
package cli

import (
	"github.com/spf13/cobra"
)

// RootCmd is the root cobra command for vanity.
var RootCmd = &cobra.Command{
	Use:     "vanity",
	Short:   "A go vanity server",
	Long:    `A go vanity server`,
	Version: getVersion(),
}
