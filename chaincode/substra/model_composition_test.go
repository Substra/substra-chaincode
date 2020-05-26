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

package substra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test makes sures that a child traintuple's state is successfully updated when a parent's state is updated
// See assertions at the bottom of the test for more details.
// The test ensures the behavior is correct with various combinations of parent traintuple types (composite, regular).
func TestModelComposition(t *testing.T) {
	testTable := []struct {
		parent1 AssetType // head for composite
		parent2 AssetType // trunk for composite
		child   AssetType
	}{
		{parent1: CompositeTraintupleType, parent2: TraintupleType, child: CompositeTraintupleType},
		{parent1: CompositeTraintupleType, parent2: CompositeTraintupleType, child: CompositeTraintupleType},
		{parent1: CompositeTraintupleType, parent2: AggregatetupleType, child: CompositeTraintupleType},
		{parent1: CompositeTraintupleType, parent2: CompositeTraintupleType, child: AggregatetupleType},
		{parent1: TraintupleType, parent2: AggregatetupleType, child: AggregatetupleType},
	}
	for _, tt := range testTable {
		for _, status := range []string{
			"successParent1",
			"successParent2",
			"successBoth",
			"failParent1",
			"failParent2",
		} {
			testName := fmt.Sprintf("TestModelComposition_%s_%sParent1_%sParent2_%sChild", status, tt.parent1, tt.parent2, tt.child)
			t.Run(testName, func(t *testing.T) {
				scc := new(Chaincode)
				mockStub := NewMockStubWithRegisterNode("substra", scc)
				registerItem(t, *mockStub, "aggregateAlgo")

				// register parents
				parent1Key, err := registerTraintuple(mockStub, tt.parent1)
				assert.NoError(t, err)
				parent2Key, err := registerTraintuple(mockStub, tt.parent2)
				assert.NoError(t, err)

				mockStub.MockTransactionStart("42")
				db := NewLedgerDB(mockStub)

				// register child traintuple...
				childKey := ""
				switch tt.child {
				case CompositeTraintupleType:
					child := inputCompositeTraintuple{}
					child.createDefault()
					child.InHeadModelKey = parent1Key
					child.InTrunkModelKey = parent2Key
					childResp, err := createCompositeTraintuple(db, assetToArgs(child))
					assert.NoError(t, err)
					childKey = childResp.Key
				case AggregatetupleType:
					child := inputAggregatetuple{}
					child.createDefault()
					child.InModels = []string{parent1Key, parent2Key}
					childResp, err := createAggregatetuple(db, assetToArgs(child))
					assert.NoError(t, err)
					childKey = childResp.Key
				default:
					assert.NoError(t, fmt.Errorf("unsupported test case %s", tt.parent2))
				}

				// ... and its testtuple
				childTesttuple := inputTesttuple{}
				childTesttuple.TraintupleKey = childKey
				childTesttuple.fillDefaults()
				childTestupleResp, err := createTesttuple(db, assetToArgs(childTesttuple))
				assert.NoError(t, err)
				childTesttupleKey := childTestupleResp.Key

				// start parents
				_, err = trainStart(db, tt.parent1, parent1Key)
				assert.NoError(t, err)
				_, err = trainStart(db, tt.parent2, parent2Key)
				assert.NoError(t, err)

				// succeed/fail parents
				switch status {
				case "successParent1", "successParent2", "successBoth":
					if status == "successBoth" || status == "successParent1" {
						_, err = trainSuccess(db, tt.parent1, parent1Key)
						assert.NoError(t, err)
					}
					if status == "successBoth" || status == "successParent2" {
						_, err = trainSuccess(db, tt.parent2, parent2Key)
						assert.NoError(t, err)

					}
				case "failParent1":
					_, err = trainFail(db, tt.parent1, parent1Key)
					assert.NoError(t, err)
				case "failParent2":
					_, err = trainFail(db, tt.parent2, parent2Key)
					assert.NoError(t, err)

				default:
					assert.NoError(t, fmt.Errorf("unsupported test case %s", status))
				}

				// check state of child traintuple/testtuple
				tuple, err := db.GetGenericTuple(childKey)
				assert.NoError(t, err)
				trainChildStatus := tuple.Status

				outChildTesttuple, err := db.GetTesttuple(childTesttupleKey)
				assert.NoError(t, err)
				testChildStatus := outChildTesttuple.Status

				switch status {
				case "successParent1":
					assert.Equal(t, StatusWaiting, trainChildStatus, "Only one parent has succeeded. The child traintuple should be Waiting")
					assert.Equal(t, StatusWaiting, testChildStatus, "Only one parent has succeeded. The child testtuple should be Waiting")
				case "successParent2":
					assert.Equal(t, StatusWaiting, trainChildStatus, "Only one parent has succeeded. The child traintuple should be Waiting")
					assert.Equal(t, StatusWaiting, testChildStatus, "Only one parent has succeeded. The child testtuple should be Waiting")
				case "successBoth":
					assert.Equal(t, StatusTodo, trainChildStatus, "Both parents have succeded. The child traintuple should be Todo")
				case "failParent1":
					assert.Equal(t, StatusFailed, trainChildStatus, "One parent has failed. The child traintuple should be Failed")
					assert.Equal(t, StatusFailed, testChildStatus, "One parent has failed. The child testtuple should be Failed")
				case "failParent2":
					assert.Equal(t, StatusFailed, trainChildStatus, "One parent has failed. The child traintuple should be Failed")
					assert.Equal(t, StatusFailed, testChildStatus, "One parent has failed. The child testtuple should be Failed")
				default:
					assert.NoError(t, fmt.Errorf("unsupported test case %s", status))
				}
			})
		}
	}
}

func trainStart(db *LedgerDB, tupleType AssetType, tupleKey string) (interface{}, error) {
	switch tupleType {
	case TraintupleType:
		return logStartTrain(db, assetToArgs(inputKey{Key: tupleKey}))
	case CompositeTraintupleType:
		return logStartCompositeTrain(db, assetToArgs(inputKey{Key: tupleKey}))
	case AggregatetupleType:
		return logStartAggregate(db, assetToArgs(inputKey{Key: tupleKey}))
	default:
		return nil, fmt.Errorf("unsupported test case %s", tupleType)
	}
}

func trainSuccess(db *LedgerDB, tupleType AssetType, tupleKey string) (interface{}, error) {
	switch tupleType {
	case TraintupleType:
		successParent1 := inputLogSuccessTrain{}
		successParent1.fillDefaults()
		successParent1.Key = tupleKey
		return logSuccessTrain(db, assetToArgs(successParent1))
	case CompositeTraintupleType:
		successParent1 := inputLogSuccessCompositeTrain{}
		successParent1.fillDefaults()
		successParent1.Key = tupleKey
		return logSuccessCompositeTrain(db, assetToArgs(successParent1))
	case AggregatetupleType:
		successParent1 := inputLogSuccessTrain{}
		successParent1.fillDefaults()
		successParent1.Key = tupleKey
		return logSuccessAggregate(db, assetToArgs(successParent1))
	default:
		return nil, fmt.Errorf("unsupported test case %s", tupleType)
	}
}

func trainFail(db *LedgerDB, tupleType AssetType, tupleKey string) (interface{}, error) {
	in := inputLogFailTrain{}
	in.Key = tupleKey
	in.fillDefaults()
	switch tupleType {
	case TraintupleType:
		return logFailTrain(db, assetToArgs(in))
	case CompositeTraintupleType:
		return logFailCompositeTrain(db, assetToArgs(in))
	case AggregatetupleType:
		return logFailAggregate(db, assetToArgs(in))
	default:
		return nil, fmt.Errorf("unsupported test case %s", tupleType)
	}
}
