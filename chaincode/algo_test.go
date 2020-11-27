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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AlgoResponse struct {
	Results  []outputAlgo `json:"results"`
	Bookmark string       `json:"bookmark"`
}

func TestAlgo(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add algo with invalid field
	inpAlgo := inputAlgo{
		DescriptionChecksum: "aaa",
	}
	args := inpAlgo.createDefault()
	resp := mockStub.MockInvoke(args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding algo with invalid checksum, status %d and message %s", resp.Status, resp.Message)

	// Properly add algo
	resp, tt := registerItem(t, *mockStub, "algo")

	inpAlgo = tt.(inputAlgo)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")

	// Query algo from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryAlgo"), keyToJSON(algoKey)}
	resp = mockStub.MockInvoke(args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying an algo with status %d and message %s", resp.Status, resp.Message)
	algo := outputAlgo{}
	err = json.Unmarshal(resp.Payload, &algo)
	assert.NoError(t, err, "when unmarshalling queried objective")
	expectedAlgo := outputAlgo{
		Key:  algoKey,
		Name: inpAlgo.Name,
		Content: &ChecksumAddress{
			Checksum:       inpAlgo.Checksum,
			StorageAddress: inpAlgo.StorageAddress,
		},
		Description: &ChecksumAddress{
			Checksum:       inpAlgo.DescriptionChecksum,
			StorageAddress: inpAlgo.DescriptionStorageAddress,
		},
		Owner: worker,
		Permissions: outputPermissions{
			Process: Permission{Public: true, AuthorizedIDs: []string{}},
		},
		Metadata: map[string]string{},
	}
	assert.Exactly(t, expectedAlgo, algo)

	// Query all algo and check consistency
	args = [][]byte{[]byte("queryAlgos")}
	resp = mockStub.MockInvoke(args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying algos - status %d and message %s", resp.Status, resp.Message)
	var algos AlgoResponse
	err = json.Unmarshal(resp.Payload, &algos)
	assert.NoError(t, err, "while unmarshalling algos")
	assert.Len(t, algos.Results, 1)
	assert.Exactly(t, expectedAlgo, algos.Results[0], "return algo different from registered one")
}
