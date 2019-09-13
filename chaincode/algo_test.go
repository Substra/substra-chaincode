package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlgo(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStub("substra", scc)

	// Add algo with invalid field
	inpAlgo := inputAlgo{
		DescriptionHash: "aaa",
	}
	args := inpAlgo.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding algo with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Properly add algo
	resp, tt := registerItem(t, *mockStub, "algo")

	inpAlgo = tt.(inputAlgo)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	algoKey := res["key"]
	assert.Equalf(t, inpAlgo.Hash, algoKey, "when adding algo, key does not corresponds to its hash - key: %s and hash %s", algoKey, inpAlgo.Hash)

	// Query algo from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryAlgo"), keyToJSON(algoKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying an algo with status %d and message %s", resp.Status, resp.Message)
	algo := outputAlgo{}
	err = bytesToStruct(resp.Payload, &algo)
	assert.NoError(t, err, "when unmarshalling queried objective")
	expectedAlgo := outputAlgo{
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
	}
	assert.Exactly(t, expectedAlgo, algo)

	// Query all algo and check consistency
	args = [][]byte{[]byte("queryAlgos")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying algos - status %d and message %s", resp.Status, resp.Message)
	var algos []outputAlgo
	err = json.Unmarshal(resp.Payload, &algos)
	assert.NoError(t, err, "while unmarshalling algos")
	assert.Len(t, algos, 1)
	assert.Exactly(t, expectedAlgo, algos[0], "return algo different from registered one")
}
