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

var (

	// Default compute plan:
	//
	//     ø     ø             ø      ø
	//     |     |             |      |
	//     hd    tr            tr     hd
	//   -----------         -----------
	//  | Composite |       | Composite |
	//   -----------         -----------
	//     hd    tr            tr     hd
	//     |      \            /      |
	//     |       \          /       |
	//     |        -----------       |
	//     |       | Aggregate |      |
	//     |        -----------       |
	//     |            |             |
	//     |     ,_____/ \_____       |
	//     |     |             |      |
	//     hd    tr            tr     hd
	//   -----------         -----------
	//  | Composite |       | Composite |
	//   -----------         -----------
	//     hd    tr            tr     hd
	//
	//
	defaultComputePlan = inputComputePlan{
		ObjectiveKey: objectiveDescriptionHash,
		TrainingTasks: []inputComputePlanTrainingTask{
			inputComputePlanTrainingTask{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash1},
				AlgoKey:        algoHash,
				ID:             traintupleID1,
			},
			inputComputePlanTrainingTask{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash2},
				ID:             traintupleID2,
				AlgoKey:        algoHash,
				InModelsIDs:    []string{traintupleID1},
			},
		},
		Testtuples: []inputComputePlanTesttuple{
			inputComputePlanTesttuple{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{testDataSampleHash1, testDataSampleHash2},
				TraintupleID:   traintupleID2,
			},
		},
	}
)

func TestCreateComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlanInternal(db, inCP)
	validateComputePlan(t, outCP, defaultComputePlan)

	// Check the traintuples
	traintuples, err := queryTraintuples(db, []string{})
	assert.NoError(t, err)
	assert.Len(t, traintuples, 2)
	require.Contains(t, outCP.TrainingTaskKeys, traintuples[0].Key)
	require.Contains(t, outCP.TrainingTaskKeys, traintuples[1].Key)
	var first, second outputTraintuple
	for _, el := range traintuples {
		switch el.Key {
		case outCP.TrainingTaskKeys[0]:
			first = el
		case outCP.TrainingTaskKeys[1]:
			second = el
		}
	}

	// check first traintuple
	assert.NotZero(t, first)
	assert.EqualValues(t, first.Key, first.ComputePlanID)
	assert.Equal(t, inCP.TrainingTasks[0].AlgoKey, first.Algo.Hash)
	assert.Equal(t, StatusTodo, first.Status)

	// check second traintuple
	assert.NotZero(t, second)
	assert.EqualValues(t, first.Key, second.InModels[0].TraintupleKey)
	assert.EqualValues(t, first.ComputePlanID, second.ComputePlanID)
	assert.Len(t, second.InModels, 1)
	assert.Equal(t, inCP.TrainingTasks[1].AlgoKey, second.Algo.Hash)
	assert.Equal(t, StatusWaiting, second.Status)

	// Check the testtuples
	testtuples, err := queryTesttuples(db, []string{})
	assert.NoError(t, err)
	require.Len(t, testtuples, 1)
	testtuple := testtuples[0]
	require.Contains(t, outCP.TesttupleKeys, testtuple.Key)
	assert.EqualValues(t, second.Key, testtuple.TraintupleKey)
	assert.True(t, testtuple.Certified)
}

func TestQueryComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlanInternal(db, inCP)
	assert.NoError(t, err)
	assert.NotNil(t, outCP)
	assert.Equal(t, outCP.ComputePlanID, outCP.TrainingTaskKeys[0])

	cp, err := queryComputePlan(db, assetToArgs(inputHash{Key: outCP.ComputePlanID}))
	assert.NoError(t, err, "calling queryComputePlan should succeed")
	assert.NotNil(t, cp)
	validateComputePlan(t, cp, defaultComputePlan)
}

func TestQueryComputePlans(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlanInternal(db, inCP)
	assert.NoError(t, err)
	assert.NotNil(t, outCP)
	assert.Equal(t, outCP.ComputePlanID, outCP.TrainingTaskKeys[0])

	cps, err := queryComputePlans(db, []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Len(t, cps, 1, "queryComputePlans should return one compute plan")
	validateComputePlan(t, cps[0], defaultComputePlan)
}

func TestComputePlanEmptyTesttuples(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	inCP := inputComputePlan{
		ObjectiveKey: objectiveDescriptionHash,
		TrainingTasks: []inputComputePlanTrainingTask{
			inputComputePlanTrainingTask{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash1},
				AlgoKey:        algoHash,
				ID:             traintupleID1,
			},
			inputComputePlanTrainingTask{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash2},
				ID:             traintupleID2,
				AlgoKey:        algoHash,
				InModelsIDs:    []string{traintupleID1},
			},
		},
		Testtuples: []inputComputePlanTesttuple{},
	}

	outCP, err := createComputePlanInternal(db, inCP)
	assert.NoError(t, err)
	assert.NotNil(t, outCP)
	assert.Equal(t, outCP.ComputePlanID, outCP.TrainingTaskKeys[0])
	assert.Equal(t, []string{}, outCP.TesttupleKeys)

	cp, err := queryComputePlan(db, assetToArgs(inputHash{Key: outCP.ComputePlanID}))
	assert.NoError(t, err, "calling queryComputePlan should succeed")
	assert.NotNil(t, cp)
	assert.Equal(t, []string{}, outCP.TesttupleKeys)

	cps, err := queryComputePlans(db, []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Len(t, cps, 1, "queryComputePlans should return one compute plan")
	assert.Equal(t, []string{}, cps[0].TesttupleKeys)
}

func TestQueryComputePlanEmpty(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	cps, err := queryComputePlans(db, []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Equal(t, []outputComputePlan{}, cps)
}

func validateComputePlan(t *testing.T, cp outputComputePlan, in inputComputePlan) {
	assert.Len(t, cp.TrainingTaskKeys, 2)
	cpID := cp.TrainingTaskKeys[0]

	assert.Equal(t, in.ObjectiveKey, cp.ObjectiveKey)

	assert.Equal(t, cpID, cp.ComputePlanID)
	assert.Equal(t, in.ObjectiveKey, cp.ObjectiveKey)

	assert.NotEmpty(t, cp.TrainingTaskKeys[0])
	assert.NotEmpty(t, cp.TrainingTaskKeys[1])

	require.Len(t, cp.TesttupleKeys, 1)
	assert.NotEmpty(t, cp.TesttupleKeys[0])
}
