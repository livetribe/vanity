/*
 * Copyright the original author or authors.
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

package server

import (
	"context"
	"fmt"
	"log/slog"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"l7e.io/yama"
)

func serverCmd(cmd *cobra.Command, _ []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("unable to bind viper to the command line flags: %w", err)
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	cfg := NewConfig()

	_, lifecycle, err := NewLifecycle(cfg, Fail(cancel))
	if err != nil {
		return fmt.Errorf("unable to build the server lifecycle: %w", err)
	}

	err = yama.RunUntilSignal(ctx, lifecycle, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("Vanity exited")

	if err != nil {
		return fmt.Errorf("unable to run the server: %w", err)
	}

	return nil
}
