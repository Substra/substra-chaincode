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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	defaultComputePlan = inputComputePlan{
		Tag: tag,
		Traintuples: []inputComputePlanTraintuple{
			inputComputePlanTraintuple{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash1},
				AlgoKey:        algoHash,
				ID:             traintupleID1,
			},
			inputComputePlanTraintuple{
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
				ObjectiveKey:   objectiveDescriptionHash,
				TraintupleID:   traintupleID2,
			},
		},
	}

	// Model-composition compute plan:
	//
	//   ,========,                ,========,
	//   | NODE A |                | NODE B |
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
	//     |       \  | NODE C |  /       |
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
	//     hd    tr                tr     hd
	//
	//
	modelCompositionComputePlan = inputComputePlan{
		CompositeTraintuples: []inputComputePlanCompositeTraintuple{
			{
				ID:             "step_1_composite_A",
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash1},
				AlgoKey:        compositeAlgoHash,
			},
			{
				ID:             "step_1_composite_B",
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash2},
				AlgoKey:        compositeAlgoHash,
			},
			{
				ID:             "step_3_composite_A",
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash1},
				AlgoKey:        compositeAlgoHash,
				InHeadModelID:  "step_1_composite_A",
				InTrunkModelID: "step_2_aggregate",
			},
			{
				ID:             "step_3_composite_B",
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash2},
				AlgoKey:        compositeAlgoHash,
				InHeadModelID:  "step_1_composite_B",
				InTrunkModelID: "step_2_aggregate",
			},
		},
		Aggregatetuples: []inputComputePlanAggregatetuple{
			{
				ID:      "step_2_aggregate",
				AlgoKey: aggregateAlgoHash,
				InModelsIDs: []string{
					"step_1_composite_A",
					"step_1_composite_B",
				},
				Worker: worker,
			},
		},
		Testtuples: []inputComputePlanTesttuple{
			inputComputePlanTesttuple{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{testDataSampleHash1, testDataSampleHash2},
				ObjectiveKey:   objectiveDescriptionHash,
				TraintupleID:   "step_1_composite_A",
			},
			inputComputePlanTesttuple{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{testDataSampleHash1, testDataSampleHash2},
				ObjectiveKey:   objectiveDescriptionHash,
				TraintupleID:   "step_1_composite_B",
			},
			inputComputePlanTesttuple{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{testDataSampleHash1, testDataSampleHash2},
				ObjectiveKey:   objectiveDescriptionHash,
				TraintupleID:   "step_2_aggregate",
			},
			inputComputePlanTesttuple{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{testDataSampleHash1, testDataSampleHash2},
				ObjectiveKey:   objectiveDescriptionHash,
				TraintupleID:   "step_3_composite_A",
			},
			inputComputePlanTesttuple{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{testDataSampleHash1, testDataSampleHash2},
				ObjectiveKey:   objectiveDescriptionHash,
				TraintupleID:   "step_3_composite_B",
			},
		},
	}
)

func TestModelCompositionComputePlanWorkflow(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, modelCompositionComputePlan)
	assert.NoError(t, err)
	assert.NotNil(t, db.tuplesEvent)
	assert.Len(t, db.tuplesEvent.CompositeTraintuples, 2)

	// ensure the returned ranks are correct
	validateTupleRank(t, db, 0, out.CompositeTraintupleKeys[0], CompositeTraintupleType)
	validateTupleRank(t, db, 0, out.CompositeTraintupleKeys[1], CompositeTraintupleType)
	validateTupleRank(t, db, 1, out.AggregatetupleKeys[0], AggregatetupleType)
	validateTupleRank(t, db, 2, out.CompositeTraintupleKeys[2], CompositeTraintupleType)
	validateTupleRank(t, db, 2, out.CompositeTraintupleKeys[3], CompositeTraintupleType)

	_, err = logStartCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[0]}))
	assert.NoError(t, err)
	_, err = logStartCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[1]}))
	assert.NoError(t, err)

	db.tuplesEvent = &TuplesEvent{}
	inpLogCompo := inputLogSuccessCompositeTrain{}
	inpLogCompo.fillDefaults()
	inpLogCompo.Key = out.CompositeTraintupleKeys[0]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)

	inpLogCompo.Key = out.CompositeTraintupleKeys[1]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)
	assert.Len(t, db.tuplesEvent.Testtuples, 2)
	for _, test := range db.tuplesEvent.Testtuples {
		assert.Equalf(t, StatusTodo, test.Status, "blame it on %+v", test)
	}
	require.Len(t, db.tuplesEvent.Aggregatetuples, 1)
	assert.Equal(t, StatusTodo, db.tuplesEvent.Aggregatetuples[0].Status)

	_, err = logStartAggregate(db, assetToArgs(inputHash{out.AggregatetupleKeys[0]}))
	assert.NoError(t, err)

	inpLogAgg := inputLogSuccessTrain{}
	inpLogAgg.fillDefaults()
	inpLogAgg.Key = out.AggregatetupleKeys[0]
	agg, err := logSuccessAggregate(db, assetToArgs(inpLogAgg))
	assert.NoError(t, err)
	assert.Equal(t, StatusDone, agg.Status)

	_, err = logStartCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[2]}))
	assert.NoError(t, err)
	_, err = logStartCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[3]}))
	assert.NoError(t, err)

	db.tuplesEvent = &TuplesEvent{}
	inpLogCompo.Key = out.CompositeTraintupleKeys[2]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)

	inpLogCompo.Key = out.CompositeTraintupleKeys[3]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)
	assert.Len(t, db.tuplesEvent.Testtuples, 2)
	for _, test := range db.tuplesEvent.Testtuples {
		assert.Equalf(t, StatusTodo, test.Status, "blame it on %+v", test)
	}
}

func validateTupleRank(t *testing.T, db *LedgerDB, expectedRank int, key string, assetType AssetType) {
	inp := inputHash{Key: key}
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
		assert.NoError(t, fmt.Errorf("not implemented: %v", assetType))
	}
	assert.Equal(t, expectedRank, rank, "Rank for tuple of type %v with key \"%s\" should be %d", assetType, key, expectedRank)
}

func TestCreateComputePlanCompositeAggregate(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	tag := []string{"compositeTraintuple1", "compositeTraintuple2", "aggregatetuple1", "aggregatetuple2"}

	inCP := inputComputePlan{
		CompositeTraintuples: []inputComputePlanCompositeTraintuple{
			{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash1},
				AlgoKey:        compositeAlgoHash,
				ID:             tag[0],
			},
			{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash1},
				AlgoKey:        compositeAlgoHash,
				ID:             tag[1],
				InTrunkModelID: tag[0],
				InHeadModelID:  tag[0],
			},
		},
		Aggregatetuples: []inputComputePlanAggregatetuple{
			{
				AlgoKey: aggregateAlgoHash,
				ID:      tag[2],
			},
			{
				AlgoKey:     aggregateAlgoHash,
				ID:          tag[3],
				InModelsIDs: []string{tag[2]},
			},
		},
	}

	outCP, err := createComputePlanInternal(db, inCP)
	assert.NoError(t, err)

	// Check the composite traintuples
	traintuples, err := queryCompositeTraintuples(db, []string{})
	assert.NoError(t, err)
	require.Len(t, traintuples, 2)
	require.Contains(t, outCP.CompositeTraintupleKeys, traintuples[0].Key)
	require.Contains(t, outCP.CompositeTraintupleKeys, traintuples[1].Key)

	// Check the aggregate traintuples
	aggtuples, err := queryAggregatetuples(db, []string{})
	assert.NoError(t, err)
	require.Len(t, aggtuples, 2)
	require.Contains(t, outCP.AggregatetupleKeys, aggtuples[0].Key)
	require.Contains(t, outCP.AggregatetupleKeys, aggtuples[1].Key)

	// Query the compute plan
	cp, err := queryComputePlan(db, assetToArgs(inputHash{Key: outCP.ComputePlanID}))
	assert.NoError(t, err, "calling queryComputePlan should succeed")
	assert.NotNil(t, cp)
	assert.Equal(t, 2, len(cp.CompositeTraintupleKeys))
	assert.Equal(t, 2, len(cp.AggregatetupleKeys))

	// Query compute plans
	cps, err := queryComputePlans(db, []string{})
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
	outCP, err := createComputePlanInternal(db, inCP)
	assert.NoError(t, err)
	validateDefaultComputePlan(t, outCP)

	// Check the traintuples
	traintuples, err := queryTraintuples(db, []string{})
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
	assert.Equal(t, inCP.Traintuples[0].AlgoKey, first.Algo.Hash)
	assert.Equal(t, StatusTodo, first.Status)

	// check second traintuple
	assert.NotZero(t, second)
	assert.EqualValues(t, first.Key, second.InModels[0].TraintupleKey)
	assert.EqualValues(t, first.ComputePlanID, second.ComputePlanID)
	assert.Len(t, second.InModels, 1)
	assert.Equal(t, inCP.Traintuples[1].AlgoKey, second.Algo.Hash)
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

	cp, err := queryComputePlan(db, assetToArgs(inputHash{Key: outCP.ComputePlanID}))
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
	outCP, err := createComputePlanInternal(db, inCP)
	assert.NoError(t, err)
	assert.NotNil(t, outCP)

	cps, err := queryComputePlans(db, []string{})
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
		Traintuples: []inputComputePlanTraintuple{
			inputComputePlanTraintuple{
				DataManagerKey: dataManagerOpenerHash,
				DataSampleKeys: []string{trainDataSampleHash1},
				AlgoKey:        algoHash,
				ID:             traintupleID1,
			},
			inputComputePlanTraintuple{
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
	assert.Len(t, outCP.TesttupleKeys, 0)

	cp, err := queryComputePlan(db, assetToArgs(inputHash{Key: outCP.ComputePlanID}))
	assert.NoError(t, err, "calling queryComputePlan should succeed")
	assert.NotNil(t, cp)
	assert.Len(t, outCP.TesttupleKeys, 0)

	cps, err := queryComputePlans(db, []string{})
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

	cps, err := queryComputePlans(db, []string{})
	assert.NoError(t, err, "calling queryComputePlans should succeed")
	assert.Equal(t, []outputComputePlan{}, cps)
}

func TestCancelComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, modelCompositionComputePlan)
	assert.NoError(t, err)
	assert.NotNil(t, db.tuplesEvent)
	assert.Len(t, db.tuplesEvent.CompositeTraintuples, 2)

	_, err = cancelComputePlan(db, assetToArgs(inputHash{Key: out.ComputePlanID}))
	assert.NoError(t, err)

	computePlan, err := getOutComputePlan(db, out.ComputePlanID)
	assert.Equal(t, StatusCanceled, computePlan.Status)

	tuples, err := queryCompositeTraintuples(db, []string{})
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

	assert.Equal(t, nbAborted, 2)
	assert.Equal(t, nbTodo, 2)

	tests, err := queryTesttuples(db, []string{})
	assert.NoError(t, err)
	for _, test := range tests {
		assert.Equal(t, StatusAborted, test.Status)
	}
}

func TestStartedTuplesOfCanceledComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, modelCompositionComputePlan)
	assert.NoError(t, err)

	logStartCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[0]}))
	logStartCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[1]}))
	logFailCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[1]}))

	_, err = cancelComputePlan(db, assetToArgs(inputHash{Key: out.ComputePlanID}))
	assert.NoError(t, err)

	computePlan, err := getOutComputePlan(db, out.ComputePlanID)
	assert.Equal(t, StatusCanceled, computePlan.Status)

	tuples, err := queryCompositeTraintuples(db, []string{})
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
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")

	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, modelCompositionComputePlan)
	assert.NoError(t, err)

	logStartCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[0]}))
	logStartCompositeTrain(db, assetToArgs(inputHash{out.CompositeTraintupleKeys[1]}))

	_, err = cancelComputePlan(db, assetToArgs(inputHash{Key: out.ComputePlanID}))
	assert.NoError(t, err)

	inp := inputLogSuccessCompositeTrain{}
	inp.fillDefaults()
	inp.Key = out.CompositeTraintupleKeys[1]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inp))
	assert.NoError(t, err)

	computePlan, err := getOutComputePlan(db, out.ComputePlanID)
	assert.Equal(t, StatusCanceled, computePlan.Status)

	tuples, err := queryCompositeTraintuples(db, []string{})
	assert.NoError(t, err)
	expected := []string{StatusDone, StatusDoing, StatusAborted, StatusAborted}
	for i, tuple := range tuples {
		assert.Equal(t, expected[i], tuple.Status)
	}
}

func TestCreateTagedEmptyComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlan(db, assetToArgs(inputComputePlan{Tag: tag}))
	assert.NoError(t, err)
	assert.Equal(t, tag, out.Tag)
}

func TestComputePlanMetrics(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "aggregateAlgo")
	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, defaultComputePlan)
	assert.NoError(t, err)
	checkComputePlanMetrics(t, db, out.ComputePlanID, 0, 3)

	traintupleToDone(t, db, out.TraintupleKeys[0])
	checkComputePlanMetrics(t, db, out.ComputePlanID, 1, 3)

	traintupleToDone(t, db, out.TraintupleKeys[1])
	checkComputePlanMetrics(t, db, out.ComputePlanID, 2, 3)

	testtupleToDone(t, db, out.TesttupleKeys[0])
	checkComputePlanMetrics(t, db, out.ComputePlanID, 3, 3)
}

func traintupleToDone(t *testing.T, db *LedgerDB, key string) {
	_, err := logStartTrain(db, assetToArgs(inputHash{Key: key}))
	assert.NoError(t, err)

	success := inputLogSuccessTrain{}
	success.Key = key
	success.fillDefaults()
	_, err = logSuccessTrain(db, assetToArgs(success))
	assert.NoError(t, err)
}
func testtupleToDone(t *testing.T, db *LedgerDB, key string) {
	_, err := logStartTest(db, assetToArgs(inputHash{Key: key}))
	assert.NoError(t, err)

	success := inputLogSuccessTest{}
	success.Key = key
	success.createDefault()
	_, err = logSuccessTest(db, assetToArgs(success))
	assert.NoError(t, err)
}

func checkComputePlanMetrics(t *testing.T, db *LedgerDB, cpID string, doneCount, tupleCount int) {
	out, err := getOutComputePlan(db, cpID)
	assert.NoError(t, err)
	assert.Equal(t, doneCount, out.DoneCount)
	assert.Equal(t, tupleCount, out.TupleCount)
}
