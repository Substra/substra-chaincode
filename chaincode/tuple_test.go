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

func TestRecursiveLogFailed(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	mockStub.MockTransactionStart("42")
	registerItem(t, *mockStub, "traintuple")
	db := NewLedgerDB(mockStub)

	childtraintuple := inputTraintuple{Key: RandomUUID()}
	childtraintuple.createDefault()
	childtraintuple.InModels = []string{traintupleKey}
	childResp, err := createTraintuple(db, assetToArgs(childtraintuple))
	assert.NoError(t, err)

	grandChildtraintuple := inputTraintuple{Key: RandomUUID()}
	grandChildtraintuple.createDefault()
	grandChildtraintuple.InModels = []string{childResp.Key}
	grandChildresp, err := createTraintuple(db, assetToArgs(grandChildtraintuple))
	assert.NoError(t, err)

	grandChildtesttuple := inputTesttuple{
		Key:           RandomUUID(),
		TraintupleKey: traintupleKey,
		ObjectiveKey:  objectiveKey,
	}
	testResp, err := createTesttuple(db, assetToArgs(grandChildtesttuple))
	assert.NoError(t, err)

	_, err = logStartTrain(db, assetToArgs(inputKey{Key: traintupleKey}))
	assert.NoError(t, err)
	_, err = logFailTrain(db, assetToArgs(inputKey{Key: traintupleKey}))
	assert.NoError(t, err)

	train2, err := db.GetTraintuple(grandChildresp.Key)
	assert.NoError(t, err)
	assert.Equal(t, StatusFailed, train2.Status)

	test, err := db.GetTesttuple(testResp.Key)
	assert.NoError(t, err)
	assert.Equal(t, StatusFailed, test.Status)
}

// myMockStub is here to simulate the fact that in real condition you cannot read
// what you just write. It should be improved and more generally used.
type myMockStub struct {
	saveWhenWriting bool
	writtenState    map[string][]byte
	*MockStub
}

func (stub *myMockStub) PutState(key string, value []byte) error {
	if !stub.saveWhenWriting {
		if stub.writtenState == nil {
			stub.writtenState = make(map[string][]byte)
		}
		stub.writtenState[key] = value
		return nil
	}
	return stub.PutState(key, value)
}

func (stub *myMockStub) saveWrittenState(t *testing.T) {
	if stub.writtenState == nil {
		return
	}
	for k, v := range stub.writtenState {
		err := stub.MockStub.PutState(k, v)
		if err != nil {
			t.Fatalf("unable to `PutState` in saveWrittenState %s", err)
		}
	}
	stub.writtenState = nil
	return
}

func TestSpecifiqArgSeq(t *testing.T) {
	t.SkipNow()
	// This test is a POC and a example of a test base on the output of the log
	// parameters directly copied in a test. It can be realy usesul for debugging
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	argSeq := [][]string{
		// []string{"registerDataManager", "Titanic", "17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223", "http://owkin.substrabac:8000/data_manager/17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223/opener/", "csv", "48c89276972363250ea949c32809020e9d7fda786547a570bcaecedcc5092627", "http://owkin.substrabac:8000/data_manager/17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223/description/", "", "all"},
		[]string{"registerDataManager", "\"{\\\"Name\\\":\\\"Titanic\\\",\\\"OpenerChecksum\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"OpenerStorageAddress\\\":\\\"http://owkin.substrabac:8000/data_manager/17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223/opener/\\\",\\\"Type\\\":\\\"csv\\\",\\\"DescriptionChecksum\\\":\\\"48c89276972363250ea949c32809020e9d7fda786547a570bcaecedcc5092627\\\",\\\"DescriptionStorageAddress\\\":\\\"http://owkin.substrabac:8000/data_manager/17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223/description/\\\",\\\"ObjectiveKey\\\":\\\"\\\",\\\"Permissions\\\":\\\"all\\\"}\""},
		[]string{"registerDataSample", "\"{\\\"Keys\\\":\\\"47f9af29d34d737acfb0e37d93bfa650979292297ed263e8536ef3d13f70c83e,df94060511117dd25da1d2b1846f9be17340128233c8b24694d5e780d909b22c,50b7a4b4f2541674958fd09a061276862e1e2ea4dbdd0e1af06e70051804e33b,1befb03ceed3ab7ec9fa4bebe9b681bbc7725a402e03f9e64f9f1677cf619183\\\",\\\"DataManagerKeys\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"TestOnly\\\":\\\"false\\\"}\""},
		[]string{"registerDataSample", "\"{\\\"Keys\\\":\\\"1a8532bd84d5ef785a4abe503a12bc7040c666a9f6264f982aa4ad77ff7217a8\\\",\\\"DataManagerKeys\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"TestOnly\\\":\\\"true\\\"}\""},
		[]string{"registerObjective", "\"{\\\"Name\\\":\\\"Titanic: Machine Learning From Disaster\\\",\\\"DescriptionChecksum\\\":\\\"1158d2f5c0cf9f80155704ca0faa28823b145b42ebdba2ca38bd726a1377e1cb\\\",\\\"DescriptionStorageAddress\\\":\\\"http://owkin.substrabac:8000/objective/1158d2f5c0cf9f80155704ca0faa28823b145b42ebdba2ca38bd726a1377e1cb/description/\\\",\\\"MetricsName\\\":\\\"accuracy\\\",\\\"MetricsChecksum\\\":\\\"0bc13ad2e481c1a52959a228984bbee2e31271d567ea55a458e9ae92d481fedb\\\",\\\"MetricsStorageAddress\\\":\\\"http://owkin.substrabac:8000/objective/1158d2f5c0cf9f80155704ca0faa28823b145b42ebdba2ca38bd726a1377e1cb/metrics/\\\",\\\"TestDataset\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223:1a8532bd84d5ef785a4abe503a12bc7040c666a9f6264f982aa4ad77ff7217a8\\\",\\\"Permissions\\\":\\\"all\\\"}\""},
		[]string{"registerAlgo", "\"{\\\"Name\\\":\\\"Constant death predictor\\\",\\\"Checksum\\\":\\\"10a16f1b96beb3c07550103a9f15b3c2a77b15046cc7c70b762606590fb99de9\\\",\\\"StorageAddress\\\":\\\"http://owkin.substrabac:8000/algo/10a16f1b96beb3c07550103a9f15b3c2a77b15046cc7c70b762606590fb99de9/file/\\\",\\\"DescriptionChecksum\\\":\\\"1dae14e339c94ae04cc8846d353c07c8de96a38d6c5b5ee4486c4102ff011450\\\",\\\"DescriptionStorageAddress\\\":\\\"http://owkin.substrabac:8000/algo/10a16f1b96beb3c07550103a9f15b3c2a77b15046cc7c70b762606590fb99de9/description/\\\",\\\"Permissions\\\":\\\"all\\\"}\""},
		[]string{"createTraintuple", "\"{\\\"AlgoKey\\\":\\\"10a16f1b96beb3c07550103a9f15b3c2a77b15046cc7c70b762606590fb99de9\\\",\\\"ObjectiveKey\\\":\\\"1158d2f5c0cf9f80155704ca0faa28823b145b42ebdba2ca38bd726a1377e1cb\\\",\\\"InModels\\\":\\\"\\\",\\\"DataManagerKey\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"DataSampleKeys\\\":\\\"47f9af29d34d737acfb0e37d93bfa650979292297ed263e8536ef3d13f70c83e,df94060511117dd25da1d2b1846f9be17340128233c8b24694d5e780d909b22c,50b7a4b4f2541674958fd09a061276862e1e2ea4dbdd0e1af06e70051804e33b\\\",\\\"FLTask\\\":\\\"\\\",\\\"Rank\\\":\\\"\\\",\\\"Tag\\\":\\\"titanic v0\\\"}\""},
		[]string{"createTesttuple", "\"{\\\"TraintupleKey\\\":\\\"8daf7d448d0318dd8b06648cf32dde35f36171b308dec8675c8ff8e718acdac4\\\",\\\"DataManagerKey\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"DataSampleKeys\\\":\\\"1befb03ceed3ab7ec9fa4bebe9b681bbc7725a402e03f9e64f9f1677cf619183\\\",\\\"Tag\\\":\\\"titanic v0\\\"}\""},
		[]string{"createTesttuple", "\"{\\\"TraintupleKey\\\":\\\"8daf7d448d0318dd8b06648cf32dde35f36171b308dec8675c8ff8e718acdac4\\\",\\\"DataManagerKey\\\":\\\"\\\",\\\"DataSampleKeys\\\":\\\"\\\",\\\"Tag\\\":\\\"\\\"}\""},
		[]string{"logStartTrain", "\"{\\\"Key\\\":\\\"8daf7d448d0318dd8b06648cf32dde35f36171b308dec8675c8ff8e718acdac4\\\"}\""},
		[]string{"logSuccessTrain", "\"{\\\"Key\\\":\\\"8daf7d448d0318dd8b06648cf32dde35f36171b308dec8675c8ff8e718acdac4\\\",\\\"Log\\\":\\\"Train - CPU:119.66 % - Mem:0.04 GB - GPU:0.00 % - GPU Mem:0.00 GB; \\\",\\\"OutModel\\\":{\\\"Checksum\\\":\\\"6f6f2c318ff95ea7de9e4c01395b78b9217ddb134279275dae7842e7d4eb4c16\\\",\\\"StorageAddress\\\":\\\"http://owkin.substrabac:8000/model/6f6f2c318ff95ea7de9e4c01395b78b9217ddb134279275dae7842e7d4eb4c16/file/\\\"},\\\"Perf\\\":0.61610484}\""},
		[]string{"logStartTest", "\"{\\\"Key\\\":\\\"81bad50d76898ba6ea5af9d0a4816726bd46b947730a1bc2dd1d6755e8ab682b\\\"}\""},
		[]string{"logSuccessTest", "\"{\\\"Key\\\":\\\"81bad50d76898ba6ea5af9d0a4816726bd46b947730a1bc2dd1d6755e8ab682b\\\",\\\"Log\\\":\\\"Test - CPU:0.00 % - Mem:0.00 GB - GPU:0.00 % - GPU Mem:0.00 GB; \\\",\\\"Perf\\\":0.6179775}\""},
	}
	for _, argList := range argSeq {
		args := [][]byte{}
		for _, arg := range argList {
			args = append(args, []byte(arg))
		}
		resp := mockStub.MockInvoke("42", args)
		assert.EqualValues(t, 200, resp.Status, resp.Message, argList[0])
	}
}

func TestTagTuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	registerItem(t, *mockStub, "algo")

	noTag := "This is not a tag because it's waaaaaaaaaaaaaaaayyyyyyyyyyyyyyyyyyyyyyy too long."

	inpTraintuple := inputTraintuple{Tag: noTag}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message)

	tag := "This is a tag"

	inpTraintuple = inputTraintuple{Tag: tag}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message)

	args = [][]byte{[]byte("queryTraintuples")}
	resp = mockStub.MockInvoke("42", args)

	traintuples := []outputTraintuple{}
	err := json.Unmarshal(resp.Payload, &traintuples)

	assert.NoError(t, err, "should be unmarshaled")
	assert.Len(t, traintuples, 1, "there should be one traintuple")
	assert.EqualValues(t, tag, traintuples[0].Tag)

	inpTesttuple := inputTesttuple{Tag: tag}
	args = inpTesttuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message)

	args = [][]byte{[]byte("queryTesttuples")}
	resp = mockStub.MockInvoke("42", args)
	testtuples := []outputTesttuple{}
	err = json.Unmarshal(resp.Payload, &testtuples)
	assert.NoError(t, err, "should be unmarshaled")
	assert.Len(t, testtuples, 1, "there should be one traintuple")
	assert.EqualValues(t, tag, testtuples[0].Tag)

	filter := inputQueryFilter{
		IndexName:  "testtuple~tag",
		Attributes: tag,
	}
	args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message)
	filtertuples := []outputTesttuple{}
	err = json.Unmarshal(resp.Payload, &filtertuples)
	assert.NoError(t, err, "should be unmarshaled")
	assert.Len(t, testtuples, 1, "there should be one traintuple")
	assert.EqualValues(t, tag, testtuples[0].Tag)

}

func TestQueryModelPermissions(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "traintuple")
	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)
	traintupleToDone(t, db, traintupleKey)
	outTrain, err := queryTraintuple(db, keyToArgs(traintupleKey))
	assert.NoError(t, err)
	outPerm, err := queryModelPermissions(db, keyToArgs(outTrain.OutModel.Key))
	assert.NoError(t, err)
	assert.NotZero(t, outPerm)
}

func TestQueryHeadModelPermissions(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	registerItem(t, *mockStub, "compositeTraintuple")
	mockStub.MockTransactionStart("42")
	db := NewLedgerDB(mockStub)

	_, err := logStartCompositeTrain(db, assetToArgs(inputKey{Key: compositeTraintupleKey}))
	assert.NoError(t, err)
	success := inputLogSuccessCompositeTrain{}
	success.Key = compositeTraintupleKey
	success.fillDefaults()
	_, err = logSuccessCompositeTrain(db, assetToArgs(success))
	assert.NoError(t, err)

	outTrain, err := queryCompositeTraintuple(db, keyToArgs(compositeTraintupleKey))
	assert.NoError(t, err)
	outPerm, err := queryModelPermissions(db, keyToArgs(outTrain.OutHeadModel.OutModel.Key))
	assert.NoError(t, err)
	assert.NotZero(t, outPerm)
	assert.False(t, outPerm.Process.Public)
	assert.Len(t, outPerm.Process.AuthorizedIDs, 1)
	assert.Contains(t, outPerm.Process.AuthorizedIDs, worker)
}
