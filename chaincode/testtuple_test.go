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
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTesttupleOnFailedTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	resp, _ := registerItem(t, *mockStub, "traintuple")

	res := outputKey{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	traintupleKey := res.Key

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
		DataSampleKeys: []string{testDataSampleKey2, testDataSampleKey1},
		DataManagerKey: dataManagerKey}
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
	inpTesttuple2 := inputTesttuple{DataSampleKeys: []string{trainDataSampleKey1}}
	args = inpTesttuple2.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status)
	assert.Contains(t, resp.Message, "invalid input: dataManagerKey and dataSampleKey should be provided together")

	// Add an uncertified testtuple successfully
	inpTesttuple3 := inputTesttuple{
		Key:            RandomUUID(),
		DataSampleKeys: []string{trainDataSampleKey1, trainDataSampleKey2},
		DataManagerKey: dataManagerKey}
	args = inpTesttuple3.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	// Fail to add the same testtuple with the same key
	inpTesttuple4 := inputTesttuple{}
	args = inpTesttuple4.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 409, resp.Status)
	assert.Contains(t, resp.Message, "already exists")
}

func TestQueryTesttuple(t *testing.T) {
	testTable := []struct {
		traintupleKey              string
		expectedTypeString         string
		expectedAlgoName           string
		expectedAlgoChecksum       string
		expectedAlgoStorageAddress string
	}{
		{
			traintupleKey:              traintupleKey,
			expectedTypeString:         "traintuple",
			expectedAlgoName:           algoName,
			expectedAlgoChecksum:       algoChecksum,
			expectedAlgoStorageAddress: algoStorageAddress,
		},
		{
			traintupleKey:              compositeTraintupleKey,
			expectedTypeString:         "composite_traintuple",
			expectedAlgoName:           compositeAlgoName,
			expectedAlgoChecksum:       compositeAlgoChecksum,
			expectedAlgoStorageAddress: compositeAlgoStorageAddress,
		},
		{
			traintupleKey:              aggregatetupleKey,
			expectedTypeString:         "aggregatetuple",
			expectedAlgoName:           aggregateAlgoName,
			expectedAlgoChecksum:       aggregateAlgoChecksum,
			expectedAlgoStorageAddress: aggregateAlgoStorageAddress,
		},
	}
	for _, tt := range testTable {
		t.Run("TestQueryTesttuple"+tt.expectedTypeString, func(t *testing.T) {
			scc := new(SubstraChaincode)
			mockStub := NewMockStubWithRegisterNode("substra", scc)
			registerItem(t, *mockStub, "aggregatetuple")

			// create testtuple
			dataSampleKeys := []string{trainDataSampleKey1, trainDataSampleKey2}
			inpTesttuple := inputTesttuple{
				TraintupleKey:  tt.traintupleKey,
				DataManagerKey: dataManagerKey,
				DataSampleKeys: dataSampleKeys,
			}
			inpTesttuple.fillDefaults()
			resp := mockStub.MockInvoke("42", inpTesttuple.getArgs())
			res := map[string]string{}
			json.Unmarshal(resp.Payload, &res)
			testtupleKey := res["key"]

			// query testtuple
			args := [][]byte{[]byte("queryTesttuple"), keyToJSON(testtupleKey)}
			resp = mockStub.MockInvoke("42", args)
			respTesttuple := resp.Payload
			testtuple := outputTesttuple{}
			json.Unmarshal(respTesttuple, &testtuple)

			// assert
			assert.Equal(t, worker, testtuple.Creator)
			assert.Equal(t, worker, testtuple.Dataset.Worker)
			assert.Equal(t, inpTesttuple.TraintupleKey, testtuple.TraintupleKey)
			assert.Equal(t, tt.expectedTypeString, testtuple.TraintupleType)
			assert.Equal(t, tt.expectedAlgoName, testtuple.Algo.Name)
			assert.Equal(t, tt.expectedAlgoChecksum, testtuple.Algo.Checksum)
			assert.Equal(t, tt.expectedAlgoStorageAddress, testtuple.Algo.StorageAddress)
			assert.Equal(t, StatusWaiting, testtuple.Status)
			assert.Equal(t, objectiveKey, testtuple.Objective.Key)
			assert.Equal(t, objectiveMetricsChecksum, testtuple.Objective.Metrics.Checksum)
			assert.Equal(t, objectiveMetricsStorageAddress, testtuple.Objective.Metrics.StorageAddress)
			assert.Equal(t, "", testtuple.Log)
			assert.Equal(t, "", testtuple.Tag)
			assert.EqualValues(t, 0, testtuple.Dataset.Perf)
			assert.Equal(t, dataManagerKey, testtuple.Dataset.Key)
			assert.Equal(t, dataSampleKeys, testtuple.Dataset.DataSampleKeys)
			assert.Equal(t, dataManagerKey, testtuple.Dataset.Key)
			assert.Equal(t, dataManagerOpenerChecksum, testtuple.Dataset.OpenerChecksum)
			assert.False(t, testtuple.Certified)
		})
	}
}

func TestTesttupleOnCompositeTraintuple(t *testing.T) {
	for _, status := range []string{StatusDone, StatusFailed} {
		testName := fmt.Sprintf("TestTesttupleOnCompositeTraintuple_%s", status)
		t.Run(testName, func(t *testing.T) {
			scc := new(SubstraChaincode)
			mockStub := NewMockStubWithRegisterNode("substra", scc)

			registerItem(t, *mockStub, "compositeTraintuple")

			inp := inputTesttuple{
				Key:           RandomUUID(),
				TraintupleKey: compositeTraintupleKey,
			}
			// Create a testtuple before training
			args := inp.createDefault()
			resp := mockStub.MockInvoke("42", args)
			assert.EqualValues(t, http.StatusOK, resp.Status, resp.Message)
			values := map[string]string{}
			json.Unmarshal(resp.Payload, &values)
			testTupleKey := values["key"]

			// Start training
			mockStub.MockTransactionStart("42")
			db := NewLedgerDB(mockStub)
			_, err := logStartCompositeTrain(db, assetToArgs(inputKey{Key: compositeTraintupleKey}))
			assert.NoError(t, err)

			// Succeed/fail training
			expectedTesttupleStatus := ""
			switch status {
			case StatusDone:
				inLog := inputLogSuccessCompositeTrain{}
				inLog.fillDefaults()
				_, err = logSuccessCompositeTrain(db, assetToArgs(inLog))
				assert.NoError(t, err)
				expectedTesttupleStatus = StatusTodo
			case StatusFailed:
				inLog := inputLogFailTrain{}
				inLog.Key = compositeTraintupleKey
				inLog.fillDefaults()
				_, err = logFailCompositeTrain(db, assetToArgs(inLog))
				assert.NoError(t, err)
				expectedTesttupleStatus = StatusFailed
			default:
				assert.NoError(t, fmt.Errorf("Unknown status %s", status))
			}

			testTuple, err := queryTesttuple(db, assetToArgs(inputKey{Key: testTupleKey}))
			assert.NoError(t, err)
			assert.Equal(t, expectedTesttupleStatus, testTuple.Status)
			assert.Equal(t, compositeTraintupleKey, testTuple.TraintupleKey)

			// Create a new testtuple *after* the traintuple has been set to failed/succeeded
			inp.Key = RandomUUID()
			inp.DataManagerKey = dataManagerKey
			inp.DataSampleKeys = []string{trainDataSampleKey1}
			args = inp.createDefault()
			resp = mockStub.MockInvoke("42", args)

			switch status {
			case StatusDone:
				assert.EqualValues(t, http.StatusOK, resp.Status, resp.Message)
				values = map[string]string{}
				json.Unmarshal(resp.Payload, &values)
				testTupleKey = values["key"]
				testTuple, err := queryTesttuple(db, assetToArgs(inputKey{Key: testTupleKey}))
				assert.NoError(t, err)
				assert.Equal(t, StatusTodo, testTuple.Status)
			case StatusFailed:
				assert.EqualValues(t, 400, resp.Status, "status should show an error since the traintuple is failed")
				assert.Contains(t, resp.Message, "could not register this testtuple")
			default:
				assert.NoError(t, fmt.Errorf("Unknown status %s", status))
			}
		})
	}
}
