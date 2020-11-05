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
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey2},
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
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{trainDataSampleKey2},
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
				Worker: worker,
			},
		},
		Testtuples: []inputComputePlanTesttuple{
			inputComputePlanTesttuple{
				Key:            computePlanTesttupleKey1,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_1_composite_A",
			},
			inputComputePlanTesttuple{
				Key:            computePlanTesttupleKey2,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_1_composite_B",
			},
			inputComputePlanTesttuple{
				Key:            computePlanTesttupleKey3,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_2_aggregate",
			},
			inputComputePlanTesttuple{
				Key:            computePlanTesttupleKey4,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
				TraintupleID:   "step_3_composite_A",
			},
			inputComputePlanTesttuple{
				Key:            computePlanTesttupleKey5,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: []string{testDataSampleKey1, testDataSampleKey2},
				ObjectiveKey:   objectiveKey,
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

	out, err := createComputePlanInternal(db, modelCompositionComputePlan, tag, map[string]string{}, false)
	assert.NoError(t, err)
	assert.NotNil(t, db.event)
	assert.Len(t, db.event.CompositeTraintuples, 2)

	// ensure the returned ranks are correct
	validateTupleRank(t, db, 0, out.CompositeTraintupleKeys[0], CompositeTraintupleType)
	validateTupleRank(t, db, 0, out.CompositeTraintupleKeys[1], CompositeTraintupleType)
	validateTupleRank(t, db, 1, out.AggregatetupleKeys[0], AggregatetupleType)
	validateTupleRank(t, db, 2, out.CompositeTraintupleKeys[2], CompositeTraintupleType)
	validateTupleRank(t, db, 2, out.CompositeTraintupleKeys[3], CompositeTraintupleType)

	_, err = logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[0]}))
	assert.NoError(t, err)
	_, err = logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[1]}))
	assert.NoError(t, err)

	db.event = &Event{}
	inpLogCompo := inputLogSuccessCompositeTrain{}
	inpLogCompo.fillDefaults()
	inpLogCompo.Key = out.CompositeTraintupleKeys[0]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)

	inpLogCompo.Key = out.CompositeTraintupleKeys[1]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)
	assert.Len(t, db.event.Testtuples, 2)
	for _, test := range db.event.Testtuples {
		assert.Equalf(t, StatusTodo, test.Status, "blame it on %+v", test)
	}
	require.Len(t, db.event.Aggregatetuples, 1)
	assert.Equal(t, StatusTodo, db.event.Aggregatetuples[0].Status)

	_, err = logStartAggregate(db, assetToArgs(inputKey{out.AggregatetupleKeys[0]}))
	assert.NoError(t, err)

	inpLogAgg := inputLogSuccessTrain{}
	inpLogAgg.fillDefaults()
	inpLogAgg.Key = out.AggregatetupleKeys[0]
	agg, err := logSuccessAggregate(db, assetToArgs(inpLogAgg))
	assert.NoError(t, err)
	assert.Equal(t, StatusDone, agg.Status)

	_, err = logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[2]}))
	assert.NoError(t, err)
	_, err = logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[3]}))
	assert.NoError(t, err)

	db.event = &Event{}
	inpLogCompo.Key = out.CompositeTraintupleKeys[2]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)

	inpLogCompo.Key = out.CompositeTraintupleKeys[3]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inpLogCompo))
	assert.NoError(t, err)
	assert.Len(t, db.event.Testtuples, 2)
	for _, test := range db.event.Testtuples {
		assert.Equalf(t, StatusTodo, test.Status, "blame it on %+v", test)
	}
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
				Worker:  worker,
			},
			{
				Key:         computePlanAggregatetupleKey2,
				AlgoKey:     aggregateAlgoKey,
				ID:          IDs[3],
				InModelsIDs: []string{IDs[2]},
				Worker:      worker,
			},
		},
	}

	outCP, err := createComputePlanInternal(db, inCP, tag, map[string]string{}, false)
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
	cp, err := queryComputePlan(db, assetToArgs(inputKey{Key: outCP.Key}))
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
	outCP, err := createComputePlanInternal(db, inCP, tag, map[string]string{}, false)
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

	out, err := createComputePlanInternal(db, modelCompositionComputePlan, tag, map[string]string{}, false)
	assert.NoError(t, err)
	assert.NotNil(t, db.event)
	assert.Len(t, db.event.CompositeTraintuples, 2)

	_, err = cancelComputePlan(db, assetToArgs(inputKey{Key: out.Key}))
	assert.NoError(t, err)

	computePlan, err := getOutComputePlan(db, out.Key)
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

	out, err := createComputePlanInternal(db, modelCompositionComputePlan, tag, map[string]string{}, false)
	assert.NoError(t, err)

	logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[0]}))
	logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[1]}))
	logFailCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[1]}))

	_, err = cancelComputePlan(db, assetToArgs(inputKey{Key: out.Key}))
	assert.NoError(t, err)

	computePlan, err := getOutComputePlan(db, out.Key)
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

	out, err := createComputePlanInternal(db, modelCompositionComputePlan, tag, map[string]string{}, false)
	assert.NoError(t, err)

	logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[0]}))
	logStartCompositeTrain(db, assetToArgs(inputKey{out.CompositeTraintupleKeys[1]}))

	_, err = cancelComputePlan(db, assetToArgs(inputKey{Key: out.Key}))
	assert.NoError(t, err)

	inp := inputLogSuccessCompositeTrain{}
	inp.fillDefaults()
	inp.Key = out.CompositeTraintupleKeys[1]
	_, err = logSuccessCompositeTrain(db, assetToArgs(inp))
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

func TestCleanModels(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	mockStub.MockTransactionStart("42")
	registerItem(t, *mockStub, "aggregateAlgo")
	db := NewLedgerDB(mockStub)

	out, err := createComputePlanInternal(db, defaultComputePlan, tag, map[string]string{}, true)
	assert.NoError(t, err)
	// Just created the compute plan so not in the event
	assert.Len(t, db.event.ComputePlans, 0)
	clearEvent(db)

	traintupleToDone(t, db, out.TraintupleKeys[0])
	// Present in the event but without any model to remove
	assert.Len(t, db.event.ComputePlans, 1)
	assert.Equal(t, db.event.ComputePlans[0].Status, StatusDoing)
	assert.Len(t, db.event.ComputePlans[0].ModelsToDelete, 0)
	clearEvent(db)

	traintupleToDone(t, db, out.TraintupleKeys[1])
	// Present in the event but with one intermediary model done to remove
	assert.Len(t, db.event.ComputePlans, 1)
	assert.Equal(t, db.event.ComputePlans[0].Status, StatusDoing)
	assert.Len(t, db.event.ComputePlans[0].ModelsToDelete, 1)
	clearEvent(db)

	testtupleToDone(t, db, out.TesttupleKeys[0])
	// Present in the event but without any model to remove because the last
	// model is not intermerdiary
	assert.Len(t, db.event.ComputePlans, 1)
	assert.Equal(t, db.event.ComputePlans[0].Status, StatusDone)
	assert.Len(t, db.event.ComputePlans[0].ModelsToDelete, 0)
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
				Worker:  worker,
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
