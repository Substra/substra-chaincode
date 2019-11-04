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

func TestCompositeAlgo(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add algo with invalid field
	inpAlgo := inputCompositeAlgo{
		inputAlgo: inputAlgo{
			DescriptionHash: "aaa",
		},
	}
	args := inpAlgo.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding algo with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Properly add algo
	resp, tt := registerItem(t, *mockStub, "compositealgo")

	inpAlgo = tt.(inputCompositeAlgo)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	algoKey := res["key"]
	assert.Equalf(t, inpAlgo.Hash, algoKey, "when adding algo, key does not corresponds to its hash - key: %s and hash %s", algoKey, inpAlgo.Hash)

	// Query algo from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryCompositeAlgo"), keyToJSON(algoKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying an algo with status %d and message %s", resp.Status, resp.Message)
	algo := outputCompositeAlgo{}
	err = bytesToStruct(resp.Payload, &algo)
	assert.NoError(t, err, "when unmarshalling queried objective")
	expectedAlgo := outputCompositeAlgo{
		outputAlgo: outputAlgo{
			Key:  algoKey,
			Name: inpAlgo.Name,
			Content: HashDress{
				Hash:           algoKey,
				StorageAddress: inpAlgo.StorageAddress,
			},
			Description: &HashDress{
				Hash:           inpAlgo.DescriptionHash,
				StorageAddress: inpAlgo.DescriptionStorageAddress,
			},
			Owner: worker,
			Permissions: outputPermissions{
				Process: Permission{Public: true, AuthorizedIDs: []string{}},
			},
		},
	}
	assert.Exactly(t, expectedAlgo, algo)

	// Query all algo and check consistency
	args = [][]byte{[]byte("queryCompositeAlgos")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying algos - status %d and message %s", resp.Status, resp.Message)
	var algos []outputCompositeAlgo
	err = json.Unmarshal(resp.Payload, &algos)
	assert.NoError(t, err, "while unmarshalling algos")
	assert.Len(t, algos, 1)
	assert.Exactly(t, expectedAlgo, algos[0], "return algo different from registered one")
}
