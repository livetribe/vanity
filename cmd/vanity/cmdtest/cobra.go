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

// Package cmdtest contains some helper functions for testing Cobra commands.
package cmdtest

import (
	"bytes"

	"github.com/spf13/cobra"
)

// ExecuteCommand executes the Cobra command cmd with arguments args.
func ExecuteCommand(cmd *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(cmd, args...)
	return output, err
}

func executeCommandC(cmd *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	//noinspection ALL
	cmd.SetOutput(buf)
	if args == nil {
		// if nil is passed then it defaults to
		// th actual command line args are passed
		args = []string{}
	}
	cmd.SetArgs(args)
	c, err = cmd.ExecuteC()
	return c, buf.String(), err
}

// NewCommand creates a new Cobra Command with run function f.
func NewCommand(f func(cmd *cobra.Command, args []string)) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "root",
		Args: cobra.NoArgs,
		Run:  f,
	}
	return cmd
}
