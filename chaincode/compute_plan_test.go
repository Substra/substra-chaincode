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
	defaultComputePlan = inputComputePlan{
		Key: computePlanKey,
		Traintuples: []inputComputePlanTraintuple{
			inputComputePlanTraintuple{
				Key:            computePlanTraintupleKey1,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        algoKey,
				ID:             traintupleID1,
			},
			inputComputePlanTraintuple{
				Key:            computePlanTraintupleKey2,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey2},
				ID:             traintupleID2,
				AlgoKey:        algoKey,
				InModelsIDs:    []string{traintupleID1},
			},
		},
		Testtuples: []inputComputePlanTesttuple{
			inputComputePlanTesttuple{
				Key:            computePlanTesttupleKey1,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   traintupleID2,
			},
		},
	}

	// Model-composition compute plan:
	//
	//   ,========,                ,========,
	//   | ORG A  |                | ORG B  |
	//   *========*                *========*
	//
	//     ø     ø                  ø      ø
	//     |     |                  |      |
	//     hd    tr                 tr     hd
	//   -----------              -----------
	//  | Composite |            | Composite |      STEP 1
	//   -----------              -----------
	//     hd    tr                 tr     hd
	//     |      \   ,========,   /      |
	//     |       \  | ORG C  |  /       |
	//     |        \ *========* /        |
	//     |       ----------------       |
	//     |      |    Aggregate   |      |         STEP 2
	//     |       ----------------       |
	//     |              |               |
	//     |     ,_______/ \_______       |
	//     |     |                 |      |
	//     hd    tr                tr     hd
	//   -----------             -----------
	//  | Composite |           | Composite |       STEP 3
	//   -----------             -----------
	//     hd    tr                 tr     hd
	//            \                /
	//             \              /
	//              \            /
	//             ----------------
	//            |    Aggregate   |                STEP 4
	//             ----------------
	//
	//
	modelCompositionComputePlan = inputComputePlan{
		Key: computePlanKey,
		CompositeTraintuples: []inputComputePlanCompositeTraintuple{
			{
				Key:            computePlanCompositeTraintupleKey1,
				ID:             "step_1_composite_A",
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        compositeAlgoKey,
			},
			{
				Key:            computePlanCompositeTraintupleKey2,
				ID:             "step_1_composite_B",
				DataManagerKey: dataManagerKey2,
				DataSampleKeys: []string{trainDataSampleKeyWorker2},
				AlgoKey:        compositeAlgoKey,
			},
			{
				Key:            computePlanCompositeTraintupleKey3,
				ID:             "step_3_composite_A",
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        compositeAlgoKey,
				InHeadModelID:  "step_1_composite_A",
				InTrunkModelID: "step_2_aggregate",
			},
			{
				Key:            computePlanCompositeTraintupleKey4,
				ID:             "step_3_composite_B",
				DataManagerKey: dataManagerKey2,
				DataSampleKeys: []string{trainDataSampleKeyWorker2},
				AlgoKey:        compositeAlgoKey,
				InHeadModelID:  "step_1_composite_B",
				InTrunkModelID: "step_2_aggregate",
			},
		},
		Aggregatetuples: []inputComputePlanAggregatetuple{
			{
				Key:     computePlanAggregatetupleKey1,
				ID:      "step_2_aggregate",
				AlgoKey: aggregateAlgoKey,
				InModelsIDs: []string{
					"step_1_composite_A",
					"step_1_composite_B",
				},
				Worker: workerC,
			}, {
				Key:     computePlanAggregatetupleKey2,
				ID:      "step_4_aggregate",
				AlgoKey: aggregateAlgoKey,
				InModelsIDs: []string{
					"step_3_composite_A",
					"step_3_composite_B",
				},
				Worker: workerC,
			},
		},
		Testtuples: []inputComputePlanTesttuple{
			{
				Key:            computePlanTesttupleKey1,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_1_composite_A",
			},
			{
				Key:            computePlanTesttupleKey2,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_1_composite_B",
			},
			{
				Key:            computePlanTesttupleKey3,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_2_aggregate",
			},
			{
				Key:            computePlanTesttupleKey4,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_3_composite_A",
			},
			{
				Key:            computePlanTesttupleKey5,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_3_composite_B",
			},
			{
				Key:            computePlanTesttupleKey6,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_4_aggregate",
			},
		},
	}
)

type TestModels struct {
	composite []TestCompositeModel
	Aggregate string
}

type TestCompositeModel struct {
	Head  string
	Trunk string
}

func TestModelCompositionComputePlanWorkflow(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := getMockStubForModelComposition(t, scc)

	// hack to be able to access internal functions directly
	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	// Create CP
	out, err := createComputePlanInternal(db, modelCompositionComputePlan, tag, map[string]string{}, true)
	assert.NoError(t, err)
	assert.NotNil(t, db.event)
	assert.Len(t, db.event.CompositeTraintuples, 2)

	// ensure the returned ranks are correct
	validateTupleRank(t, db, 0, out.CompositeTraintupleKeys[0], CompositeTraintupleType)
	validateTupleRank(t, db, 0, out.CompositeTraintupleKeys[1], CompositeTraintupleType)
	validateTupleRank(t, db, 1, out.AggregatetupleKeys[0], AggregatetupleType)
	validateTupleRank(t, db, 2, out.CompositeTraintupleKeys[2], CompositeTraintupleType)
	validateTupleRank(t, db, 2, out.CompositeTraintupleKeys[3], CompositeTraintupleType)

	// Generate some random out-model hashes
	step := map[int]TestModels{
		1: {composite: []TestCompositeModel{
			{Head: RandomUUID(), Trunk: RandomUUID()},
			{Head: RandomUUID(), Trunk: RandomUUID()}}},
		2: {Aggregate: RandomUUID()},
		3: {composite: []TestCompositeModel{
			{Head: RandomUUID(), Trunk: RandomUUID()},
			{Head: RandomUUID(), Trunk: RandomUUID()}}},
		4: {Aggregate: RandomUUID()},
	}

	// Step 1
	compositeToDone(t, mockStub, workerA, db, out.CompositeTraintupleKeys[0], step[1].composite[0].Head, step[1].composite[0].Trunk)
	assert.Len(t, db.event.Testtuples, 1)
	assert.Equal(t, StatusTodo, db.event.Testtuples[0].Status)

	compositeToDone(t, mockStub, workerB, db, out.CompositeTraintupleKeys[1], step[1].composite[1].Head, step[1].composite[1].Trunk)
	assert.Len(t, db.event.Testtuples, 1)
	assert.Len(t, db.event.Aggregatetuples, 1)
	assert.Equal(t, StatusTodo, db.event.Testtuples[0].Status)
	assert.Equal(t, StatusTodo, db.event.Aggregatetuples[0].Status)

	assert.Len(t, db.event.ComputePlans, 0)

	testtupleToDone(t, db, out.TesttupleKeys[0])
	testtupleToDone(t, db, out.TesttupleKeys[1])

	// Step 2
	aggregateToDone(t, mockStub, workerC, db, out.AggregatetupleKeys[0], step[2].Aggregate)
	testtupleToDone(t, db, out.TesttupleKeys[2])
	assert.Len(t, db.event.ComputePlans, 0)

	// Step 3
	compositeToDone(t, mockStub, workerA, db, out.CompositeTraintupleKeys[2], step[3].composite[0].Head, step[3].composite[0].Trunk)
	assert.Len(t, db.event.Testtuples, 1)
	assert.Equal(t, StatusTodo, db.event.Testtuples[0].Status)
	assert.Len(t, db.event.ComputePlans, 1)
	assert.Len(t, db.event.ComputePlans[0].ModelsToDelete, 2)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[1].composite[0].Head)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[1].composite[0].Trunk)

	compositeToDone(t, mockStub, workerB, db, out.CompositeTraintupleKeys[3], step[3].composite[1].Head, step[3].composite[1].Trunk)
	assert.Len(t, db.event.Testtuples, 1)
	assert.Equal(t, StatusTodo, db.event.Testtuples[0].Status)
	assert.Len(t, db.event.ComputePlans, 1)
	assert.Len(t, db.event.ComputePlans[0].ModelsToDelete, 2)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[1].composite[1].Head)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[1].composite[1].Trunk)

	testtupleToDone(t, db, out.TesttupleKeys[3])
	testtupleToDone(t, db, out.TesttupleKeys[4])

	// Step 4
	aggregateToDone(t, mockStub, workerC, db, out.AggregatetupleKeys[1], step[4].Aggregate)
	assert.Len(t, db.event.ComputePlans, 1)
	assert.Len(t, db.event.ComputePlans[0].ModelsToDelete, 1)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[2].Aggregate)

	testtupleToDone(t, db, out.TesttupleKeys[5])
	assert.Len(t, db.event.ComputePlans, 1)
	assert.Len(t, db.event.ComputePlans[0].ModelsToDelete, 4)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[3].composite[0].Head)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[3].composite[0].Trunk)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[3].composite[1].Head)
	assert.Contains(t, db.event.ComputePlans[0].ModelsToDelete, step[3].composite[1].Trunk)
}

func validateTupleRank(t *testing.T, db *LedgerDB, expectedRank int, key string, assetType AssetType) {
	inp := inputKey{Key: key}
	rank := -42
	switch assetType {
	case CompositeTraintupleType:
		tuple, err := queryCompositeTraintuple(db, assetToArgs(inp))
		assert.NoError(t, err)
		rank = tuple.Rank
	case AggregatetupleType:
		tuple, err := queryAggregatetuple(db, assetToArgs(inp))
		assert.NoError(t, err)
		rank = tuple.Rank
	default:
		t.Errorf("not implemented: %s", assetType)
	}
	assert.Equal(t, expectedRank, rank, "Rank for tuple of type %v with key \"%s\" should be %d", assetType, key, expectedRank)
}

func TestCreateComputePlanCompositeAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	IDs := []string{"compositeTraintuple1", "compositeTraintuple2", "aggregatetuple1", "aggregatetuple2"}

	inCP := inputComputePlan{
		Key: computePlanKey,
		CompositeTraintuples: []inputComputePlanCompositeTraintuple{
			{
				Key:            computePlanCompositeTraintupleKey1,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        compositeAlgoKey,
				ID:             IDs[0],
			},
			{
				Key:            computePlanCompositeTraintupleKey2,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        compositeAlgoKey,
				ID:             IDs[1],
				InTrunkModelID: IDs[0],
				InHeadModelID:  IDs[0],
			},
		},
		Aggregatetuples: []inputComputePlanAggregatetuple{
			{
				Key:     computePlanAggregatetupleKey1,
				AlgoKey: aggregateAlgoKey,
				ID:      IDs[2],
				Worker:  workerA,
			},
			{
				Key:         computePlanAggregatetupleKey2,
				AlgoKey:     aggregateAlgoKey,
				ID:          IDs[3],
				InModelsIDs: []string{IDs[2]},
				Worker:      workerA,
			},
		},
	}

	outCP, err := createComputePlanInternal(db, inCP, tag, map[string]string{}, false)
	assert.NoError(t, err)

	// Check the composite traintuples
	traintuples, _, err := queryCompositeTraintuples(db, []string{})
	assert.NoError(t, err)
	require.Len(t, traintuples, 2)
	require.Contains(t, outCP.CompositeTraintupleKeys, traintuples[0].Key)
	require.Contains(t, outCP.CompositeTraintupleKeys, traintuples[1].Key)

	// Check the aggregate traintuples
	aggtuples, _, err := queryAggregatetuples(db, []string{})
	assert.NoError(t, err)
	require.Len(t, aggtuples, 2)
	require.Contains(t, outCP.AggregatetupleKeys, aggtuples[0].Key)
	require.Contains(t, outCP.AggregatetupleKeys, aggtuples[1].Key)

	// Query the compute plan
	cp, err := queryComputePlan(db, assetToArgs(inputKey{Key: outCP.Key}))
	assert.NoError(t, err, "calling queryComputePlan should succeed")
	assert.NotNil(t, cp)
	assert.Equal(t, 2, len(cp.CompositeTraintupleKeys))
	assert.Equal(t, 2, len(cp.AggregatetupleKeys))

	// Query compute plans
	cps, _, err := queryComputePlans(db, []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Len(t, cps, 1, "queryComputePlans should return one compute plan")
	assert.Equal(t, 2, len(cps[0].CompositeTraintupleKeys))
	assert.Equal(t, 2, len(cps[0].AggregatetupleKeys))
}

func TestCreateComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlanInternal(db, inCP, tag, map[string]string{}, false)
	assert.NoError(t, err)
	validateDefaultComputePlan(t, outCP)

	// Check the traintuples
	traintuples, _, err := queryTraintuples(db, []string{})
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

	// check first traintuple
	assert.NotZero(t, first)
	assert.Equal(t, inCP.Traintuples[0].AlgoKey, first.Algo.Key)
	algo1, err := queryAlgo(db, assetToArgs(inputKey{Key: inCP.Traintuples[0].AlgoKey}))
	assert.NoError(t, err)
	assert.Equal(t, algo1.Content.Checksum, first.Algo.Checksum)
	assert.Equal(t, StatusTodo, first.Status)

	// check second traintuple
	assert.NotZero(t, second)
	assert.EqualValues(t, first.Key, second.InModels[0].TraintupleKey)
	assert.EqualValues(t, first.ComputePlanKey, second.ComputePlanKey)
	assert.Len(t, second.InModels, 1)
	assert.Equal(t, inCP.Traintuples[1].AlgoKey, second.Algo.Key)
	algo2, err := queryAlgo(db, assetToArgs(inputKey{Key: inCP.Traintuples[1].AlgoKey}))
	assert.NoError(t, err)
	assert.Equal(t, algo2.Content.Checksum, second.Algo.Checksum)
	assert.Equal(t, StatusWaiting, second.Status)

	// Check the testtuples
	testtuples, _, err := queryTesttuples(db, []string{})
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
	outCP, err := createComputePlanInternal(db, inCP, tag, map[string]string{}, false)
	assert.NoError(t, err)
	assert.NotNil(t, outCP)

	cp, err := queryComputePlan(db, assetToArgs(inputKey{Key: outCP.Key}))
	assert.NoError(t, err, "calling queryComputePlan should succeed")
	assert.NotNil(t, cp)
	validateDefaultComputePlan(t, cp)
}

func TestQueryComputePlans(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlanInternal(db, inCP, tag, map[string]string{}, false)
	assert.NoError(t, err)
	assert.NotNil(t, outCP)

	cps, _, err := queryComputePlans(db, []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Len(t, cps, 1, "queryComputePlans should return one compute plan")
	validateDefaultComputePlan(t, cps[0])
}

func validateDefaultComputePlan(t *testing.T, cp outputComputePlan) {
	assert.Equal(t, tag, cp.Tag)
	assert.Len(t, cp.TraintupleKeys, 2)

	assert.NotEmpty(t, cp.TraintupleKeys[0])
	assert.NotEmpty(t, cp.TraintupleKeys[1])

	require.Len(t, cp.TesttupleKeys, 1)
	assert.NotEmpty(t, cp.TesttupleKeys[0])
}

func TestComputePlanEmptyTesttuples(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	inCP := inputComputePlan{
		Key: computePlanKey,
		Traintuples: []inputComputePlanTraintuple{
			inputComputePlanTraintuple{
				Key:            computePlanTraintupleKey1,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        algoKey,
				ID:             traintupleID1,
			},
			inputComputePlanTraintuple{
				Key:            computePlanTraintupleKey2,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey2},
				ID:             traintupleID2,
				AlgoKey:        algoKey,
				InModelsIDs:    []string{traintupleID1},
			},
		},
		Testtuples: []inputComputePlanTesttuple{},
	}

	outCP, err := createComputePlanInternal(db, inCP, tag, map[string]string{}, false)
	assert.NoError(t, err)
	assert.NotNil(t, outCP)
	assert.Len(t, outCP.TesttupleKeys, 0)

	cp, err := queryComputePlan(db, assetToArgs(inputKey{Key: outCP.Key}))
	assert.NoError(t, err, "calling queryComputePlan should succeed")
	assert.NotNil(t, cp)
	assert.Len(t, outCP.TesttupleKeys, 0)

	cps, _, err := queryComputePlans(db, []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Len(t, cps, 1, "queryComputePlans should return one compute plan")
	assert.Len(t, cps[0].TesttupleKeys, 0)
}

func TestQueryComputePlanEmpty(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "algo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	cps, _, err := queryComputePlans(db, []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Equal(t, []outputComputePlan{}, cps)
}

func TestCancelComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, defaultComputePlan, tag, map[string]string{}, false)
	assert.NoError(t, err)

	_, err = cancelComputePlan(db, assetToArgs(inputKey{Key: out.Key}))
	assert.NoError(t, err)

	computePlan, err := getOutComputePlan(db, out.Key)
	assert.Equal(t, StatusCanceled, computePlan.Status)

	tuples, _, err := queryTraintuples(db, []string{})
	assert.NoError(t, err)

	nbAborted, nbTodo := 0, 0
	for _, tuple := range tuples {
		if tuple.Status == StatusAborted {
			nbAborted = nbAborted + 1
		}
		if tuple.Status == StatusTodo {
			nbTodo = nbTodo + 1
		}
	}

	assert.Equal(t, nbAborted, 1)
	assert.Equal(t, nbTodo, 1)

	tests, _, err := queryTesttuples(db, []string{})
	assert.NoError(t, err)
	for _, test := range tests {
		assert.Equal(t, StatusAborted, test.Status)
	}
}

func TestStartedTuplesOfCanceledComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := getMockStubForModelComposition(t, scc)

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, modelCompositionComputePlan, tag, map[string]string{}, false)
	assert.NoError(t, err)

	logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[0]}))
	logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[1]}))
	logFailCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[1]}))

	_, err = cancelComputePlan(db, assetToArgs(inputKey{Key: out.Key}))
	assert.NoError(t, err)

	computePlan, err := getOutComputePlan(db, out.Key)
	assert.Equal(t, StatusCanceled, computePlan.Status)

	tuples, _, err := queryCompositeTraintuples(db, []string{})
	assert.NoError(t, err)
	for _, tuple := range tuples {
		if tuple.Rank == 0 {
			assert.NotEqual(t, StatusAborted, tuple.Status, tuple.Rank)
			continue
		}
		assert.Equal(t, StatusAborted, tuple.Status, tuple.Rank)
	}
}

func TestLogSuccessAfterCancel(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := getMockStubForModelComposition(t, scc)

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, modelCompositionComputePlan, tag, map[string]string{}, false)
	assert.NoError(t, err)

	logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[0]}))

	mockStub.Creator = workerB // log start as org B
	logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[1]}))
	mockStub.Creator = workerA

	_, err = cancelComputePlan(db, assetToArgs(inputKey{Key: out.Key}))
	assert.NoError(t, err)

	inp := inputLogSuccessCompositeTrain{}
	inp.fillDefaults()
	inp.Key = out.CompositeTraintupleKeys[1]

	mockStub.Creator = workerB // log success as org B
	_, err = logSuccessCompositeTrain(db, assetToArgs(inp))
	mockStub.Creator = workerA

	assert.NoError(t, err)

	computePlan, err := getOutComputePlan(db, out.Key)
	assert.Equal(t, StatusCanceled, computePlan.Status)

	expected := []string{StatusDoing, StatusDone, StatusAborted, StatusAborted}
	for i, tuplekey := range out.CompositeTraintupleKeys {
		tuple, err := queryCompositeTraintuple(db, keyToArgs(tuplekey))
		assert.NoError(t, err)
		assert.Equal(t, expected[i], tuple.Status)
	}
}

func TestCreateTagedEmptyComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	inp := inputNewComputePlan{
		inputComputePlan: inputComputePlan{
			Key: computePlanKey,
		},
		Tag: tag}
	out, err := createComputePlan(db, assetToArgs(inp))
	assert.NoError(t, err)
	assert.Equal(t, tag, out.Tag)
}

func TestComputePlanMetrics(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")
	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, defaultComputePlan, tag, map[string]string{}, false)
	assert.NoError(t, err)
	checkComputePlanMetrics(t, db, out.Key, 0, 3)

	traintupleToDone(t, db, out.TraintupleKeys[0])
	checkComputePlanMetrics(t, db, out.Key, 1, 3)

	traintupleToDone(t, db, out.TraintupleKeys[1])
	checkComputePlanMetrics(t, db, out.Key, 2, 3)

	testtupleToDone(t, db, out.TesttupleKeys[0])
	checkComputePlanMetrics(t, db, out.Key, 3, 3)
}

func traintupleToDone(t *testing.T, db *LedgerDB, key string) {
	_, err := logStartTrain(db, assetToArgs(inputKey{Key: key}))
	assert.NoError(t, err)
	clearEvent(db)

	success := inputLogSuccessTrain{}
	success.Key = key
	success.OutModel.Key = RandomUUID()
	success.OutModel.Checksum = GetRandomHash()
	success.fillDefaults()
	_, err = logSuccessTrain(db, assetToArgs(success))
	assert.NoError(t, err)
}

func testtupleToDone(t *testing.T, db *LedgerDB, key string) {
	_, err := logStartTest(db, assetToArgs(inputKey{Key: key}))
	assert.NoError(t, err)
	clearEvent(db)

	success := inputLogSuccessTest{}
	success.Key = key
	success.createDefault()
	_, err = logSuccessTest(db, assetToArgs(success))
	assert.NoError(t, err)
}

func checkComputePlanMetrics(t *testing.T, db *LedgerDB, cpKey string, doneCount, tupleCount int) {
	out, err := getOutComputePlan(db, cpKey)
	assert.NoError(t, err)
	assert.Equal(t, doneCount, out.DoneCount)
	assert.Equal(t, tupleCount, out.TupleCount)
}

func TestUpdateComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	mockStub.MockTransactionStart("42")
	registerItem(t, *mockStub, "aggregateAlgo")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, inputComputePlan{Key: computePlanKey}, tag, map[string]string{}, false)
	assert.NoError(t, err)
	assert.Equal(t, tag, out.Tag)

	inp := defaultComputePlan
	inp.Key = out.Key
	out, err = updateComputePlanInternal(db, inp)
	assert.NoError(t, err)
	validateDefaultComputePlan(t, out)
	for _, train := range defaultComputePlan.Traintuples {
		assert.Contains(t, out.IDToKey, train.ID)
	}

	NewID := "Update"
	up := inputComputePlan{
		Key: out.Key,
		Traintuples: []inputComputePlanTraintuple{
			{
				Key:            computePlanTraintupleKey3,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        algoKey,
				ID:             NewID,
				InModelsIDs:    []string{traintupleID1, traintupleID2},
			},
		},
	}
	out, err = updateComputePlanInternal(db, up)
	assert.NoError(t, err)
	assert.Contains(t, out.IDToKey, NewID)
	assert.Len(
		t,
		out.IDToKey,
		1,
		"IDToKey should match the newly created tuple keys to its ID")
	assert.Equal(t, 4, out.TupleCount)
}

// When the smart contracts are called directly the event object is never reset
// so we need to empty it by hand after each transaction when testing the event content
func clearEvent(db *LedgerDB) {
	if db.event == nil {
		return
	}
	*(db.event) = Event{}
}

func TestCreateSameComputePlanTwice(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	mockStub.MockTransactionStart("42")
	registerItem(t, *mockStub, "aggregateAlgo")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, inputComputePlan{Key: computePlanKey}, tag, map[string]string{}, false)
	assert.NoError(t, err)
	assert.Equal(t, tag, out.Tag)

	up := inputComputePlan{
		Key: out.Key,
		Traintuples: []inputComputePlanTraintuple{
			{
				Key:            computePlanTraintupleKey3,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        algoKey,
				ID:             "traintuple",
			},
		},
		CompositeTraintuples: []inputComputePlanCompositeTraintuple{
			{
				Key:            computePlanCompositeTraintupleKey1,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey1},
				AlgoKey:        compositeAlgoKey,
				ID:             "CompositeTraintuple",
			},
		},
		Aggregatetuples: []inputComputePlanAggregatetuple{
			{
				Key:     computePlanAggregatetupleKey1,
				AlgoKey: aggregateAlgoKey,
				ID:      "Aggregatetuple",
				Worker:  workerA,
			},
		},
		Testtuples: []inputComputePlanTesttuple{
			{
				Key:            computePlanTesttupleKey2,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "traintuple",
			},
		},
	}
	out, err = updateComputePlanInternal(db, up)
	assert.NoError(t, err)

	// Upload the same tuples inside another compute plan
	out, err = createComputePlanInternal(db, inputComputePlan{Key: computePlanKey2}, tag, map[string]string{}, false)
	assert.NoError(t, err)
	assert.Equal(t, tag, out.Tag)

	inp := defaultComputePlan
	inp.Key = out.Key
	out, err = updateComputePlanInternal(db, inp)
	assert.NoError(t, err)
}

/////////////////////////////////
//                             //
// Helper types and functions  //
//                             //
/////////////////////////////////

func getMockStubForModelComposition(t *testing.T, scc *SubstraChaincode) *MockStub {
	mockStub := NewMockStub("substra", scc)
	registerWorker(mockStub, workerA)
	registerWorker(mockStub, workerB)
	registerWorker(mockStub, workerC)
	registerItem(t, *mockStub, "aggregateAlgo")

	// Add data manager and data samples for node B
	savedCreator := mockStub.Creator
	mockStub.Creator = workerB
	registerTestDataManager(t, mockStub, dataManagerKey2, trainDataSampleKeyWorker2)
	mockStub.Creator = savedCreator
	return mockStub
}

func registerTestDataManager(t *testing.T, mockStub *MockStub, key string, datasampleKeys ...string) {
	inpDataManager := inputDataManager{Key: key}
	args := inpDataManager.createDefault()
	resp := mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding dataManager with status %d and message %s", resp.Status, resp.Message)

	inpDataSample := inputDataSample{
		Keys:            datasampleKeys,
		DataManagerKeys: []string{key},
		TestOnly:        "false",
	}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding test dataSample with status %d and message %s", resp.Status, resp.Message)
}

func registerWorker(mockStub *MockStub, worker string) {
	savedCreator := mockStub.Creator
	mockStub.Creator = worker
	mockStub.MockInvoke([][]byte{[]byte("registerNode")})
	mockStub.Creator = savedCreator
}

func compositeToDone(t *testing.T, mockStub *MockStub, worker string, db *LedgerDB, key string, headModelKey string, trunkModelKey string) {
	mockStub.Creator = worker

	_, err := logStartCompositeTrain(db, assetToArgs(inputKey{key}))
	assert.NoError(t, err)
	clearEvent(db)

	inpLogCompo := inputLogSuccessCompositeTrain{}
	inpLogCompo.fillDefaults()
	inpLogCompo.OutHeadModel.Key = headModelKey
	inpLogCompo.OutTrunkModel.Key = trunkModelKey
	inpLogCompo.Key = key
	comp, err := logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)
	assert.Equal(t, StatusDone, comp.Status)

	mockStub.Creator = workerA // reset worker to default
}

func aggregateToDone(t *testing.T, mockStub *MockStub, worker string, db *LedgerDB, key string, modelKey string) {
	mockStub.Creator = worker

	_, err := logStartAggregate(db, assetToArgs(inputKey{key}))
	assert.NoError(t, err)
	clearEvent(db)

	inpLogAgg := inputLogSuccessTrain{}
	inpLogAgg.fillDefaults()
	inpLogAgg.OutModel.Key = modelKey
	inpLogAgg.Key = key
	agg, err := logSuccessAggregate(db, assetToArgs(inpLogAgg))
	assert.NoError(t, err)
	assert.Equal(t, StatusDone, agg.Status)

	mockStub.Creator = workerA // reset worker to default
}
