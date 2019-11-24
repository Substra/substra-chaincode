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
)

var modelCompositionTests = []struct {
	head  AssetType
	trunk AssetType
	child AssetType
}{
	{head: CompositeTraintupleType, trunk: TraintupleType, child: CompositeTraintupleType},
	{head: CompositeTraintupleType, trunk: CompositeTraintupleType, child: CompositeTraintupleType},
	// TODO (aggregate)
}

// This test makes sures that a child traintuple's state is successfully updated when a parent's state is updated
// See assertions at the bottom of the test for more details.
// The test ensures the behavior is correct with various combinations of parent traintuple types (composite, regular).
func TestModelComposition(t *testing.T) {
	for _, tt := range modelCompositionTests {
		for _, status := range []string{"successHead", "successTrunk", "successBoth", "failHead", "failTrunk"} {
			testName := fmt.Sprintf("TestModelComposition_%s_%sHead_%sTrunk_%sChild", status, tt.head, tt.trunk, tt.child)
			t.Run(testName, func(t *testing.T) {
				scc := new(SubstraChaincode)
				mockStub := NewMockStubWithRegisterNode("substra", scc)
				registerItem(t, *mockStub, "compositeAlgo")

				// register head
				headKey, err := registerTraintuple(mockStub, tt.head)
				assert.NoError(t, err)

				// register trunk
				trunkKey, err := registerTraintuple(mockStub, tt.trunk)
				assert.NoError(t, err)

				mockStub.MockTransactionStart("42")
				db := NewLedgerDB(mockStub)

				// register child traintuple...
				child := inputCompositeTraintuple{}
				child.createDefault()
				child.InHeadModelKey = headKey
				child.InTrunkModelKey = trunkKey
				childResp, err := createCompositeTraintuple(db, assetToArgs(child))
				assert.NoError(t, err)
				childKey := childResp["key"]

				// ... and its testtuple
				childTesttuple := inputTesttuple{}
				childTesttuple.TraintupleKey = childKey
				childTesttuple.fillDefaults()
				childTestupleResp, err := createTesttuple(db, assetToArgs(childTesttuple))
				assert.NoError(t, err)
				childTesttupleKey := childTestupleResp["key"]

				// start parents
				_, err = logStartCompositeTrain(db, assetToArgs(inputHash{Key: headKey}))
				assert.NoError(t, err)
				switch tt.trunk {
				case TraintupleType:
					_, err = logStartTrain(db, assetToArgs(inputHash{Key: trunkKey}))
				case CompositeTraintupleType:
					_, err = logStartCompositeTrain(db, assetToArgs(inputHash{Key: trunkKey}))
				default:
					assert.NoError(t, fmt.Errorf("unsupported test case %s", tt.trunk))
				}
				assert.NoError(t, err)

				// succeed/fail parents
				switch status {
				case "successHead":
					fallthrough
				case "successTrunk":
					fallthrough
				case "successBoth":
					if status == "successBoth" || status == "successHead" {
						successHead := inputLogSuccessCompositeTrain{}
						successHead.Key = headKey
						successHead.fillDefaults()
						_, err = logSuccessCompositeTrain(db, assetToArgs(successHead))
						assert.NoError(t, err)
					}
					if status == "successBoth" || status == "successTrunk" {
						switch tt.trunk {
						case TraintupleType:
							successTrunk := inputLogSuccessTrain{}
							successTrunk.fillDefaults()
							successTrunk.Key = trunkKey
							_, err = logSuccessTrain(db, assetToArgs(successTrunk))
						case CompositeTraintupleType:
							successTrunk := inputLogSuccessCompositeTrain{}
							successTrunk.fillDefaults()
							successTrunk.Key = trunkKey
							_, err = logSuccessCompositeTrain(db, assetToArgs(successTrunk))
						default:
							assert.NoError(t, fmt.Errorf("unsupported test case %s", tt.trunk))
						}
						assert.NoError(t, err)
					}
				case "failHead":
					failHead := inputLogFailTrain{}
					failHead.Key = headKey
					failHead.fillDefaults()
					_, err = logFailCompositeTrain(db, assetToArgs(failHead))
					assert.NoError(t, err)
				case "failTrunk":
					failTrunk := inputLogFailTrain{}
					failTrunk.Key = trunkKey
					failTrunk.fillDefaults()
					switch tt.trunk {
					case TraintupleType:
						_, err = logFailTrain(db, assetToArgs(failTrunk))
						assert.NoError(t, err)
					case CompositeTraintupleType:
						_, err = logFailCompositeTrain(db, assetToArgs(failTrunk))
						assert.NoError(t, err)
					default:
						assert.NoError(t, fmt.Errorf("unsupported test case %s", tt.trunk))
					}
				default:
					assert.NoError(t, fmt.Errorf("unsupported test case %s", status))
				}

				// check state of child traintuple/testtuple
				outChild, err := db.GetCompositeTraintuple(childKey)
				assert.NoError(t, err)
				trainChildStatus := outChild.Status

				outChildTesttuple, err := db.GetTesttuple(childTesttupleKey)
				assert.NoError(t, err)
				testChildStatus := outChildTesttuple.Status

				switch status {
				case "successHead":
					assert.Equal(t, StatusWaiting, trainChildStatus, "Only one parent has succeeded. The child traintuple should be Waiting")
					assert.Equal(t, StatusWaiting, testChildStatus, "Only one parent has succeeded. The child testtuple should be Waiting")
				case "successTrunk":
					assert.Equal(t, StatusWaiting, trainChildStatus, "Only one parent has succeeded. The child traintuple should be Waiting")
					assert.Equal(t, StatusWaiting, testChildStatus, "Only one parent has succeeded. The child testtuple should be Waiting")
				case "successBoth":
					assert.Equal(t, StatusTodo, trainChildStatus, "Both parents have succeded. The child traintuple should be Todo")
				case "failHead":
					assert.Equal(t, StatusFailed, trainChildStatus, "One parent has failed. The child traintuple should be Failed")
					assert.Equal(t, StatusFailed, testChildStatus, "One parent has failed. The child testtuple should be Failed")
				case "failTrunk":
					assert.Equal(t, StatusFailed, trainChildStatus, "One parent has failed. The child traintuple should be Failed")
					assert.Equal(t, StatusFailed, testChildStatus, "One parent has failed. The child testtuple should be Failed")
				default:
					assert.NoError(t, fmt.Errorf("unsupported test case %s", status))
				}
			})
		}
	}
}
