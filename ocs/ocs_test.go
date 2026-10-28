// Copyright (c) 2026 DistintiveLab, All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ocs

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusDetail(t *testing.T) {
	assert.Equal(t, "401 Unauthorized", StatusDetail(http.StatusUnauthorized, nil))
	assert.Equal(t, "401 Unauthorized", StatusDetail(http.StatusUnauthorized, []byte("  \n")))
	assert.Equal(t, `500 Internal Server Error: {"message":"boom"}`,
		StatusDetail(http.StatusInternalServerError, []byte(`{"message":"boom"}`)))

	long := strings.Repeat("a", maxStatusDetailBody+50)
	detail := StatusDetail(http.StatusBadGateway, []byte(long))
	assert.True(t, strings.HasSuffix(detail, "..."))
	assert.Len(t, detail, len(http.StatusText(http.StatusBadGateway))+len("502 ")+maxStatusDetailBody+len(": ..."))
}

func TestOCSStatusDetail(t *testing.T) {
	assert.Equal(t, "OCS status 997", OCSStatusDetail(997, ""))
	assert.Equal(t, "OCS status 997: Unauthorised", OCSStatusDetail(997, "Unauthorised"))
}

func capabilitiesFrom(body string) (Capabilities, error) {
	var wrapper struct {
		Ocs Capabilities `json:"ocs"`
	}
	err := json.Unmarshal([]byte(body), &wrapper)
	return wrapper.Ocs, err
}

func TestCapabilitiesUnmarshalObjectData(t *testing.T) {
	c, err := capabilitiesFrom(`{"ocs":{"meta":{"status":"ok","statuscode":100,"message":"OK"},"data":{"capabilities":{"spreed":{"features":["chat-v2","conversation-v4"]}}}}}`)
	assert.NoError(t, err)
	assert.Equal(t, StatusOK, c.OCSMeta.StatusCode)
	assert.Equal(t, []string{"chat-v2", "conversation-v4"}, c.Data.Capabilities.SpreedCapabilities.Features)
}

func TestCapabilitiesUnmarshalArrayData(t *testing.T) {
	c, err := capabilitiesFrom(`{"ocs":{"meta":{"status":"ok","statuscode":100,"message":"OK"},"data":[{"capabilities":{"spreed":{"features":["chat-v2","delete-messages"]}}}]}}`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"chat-v2", "delete-messages"}, c.Data.Capabilities.SpreedCapabilities.Features)
}

func TestCapabilitiesUnmarshalUnauthorized(t *testing.T) {
	c, err := capabilitiesFrom(`{"ocs":{"meta":{"status":"failure","statuscode":997,"message":"Unauthorised"},"data":[]}}`)
	assert.NoError(t, err)
	assert.Equal(t, 997, c.OCSMeta.StatusCode)
	assert.Equal(t, "Unauthorised", c.OCSMeta.Message)
	assert.Empty(t, c.Data.Capabilities.SpreedCapabilities.Features)
}
