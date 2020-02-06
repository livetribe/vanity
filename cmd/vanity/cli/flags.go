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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// FlagSet extends pflag.FlagSet with methods that obtain values from
// the environment.
type FlagSet struct {
	*pflag.FlagSet
}

// Flags returns the a wrapped FlagSet that applies to this command (local
// and persistent declared here and by all parents).
func Flags(cmd *cobra.Command) *FlagSet {
	return &FlagSet{
		cmd.Flags(),
	}
}

// GetValue is a utility that attempts to obtain a value from the command line
// and if there's no value, then it attempts to obtain the value from the
// environmental variable.
func (f *FlagSet) GetValue(flagName string) (value string, found bool) {
	value, err := f.GetString(flagName)
	if err == nil && len(value) > 0 {
		return value, true
	}
	if len(value) > 0 {
		return value, true
	}

	v := viper.Get(flagName)
	if v == nil {
		return "", false
	}

	return v.(string), true
}
