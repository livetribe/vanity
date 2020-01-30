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

// Package gcp declares common Google API clients options.
package gcp

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

const (
	apiKey          = "api-key"
	credentialsFile = "credentials-file"
	credentialsJSON = "credentials-json"
)

// InitFlags initializes the Cobra command with flags common to Google API clients.
func InitFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.StringP(apiKey, "", "", "API key to be used as the basis for authentication (optional)")
	flags.StringP(credentialsFile, "", "", "service account or refresh token JSON credentials file (optional)")
	flags.StringP(credentialsJSON, "", "", "service account or refresh token JSON credentials (optional)")
}

// Helper provides convenient access to the Google API Cobra flags.
type Helper struct {
	*pflag.FlagSet
}

// NewHelper wraps the Cobra command's flags with a utility wrapper to assist in
// the collection of Google API client options.
func NewHelper(cmd *cobra.Command) *Helper {
	return &Helper{cmd.Flags()}
}

// GetClientOptions obtains the Google API client options.
func (h *Helper) GetClientOptions() []option.ClientOption {
	var co []option.ClientOption

	co = h.CollectAPIKeyOption(co)
	co = h.CollectCredentialsFileOption(co)
	co = h.CollectCredentialsOption(co)

	return co
}

// CollectAPIKeyOption adds the Google API key, if passed as a command flag, to
// the Google API client options.
func (h *Helper) CollectAPIKeyOption(options []option.ClientOption) []option.ClientOption {
	key := viper.GetString(apiKey)

	if key == "" {
		return options
	}

	return append(options, option.WithAPIKey(key))
}

// CollectCredentialsFileOption adds the Google credentials file, if passed as
// a command flag, to the Google API client options.
func (h *Helper) CollectCredentialsFileOption(options []option.ClientOption) []option.ClientOption {
	cf := viper.GetString(credentialsFile)

	if cf == "" {
		return options
	}

	return append(options, option.WithCredentialsFile(cf))
}

// CollectCredentialsOption adds the Google credentials, if passed as a
// command flag, to the Google API client options.
func (h *Helper) CollectCredentialsOption(options []option.ClientOption) []option.ClientOption {
	cj := viper.GetString(credentialsJSON)

	if cj == "" {
		return options
	}

	return append(options, option.WithCredentialsJSON([]byte(cj)))
}
