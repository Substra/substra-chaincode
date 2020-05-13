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

func TestAggregateAlgo(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add algo with invalid field
	inpAlgo := inputAggregateAlgo{
		inputAlgo: inputAlgo{
			DescriptionHash: "aaa",
		},
	}
	args := inpAlgo.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding algo with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Properly add algo
	resp, tt := registerItem(t, *mockStub, "aggregateAlgo")

	inpAlgo = tt.(inputAggregateAlgo)
	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	algoKey := res.Key
	assert.Equalf(t, inpAlgo.Hash, algoKey, "when adding aggregate algo, key does not corresponds to its hash - key: %s and hash %s", algoKey, inpAlgo.Hash)

	// Query algo from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryAggregateAlgo"), keyToJSON(algoKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying an aggregate algo with status %d and message %s", resp.Status, resp.Message)
	algo := outputAggregateAlgo{}
	err = json.Unmarshal(resp.Payload, &algo)
	assert.NoError(t, err, "when unmarshalling queried aggregate  algo")
	expectedAlgo := outputAggregateAlgo{
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
	args = [][]byte{[]byte("queryAggregateAlgos")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying aggregate algos - status %d and message %s", resp.Status, resp.Message)
	var algos []outputAggregateAlgo
	err = json.Unmarshal(resp.Payload, &algos)
	assert.NoError(t, err, "while unmarshalling aggregate algos")
	assert.Len(t, algos, 1)
	assert.Exactly(t, expectedAlgo, algos[0], "return aggregate algo different from registered one")
}
