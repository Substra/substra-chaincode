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
	"github.com/stretchr/testify/require"
)

func TestCreateComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	myStub := myMockStub{MockStub: mockStub}
	myStub.saveWhenWriting = true
	registerItem(t, *mockStub, "algo")
	myStub.MockTransactionStart("42")
	myStub.saveWhenWriting = false

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlan(NewLedgerDB(&myStub), assetToArgs(inCP))
	assert.NoError(t, err)
	assert.NotNil(t, outCP)
	assert.EqualValues(t, outCP.ComputePlanID, outCP.TraintupleKeys[0])

	// Save all that was written in the mocked ledger
	myStub.saveWrittenState(t)

	// Check the traintuples
	traintuples, err := queryTraintuples(NewLedgerDB(&myStub), []string{})
	assert.NoError(t, err)
	assert.Len(t, traintuples, 2)
	require.Contains(t, outCP.TraintupleKeys, traintuples[0].Key)
	require.Contains(t, outCP.TraintupleKeys, traintuples[1].Key)
	var first, second outputTraintuple
	for _, el := range traintuples {
		switch el.Key {
		case outCP.TraintupleKeys[0]:
			first = el
		case outCP.TraintupleKeys[1]:
			second = el
		}
	}
	assert.NotZero(t, first)
	assert.NotZero(t, second)
	assert.EqualValues(t, first.Key, first.ComputePlanID)
	assert.EqualValues(t, first.ComputePlanID, second.ComputePlanID)
	assert.Len(t, second.InModels, 1)
	assert.EqualValues(t, first.Key, second.InModels[0].TraintupleKey)
	assert.Equal(t, first.Status, StatusTodo)
	assert.Equal(t, second.Status, StatusWaiting)

	// Check the testtuples
	testtuples, err := queryTesttuples(NewLedgerDB(&myStub), []string{})
	assert.NoError(t, err)
	require.Len(t, testtuples, 1)
	testtuple := testtuples[0]
	require.Contains(t, outCP.TesttupleKeys, testtuple.Key)
	assert.EqualValues(t, second.Key, testtuple.Model.TraintupleKey)
	assert.True(t, testtuple.Certified)
}

func TestQueryComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	myStub := myMockStub{MockStub: mockStub}
	myStub.saveWhenWriting = true
	registerItem(t, *mockStub, "algo")
	myStub.MockTransactionStart("42")
	myStub.saveWhenWriting = false

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlan(NewLedgerDB(&myStub), assetToArgs(inCP))
	assert.NoError(t, err)
	assert.NotNil(t, outCP)
	assert.Equal(t, outCP.ComputePlanID, outCP.TraintupleKeys[0])

	// Save all that was written in the mocked ledger
	myStub.saveWrittenState(t)

	cp, err := queryComputePlan(NewLedgerDB(&myStub), assetToArgs(inputHash{Key: outCP.ComputePlanID}))
	assert.NoError(t, err, "calling queryComputePlan should succeed")
	assert.NotNil(t, cp)
	validateComputePlan(t, cp, defaultComputePlan)
}

func TestQueryComputePlans(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	myStub := myMockStub{MockStub: mockStub}
	myStub.saveWhenWriting = true
	registerItem(t, *mockStub, "algo")
	myStub.MockTransactionStart("42")
	myStub.saveWhenWriting = false

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlan(NewLedgerDB(&myStub), assetToArgs(inCP))
	assert.NoError(t, err)
	assert.NotNil(t, outCP)
	assert.Equal(t, outCP.ComputePlanID, outCP.TraintupleKeys[0])

	// Save all that was written in the mocked ledger
	myStub.saveWrittenState(t)

	cps, err := queryComputePlans(NewLedgerDB(&myStub), []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Len(t, cps, 1, "queryComputePlans should return one compute plan")
	validateComputePlan(t, cps[0], defaultComputePlan)
}

func validateComputePlan(t *testing.T, cp outputComputePlanDetails, in inputComputePlan) {
	assert.Len(t, cp.Traintuples, 2)
	cpID := cp.Traintuples[0]

	assert.Equal(t, cpID, cp.ComputePlanID)
	assert.Equal(t, in.AlgoKey, cp.AlgoKey)
	assert.Equal(t, in.ObjectiveKey, cp.ObjectiveKey)

	assert.NotEmpty(t, cp.Traintuples[0])
	assert.NotEmpty(t, cp.Traintuples[1])

	require.Len(t, cp.Testtuples, 1)
	assert.NotEmpty(t, cp.Testtuples[0])
}
