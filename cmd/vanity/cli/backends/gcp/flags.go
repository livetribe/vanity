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
	"encoding/base64"
	"fmt"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	"l7e.io/vanity/cmd/vanity/cli"
	"l7e.io/vanity/cmd/vanity/cli/log"
	"l7e.io/vanity/cmd/vanity/server/interceptors"
	"l7e.io/vanity/pkg/gcp"
)

const (
	glbHCKey        = "glb-healthcheck"
	apiKey          = "api-key"
	credentialsFile = "credentials-file"
	credentialsJSON = "credentials-json"
)

// InitFlags initializes the Cobra command with flags common to Google API clients.
func InitFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.BoolP(glbHCKey, "", true, "Add Google load balancer health-check interceptor")
	flags.StringP(apiKey, "", "", "API key to be used as the basis for authentication (optional)")
	flags.StringP(credentialsFile, "", "", "service account or refresh token JSON credentials file (optional)")
	flags.StringP(credentialsJSON, "", "", "service account or refresh token JSON credentials in base64 (optional)")

	viper.RegisterAlias(apiKey, "google-api.key")
	viper.RegisterAlias(credentialsFile, "google-api."+credentialsFile)
	viper.RegisterAlias(credentialsJSON, "google-api."+credentialsJSON)
}

// Helper provides convenient access to the Google API Cobra flags.
type Helper struct {
	*cli.FlagSet
}

// NewHelper wraps the Cobra command's flags with a utility wrapper to assist in
// the collection of Google API client options.
func NewHelper(cmd *cobra.Command) *Helper {
	return &Helper{cli.Flags(cmd)}
}

// AddInterceptors will register interceptors which handle Google load balancer
// user agent checks for a GCE ingress agent.
func (h *Helper) AddInterceptors() error {
	glb, err := h.GetBool(glbHCKey)
	if err != nil {
		return err
	}

	if glb {
		interceptors.RegisterInterceptor(gcp.GLB)
		glog.V(log.Debug).Info("Added GLB interceptor")
	}
	return nil
}

// GetClientOptions obtains the Google API client options.
func (h *Helper) GetClientOptions() ([]option.ClientOption, error) {
	var co []option.ClientOption

	co = h.CollectAPIKeyOption(co)
	co, err := h.CollectCredentialsFileOption(co)
	if err != nil {
		return nil, err
	}
	co, err = h.CollectCredentialsOption(co)
	if err != nil {
		return nil, err
	}

	return co, nil
}

// CollectAPIKeyOption adds the Google API key, if passed as a command flag, to
// the Google API client options.
func (h *Helper) CollectAPIKeyOption(options []option.ClientOption) []option.ClientOption {
	key, found := h.GetValue(apiKey)
	if !found {
		return options
	}

	glog.V(log.Debug).Infof("Google API key (sha1): %s", cli.SHA1FromString(key))

	return append(options, option.WithAPIKey(key))
}

// CollectCredentialsFileOption adds the Google credentials file, if passed as
// a command flag, to the Google API client options.
func (h *Helper) CollectCredentialsFileOption(options []option.ClientOption) ([]option.ClientOption, error) {
	cf, found := h.GetValue(credentialsFile)
	if !found {
		return options, nil
	}

	if !cli.FileExists(cf) {
		return options, fmt.Errorf("file %s does not exist", cf)
	}

	glog.V(log.Debug).Infof("Google API credentials file: %s", cf)

	return append(options, option.WithCredentialsFile(cf)), nil
}

// CollectCredentialsOption adds the Google credentials, if passed as a
// command flag, to the Google API client options.
func (h *Helper) CollectCredentialsOption(options []option.ClientOption) ([]option.ClientOption, error) {
	b64, found := h.GetValue(credentialsJSON)
	if !found {
		return options, nil
	}

	cj, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid base64 encoding of %s", credentialsJSON)
	}

	// if cj is not valid JSON, this function will fail
	sha1, err := cli.SHA1FromJSON(string(cj))
	if err != nil {
		return nil, err
	}

	glog.V(log.Debug).Infof("Google API credentials (sha1): %s", sha1)

	return append(options, option.WithCredentialsJSON(cj)), nil
}
