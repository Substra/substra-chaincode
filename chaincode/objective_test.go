package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestRegisterObjectiveWhitoutDataset(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	inpObjective := inputObjective{TestDataset: ":"}
	args := inpObjective.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)
}
func TestRegisterObjectiveWithDataSampleKeyNotDataManagerKey(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a dataManager and some dataSample successfuly
	inpDataManager := inputDataManager{}
	args := inpDataManager.createDefault()
	mockStub.MockInvoke("42", args)
	inpDataSample := inputDataSample{
		Hashes:          testDataSampleHash1,
		DataManagerKeys: dataManagerOpenerHash,
		TestOnly:        "true",
	}
	args = inpDataSample.createDefault()
	mockStub.MockInvoke("42", args)
	inpDataSample = inputDataSample{
		Hashes:          testDataSampleHash2,
		DataManagerKeys: dataManagerOpenerHash,
		TestOnly:        "true",
	}
	args = inpDataSample.createDefault()
	r := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, r.Status)

	// Fail to insert the objective
	inpObjective := inputObjective{TestDataset: testDataSampleHash1 + ":" + testDataSampleHash2}
	args = inpObjective.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 500, resp.Status, "status should indicate an error since the dataManager key is a dataSample key")
}
func TestObjective(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add objective with invalid field
	inpObjective := inputObjective{
		DescriptionHash: "aaa",
	}
	args := inpObjective.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 500, resp.Status, "when adding objective with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Add objective with unexisting test dataSample
	inpObjective = inputObjective{}
	args = inpObjective.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 500, resp.Status, "when adding objective with unexisting test dataSample, status %d and message %s", resp.Status, resp.Message)

	// Properly add objective
	resp, tt := registerItem(t, *mockStub, "objective")

	inpObjective = tt.(inputObjective)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	objectiveKey := res["key"]
	assert.EqualValuesf(
		t,
		inpObjective.DescriptionHash,
		objectiveKey,
		"when adding objective: unexpected returned objective key - %s / %s",
		objectiveKey,
		inpObjective.DescriptionHash)

	// Query objective from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryObjective"), keyToJSON(objectiveKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying a dataManager with status %d and message %s", resp.Status, resp.Message)
	objective := outputObjective{}
	err = bytesToStruct(resp.Payload, &objective)
	assert.NoError(t, err, "when unmarshalling queried objective")
	expectedObjective := outputObjective{
		Key:   objectiveKey,
		Owner: "bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0",
		TestDataset: &Dataset{
			DataManagerKey: dataManagerOpenerHash,
			DataSampleKeys: []string{testDataSampleHash1, testDataSampleHash2},
		},
		Name: inpObjective.Name,
		Description: HashDress{
			StorageAddress: inpObjective.DescriptionStorageAddress,
			Hash:           objectiveKey,
		},
		Permissions: inpObjective.Permissions,
		Metrics: &HashDressName{
			Hash:           inpObjective.MetricsHash,
			Name:           inpObjective.MetricsName,
			StorageAddress: inpObjective.MetricsStorageAddress,
		},
	}
	assert.Exactly(t, expectedObjective, objective)

	// Query all objectives and check consistency
	args = [][]byte{[]byte("queryObjectives")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying objectives - status %d and message %s", resp.Status, resp.Message)
	var objectives []outputObjective
	err = json.Unmarshal(resp.Payload, &objectives)
	assert.NoError(t, err, "while unmarshalling objectives")
	assert.Len(t, objectives, 1)
	assert.Exactly(t, expectedObjective, objectives[0], "return objective different from registered one")
}
