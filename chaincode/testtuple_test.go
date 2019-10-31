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

func TestTesttupleOnFailedTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	resp, _ := registerItem(t, *mockStub, "traintuple")

	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	traintupleKey := res["key"]

	// Mark the traintuple as failed
	fail := inputLogFailTrain{}
	fail.Key = traintupleKey
	args := fail.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "should be able to log traintuple as failed")

	// Fail to add a testtuple to this failed traintuple
	inpTesttuple := inputTesttuple{}
	args = inpTesttuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, "status should show an error since the traintuple is failed")
	assert.Contains(t, resp.Message, "could not register this testtuple")
}

func TestCertifiedExplicitTesttuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "traintuple")

	// Add a testtuple that shoulb be certified since it's the same dataManager and
	// dataSample than the objective but explicitly pass as arguments and in disorder
	inpTesttuple := inputTesttuple{
		DataSampleKeys: []string{testDataSampleHash2, testDataSampleHash1},
		DataManagerKey: dataManagerOpenerHash}
	args := inpTesttuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	args = [][]byte{[]byte("queryTesttuples")}
	resp = mockStub.MockInvoke("42", args)
	testtuples := [](map[string]interface{}){}
	err := json.Unmarshal(resp.Payload, &testtuples)
	assert.NoError(t, err, "should be unmarshaled")
	assert.Len(t, testtuples, 1, "there should be only one testtuple...")
	assert.True(t, testtuples[0]["certified"].(bool), "... and it should be certified")

}

func TestConflictCertifiedNonCertifiedTesttuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "traintuple")

	// Add a certified testtuple
	inpTesttuple1 := inputTesttuple{}
	args := inpTesttuple1.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	// Fail to add an incomplete uncertified testtuple
	inpTesttuple2 := inputTesttuple{DataSampleKeys: []string{trainDataSampleHash1}}
	args = inpTesttuple2.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status)
	assert.Contains(t, resp.Message, "invalid input: dataManagerKey and dataSampleKey should be provided together")

	// Add an uncertified testtuple successfully
	inpTesttuple3 := inputTesttuple{
		DataSampleKeys: []string{trainDataSampleHash1, trainDataSampleHash2},
		DataManagerKey: dataManagerOpenerHash}
	args = inpTesttuple3.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	// Fail to add the same testtuple with a different order for dataSampleKeys
	inpTesttuple4 := inputTesttuple{
		DataSampleKeys: []string{trainDataSampleHash2, trainDataSampleHash1},
		DataManagerKey: dataManagerOpenerHash}
	args = inpTesttuple4.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 409, resp.Status)
	assert.Contains(t, resp.Message, "already exists")
}
