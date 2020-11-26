// Copyright 2018 Owkin, inc.
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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStub("substra", scc)
	mockTxID := "fa0f757bc278fdf6a32d00975602eb853e23a86a156781588d99ddef5b80720f"

	args := append([][]byte{[]byte("registerNode")}, []byte{})

	resp := mockStub.MockInvoke(mockTxID, args)
	assert.EqualValuesf(t, 200, resp.Status, "Node created")
	assert.Contains(t, string(resp.Payload), "\"id\":\"SampleOrg\"", "Node created")

	resp = mockStub.MockInvoke(mockTxID, args)
	assert.EqualValuesf(t, 200, resp.Status, "Node registered twice")
	assert.Contains(t, string(resp.Payload), "\"id\":\"SampleOrg\"", "Node registered twice")
}

func TestQueryNodes(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	mockTxID := "fa0f757bc278fdf6a32d00975602eb853e23a86a156781588d99ddef5b80720f"

	response := mockStub.MockInvoke(mockTxID, [][]byte{[]byte("queryNodes")})

	assert.EqualValuesf(t, 200, response.Status, "Node Created")
	assert.Contains(t, string(response.Payload), "\"id\":\"SampleOrg\"", "Query nodes")
}
