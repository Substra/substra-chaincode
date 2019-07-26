package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// myMockStub is here to simulate the fact that in real condition you cannot read
// what you just write. It should be improved and more generally used.
type myMockStub struct {
	saveWhenWriting bool
	writenState     map[string][]byte
	*shim.MockStub
}

func (stub *myMockStub) PutState(key string, value []byte) error {
	if !stub.saveWhenWriting {
		if stub.writenState == nil {
			stub.writenState = make(map[string][]byte)
		}
		stub.writenState[key] = value
		return nil
	}
	return stub.PutState(key, value)
}

func (stub *myMockStub) saveWritenState(t *testing.T) {
	if stub.writenState == nil {
		return
	}
	for k, v := range stub.writenState {
		err := stub.MockStub.PutState(k, v)
		if err != nil {
			t.Fatalf("ennable to `PutStae` in saveWritenState %s", err)
		}
	}
	stub.writenState = nil
	return
}

func TestCreateComputePlan(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	myStub := myMockStub{MockStub: mockStub}
	myStub.saveWhenWriting = true
	registerItem(t, *mockStub, "algo")
	myStub.MockTransactionStart("42")
	myStub.saveWhenWriting = false

	// Simply test method and return values
	inCP := defaultComputePlan
	outCP, err := createComputePlan(&myStub, assetToArgs(inCP))
	assert.NoError(t, err)
	assert.NotNil(t, outCP)
	assert.Contains(t, outCP.TupleKeys, traintupleID1)
	assert.Contains(t, outCP.TupleKeys, traintupleID2)
	assert.Contains(t, outCP.TupleKeys, testtupleID)
	assert.EqualValues(t, outCP.TupleKeys[traintupleID1], outCP.FLTask)

	// Save all that was written in the mocked ledger
	myStub.saveWritenState(t)

	// Check the traintuples
	traintuples, err := queryTraintuples(&myStub, []string{})
	assert.NoError(t, err)
	assert.Len(t, traintuples, 2)
	var first, second outputTraintuple
	for _, el := range traintuples {
		switch el.Key {
		case outCP.TupleKeys[traintupleID1]:
			first = el
		case outCP.TupleKeys[traintupleID2]:
			second = el
		}
	}
	assert.NotZero(t, first)
	assert.NotZero(t, second)
	assert.EqualValues(t, first.Key, first.FLTask)
	assert.EqualValues(t, first.FLTask, second.FLTask)
	assert.Len(t, second.InModels, 1)
	assert.EqualValues(t, first.Key, second.InModels[0].TraintupleKey)
	assert.Equal(t, first.Status, StatusTodo)
	assert.Equal(t, second.Status, StatusWaiting)
}
func TestSpecifiqArgSeq(t *testing.T) {
	t.SkipNow()
	// This test is a POC and a example of a test base on the output of the log
	// parameters directly copied in a test. It can be realy usesul for debugging
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	argSeq := [][]string{
		// []string{"registerDataManager", "Titanic", "17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223", "http://owkin.substrabac:8000/data_manager/17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223/opener/", "csv", "48c89276972363250ea949c32809020e9d7fda786547a570bcaecedcc5092627", "http://owkin.substrabac:8000/data_manager/17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223/description/", "", "all"},
		[]string{"registerDataManager", "\"{\\\"Name\\\":\\\"Titanic\\\",\\\"OpenerHash\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"OpenerStorageAddress\\\":\\\"http://owkin.substrabac:8000/data_manager/17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223/opener/\\\",\\\"Type\\\":\\\"csv\\\",\\\"DescriptionHash\\\":\\\"48c89276972363250ea949c32809020e9d7fda786547a570bcaecedcc5092627\\\",\\\"DescriptionStorageAddress\\\":\\\"http://owkin.substrabac:8000/data_manager/17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223/description/\\\",\\\"ObjectiveKey\\\":\\\"\\\",\\\"Permissions\\\":\\\"all\\\"}\""},
		[]string{"registerDataSample", "\"{\\\"Hashes\\\":\\\"47f9af29d34d737acfb0e37d93bfa650979292297ed263e8536ef3d13f70c83e,df94060511117dd25da1d2b1846f9be17340128233c8b24694d5e780d909b22c,50b7a4b4f2541674958fd09a061276862e1e2ea4dbdd0e1af06e70051804e33b,1befb03ceed3ab7ec9fa4bebe9b681bbc7725a402e03f9e64f9f1677cf619183\\\",\\\"DataManagerKeys\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"TestOnly\\\":\\\"false\\\"}\""},
		[]string{"registerDataSample", "\"{\\\"Hashes\\\":\\\"1a8532bd84d5ef785a4abe503a12bc7040c666a9f6264f982aa4ad77ff7217a8\\\",\\\"DataManagerKeys\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"TestOnly\\\":\\\"true\\\"}\""},
		[]string{"registerObjective", "\"{\\\"Name\\\":\\\"Titanic: Machine Learning From Disaster\\\",\\\"DescriptionHash\\\":\\\"1158d2f5c0cf9f80155704ca0faa28823b145b42ebdba2ca38bd726a1377e1cb\\\",\\\"DescriptionStorageAddress\\\":\\\"http://owkin.substrabac:8000/objective/1158d2f5c0cf9f80155704ca0faa28823b145b42ebdba2ca38bd726a1377e1cb/description/\\\",\\\"MetricsName\\\":\\\"accuracy\\\",\\\"MetricsHash\\\":\\\"0bc13ad2e481c1a52959a228984bbee2e31271d567ea55a458e9ae92d481fedb\\\",\\\"MetricsStorageAddress\\\":\\\"http://owkin.substrabac:8000/objective/1158d2f5c0cf9f80155704ca0faa28823b145b42ebdba2ca38bd726a1377e1cb/metrics/\\\",\\\"TestDataset\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223:1a8532bd84d5ef785a4abe503a12bc7040c666a9f6264f982aa4ad77ff7217a8\\\",\\\"Permissions\\\":\\\"all\\\"}\""},
		[]string{"registerAlgo", "\"{\\\"Name\\\":\\\"Constant death predictor\\\",\\\"Hash\\\":\\\"10a16f1b96beb3c07550103a9f15b3c2a77b15046cc7c70b762606590fb99de9\\\",\\\"StorageAddress\\\":\\\"http://owkin.substrabac:8000/algo/10a16f1b96beb3c07550103a9f15b3c2a77b15046cc7c70b762606590fb99de9/file/\\\",\\\"DescriptionHash\\\":\\\"1dae14e339c94ae04cc8846d353c07c8de96a38d6c5b5ee4486c4102ff011450\\\",\\\"DescriptionStorageAddress\\\":\\\"http://owkin.substrabac:8000/algo/10a16f1b96beb3c07550103a9f15b3c2a77b15046cc7c70b762606590fb99de9/description/\\\",\\\"Permissions\\\":\\\"all\\\"}\""},
		[]string{"createTraintuple", "\"{\\\"AlgoKey\\\":\\\"10a16f1b96beb3c07550103a9f15b3c2a77b15046cc7c70b762606590fb99de9\\\",\\\"ObjectiveKey\\\":\\\"1158d2f5c0cf9f80155704ca0faa28823b145b42ebdba2ca38bd726a1377e1cb\\\",\\\"InModels\\\":\\\"\\\",\\\"DataManagerKey\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"DataSampleKeys\\\":\\\"47f9af29d34d737acfb0e37d93bfa650979292297ed263e8536ef3d13f70c83e,df94060511117dd25da1d2b1846f9be17340128233c8b24694d5e780d909b22c,50b7a4b4f2541674958fd09a061276862e1e2ea4dbdd0e1af06e70051804e33b\\\",\\\"FLTask\\\":\\\"\\\",\\\"Rank\\\":\\\"\\\",\\\"Tag\\\":\\\"titanic v0\\\"}\""},
		[]string{"createTesttuple", "\"{\\\"TraintupleKey\\\":\\\"8daf7d448d0318dd8b06648cf32dde35f36171b308dec8675c8ff8e718acdac4\\\",\\\"DataManagerKey\\\":\\\"17dbc4ece248304cab7b1dd53ec7edf1ebf8a5e12ff77a26dc6e8da9db4da223\\\",\\\"DataSampleKeys\\\":\\\"1befb03ceed3ab7ec9fa4bebe9b681bbc7725a402e03f9e64f9f1677cf619183\\\",\\\"Tag\\\":\\\"titanic v0\\\"}\""},
		[]string{"createTesttuple", "\"{\\\"TraintupleKey\\\":\\\"8daf7d448d0318dd8b06648cf32dde35f36171b308dec8675c8ff8e718acdac4\\\",\\\"DataManagerKey\\\":\\\"\\\",\\\"DataSampleKeys\\\":\\\"\\\",\\\"Tag\\\":\\\"\\\"}\""},
		[]string{"logStartTrain", "\"{\\\"Key\\\":\\\"8daf7d448d0318dd8b06648cf32dde35f36171b308dec8675c8ff8e718acdac4\\\"}\""},
		[]string{"logSuccessTrain", "\"{\\\"Key\\\":\\\"8daf7d448d0318dd8b06648cf32dde35f36171b308dec8675c8ff8e718acdac4\\\",\\\"Log\\\":\\\"Train - CPU:119.66 % - Mem:0.04 GB - GPU:0.00 % - GPU Mem:0.00 GB; \\\",\\\"OutModel\\\":{\\\"Hash\\\":\\\"6f6f2c318ff95ea7de9e4c01395b78b9217ddb134279275dae7842e7d4eb4c16\\\",\\\"StorageAddress\\\":\\\"http://owkin.substrabac:8000/model/6f6f2c318ff95ea7de9e4c01395b78b9217ddb134279275dae7842e7d4eb4c16/file/\\\"},\\\"Perf\\\":0.61610484}\""},
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

func TestTraintupleWithNoTestDataset(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	registerItem(t, *mockStub, "trainDataset")

	objHash := strings.ReplaceAll(objectiveDescriptionHash, "1", "2")
	inpObjective := inputObjective{DescriptionHash: objHash}
	inpObjective.createDefault()
	inpObjective.TestDataset = inputDataset{}
	resp := mockStub.MockInvoke("42", methodAndAssetToByte("registerObjective", inpObjective))
	assert.EqualValues(t, 200, resp.Status, "when adding objective without dataset it should work: ", resp.Message)

	inpAlgo := inputAlgo{}
	args := inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding algo it should work: ", resp.Message)

	inpTraintuple := inputTraintuple{ObjectiveKey: objHash}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "when adding traintuple without test dataset it should work: ", resp.Message)

	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, "It should find the traintuple without error ", resp.Message)
}
func TestTagTuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

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
func TestNoPanicWhileQueryingIncompleteTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "traintuple")

	// Manually open a ledger transaction
	mockStub.MockTransactionStart("42")
	defer mockStub.MockTransactionEnd("42")

	// Retreive and alter existing objectif to pass Metrics at nil
	objective := Objective{}
	err := getElementStruct(mockStub, objectiveDescriptionHash, &objective)
	assert.NoError(t, err)
	objective.Metrics = nil
	objBytes, err := json.Marshal(objective)
	assert.NoError(t, err)
	err = mockStub.PutState(objectiveDescriptionHash, objBytes)
	assert.NoError(t, err)
	// It should not panic
	require.NotPanics(t, func() {
		getOutputTraintuple(mockStub, traintupleKey)
	})
}
func TestTraintupleFLTaskCreation(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add dataManager, dataSample and algo
	registerItem(t, *mockStub, "algo")

	inpTraintuple := inputTraintuple{FLTask: "someFLTask"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for missing rank")
	require.Contains(t, resp.Message, "invalit inputs, a FLTask should have a rank", "invalid error message")

	inpTraintuple = inputTraintuple{Rank: "1"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 400, resp.Status, "should failed for invalid rank")
	require.Contains(t, resp.Message, "invalid inputs, a new FLTask should have a rank 0")

	inpTraintuple = inputTraintuple{Rank: "0"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	key := res["key"]
	require.EqualValues(t, key, traintupleKey)

	inpTraintuple = inputTraintuple{Rank: "0"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValues(t, 409, resp.Status, "should failed for existing FLTask")
	require.Contains(t, resp.Message, "this traintuple already exists")
}

func TestTraintupleMultipleFLTaskCreations(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add a some dataManager, dataSample and traintuple
	registerItem(t, *mockStub, "algo")

	inpTraintuple := inputTraintuple{Rank: "0"}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	key := res["key"]
	// Failed to add a traintuple with the same rank
	inpTraintuple = inputTraintuple{
		InModels: []string{key},
		Rank:     "0",
		FLTask:   key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add a traintuple of the same rank")

	// Failed to add a traintuple to an unexisting Fltask
	inpTraintuple = inputTraintuple{
		InModels: []string{key},
		Rank:     "1",
		FLTask:   "notarealone"}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "should failed to add a traintuple to an unexisting FLTask")

	// Succesfully add a traintuple to the same FLTask
	inpTraintuple = inputTraintuple{
		InModels: []string{key},
		Rank:     "1",
		FLTask:   key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status, resp.Message, "should be able do create a traintuple with the same FLTask")
	err = json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "should unmarshal without problem")
	assert.Contains(t, res, "key")
	ttkey := res["key"]
	// Add new algo to check all fltask algo consistency
	newAlgoHash := strings.Replace(algoHash, "a", "b", 1)
	inpAlgo := inputAlgo{Hash: newAlgoHash}
	args = inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 200, resp.Status)

	inpTraintuple = inputTraintuple{
		AlgoKey:  newAlgoHash,
		InModels: []string{ttkey},
		Rank:     "2",
		FLTask:   key}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValues(t, 400, resp.Status, resp.Message, "sould fail for it doesn't have the same algo key")
	assert.Contains(t, resp.Message, "does not have the same algo key")
}

func TestTesttupleOnFailedTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

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
	mockStub := shim.NewMockStub("substra", scc)

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
	mockStub := shim.NewMockStub("substra", scc)

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
	assert.Contains(t, resp.Message, "this testtuple already exists")
}

func TestTraintuple(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)

	// Add traintuple with invalid field
	inpTraintuple := inputTraintuple{
		AlgoKey: "aaa",
	}
	args := inpTraintuple.createDefault()
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding objective with invalid hash, status %d and message %s", resp.Status, resp.Message)

	// Add traintuple with unexisting algo
	inpTraintuple = inputTraintuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 400, resp.Status, "when adding traintuple with unexisting algo, status %d and message %s", resp.Status, resp.Message)

	// Properly add traintuple
	resp, tt := registerItem(t, *mockStub, "traintuple")

	inpTraintuple = tt.(inputTraintuple)
	res := map[string]string{}
	err := json.Unmarshal(resp.Payload, &res)
	assert.NoError(t, err, "traintuple should unmarshal without problem")
	assert.Contains(t, res, "key")
	traintupleKey := res["key"]
	// Query traintuple from key and check the consistency of returned arguments
	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)
	out := outputTraintuple{}
	err = json.Unmarshal(resp.Payload, &out)
	assert.NoError(t, err, "when unmarshalling queried traintuple")
	expected := outputTraintuple{
		Key: traintupleKey,
		Algo: &HashDressName{
			Hash:           algoHash,
			Name:           algoName,
			StorageAddress: algoStorageAddress,
		},
		Creator: "bbd157aa8e85eb985aeedb79361cd45739c92494dce44d351fd2dbd6190e27f0",
		Dataset: &TtDataset{
			DataSampleKeys: []string{trainDataSampleHash1, trainDataSampleHash2},
			OpenerHash:     dataManagerOpenerHash,
			Perf:           0.0,
			Worker:         worker,
		},
		Objective: &TtObjective{
			Key: objectiveDescriptionHash,
			Metrics: &HashDress{
				Hash:           objectiveMetricsHash,
				StorageAddress: objectiveMetricsStorageAddress,
			},
		},
		Permissions: "all",
		Status:      StatusTodo,
	}
	assert.Exactly(t, expected, out, "the traintuple queried from the ledger differ from expected")

	// Query all traintuples and check consistency
	args = [][]byte{[]byte("queryTraintuples")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuples - status %d and message %s", resp.Status, resp.Message)
	// TODO add traintuple key to output struct
	// For now we test it as cleanly as its added to the query response
	assert.Contains(t, string(resp.Payload), "key\":\""+traintupleKey)
	var queryTraintuples []outputTraintuple
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "traintuples should unmarshal without problem")
	assert.Exactly(t, out, queryTraintuples[0])

	// Add traintuple with inmodel from the above-submitted traintuple
	inpWaitingTraintuple := inputTraintuple{
		InModels: []string{string(traintupleKey)},
	}
	args = inpWaitingTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	//waitingTraintupleKey := string(resp.Payload)

	// Query traintuple with status todo and worker as trainworker and check consistency
	filter := inputQueryFilter{
		IndexName:  "traintuple~worker~status",
		Attributes: worker + ", todo",
	}
	args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple of worker with todo status - status %d and message %s", resp.Status, resp.Message)
	err = json.Unmarshal(resp.Payload, &queryTraintuples)
	assert.NoError(t, err, "traintuples should unmarshal without problem")
	assert.Exactly(t, out, queryTraintuples[0])

	// Update status and check consistency
	success := inputLogSuccessTrain{}
	success.Key = traintupleKey

	argsSlice := [][][]byte{
		[][]byte{[]byte("logStartTrain"), keyToJSON(traintupleKey)},
		success.createDefault(),
	}
	traintupleStatus := []string{StatusDoing, StatusDone}
	for i := range traintupleStatus {
		resp = mockStub.MockInvoke("42", argsSlice[i])
		require.EqualValuesf(t, 200, resp.Status, "when logging start %s with message %s", traintupleStatus[i], resp.Message)
		filter := inputQueryFilter{
			IndexName:  "traintuple~worker~status",
			Attributes: worker + ", " + traintupleStatus[i],
		}
		args = [][]byte{[]byte("queryFilter"), assetToJSON(filter)}
		resp = mockStub.MockInvoke("42", args)
		assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple of worker with %s status - message %s", traintupleStatus[i], resp.Message)
		sPayload := make([]map[string]interface{}, 1)
		assert.NoError(t, json.Unmarshal(resp.Payload, &sPayload), "when unmarshal queried traintuples")
		assert.EqualValues(t, traintupleKey, sPayload[0]["key"], "wrong retrieved key when querying traintuple of worker with %s status ", traintupleStatus[i])
		assert.EqualValues(t, traintupleStatus[i], sPayload[0]["status"], "wrong retrieved status when querying traintuple of worker with %s status ", traintupleStatus[i])
	}

	// Query Traintuple From key
	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying traintuple with status %d and message %s", resp.Status, resp.Message)
	endTraintuple := outputTraintuple{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &endTraintuple))
	expected.Dataset.Perf = success.Perf
	expected.Log = success.Log
	expected.OutModel = &HashDress{
		Hash:           modelHash,
		StorageAddress: modelAddress}
	expected.Status = traintupleStatus[1]
	assert.Exactly(t, expected, endTraintuple, "retreived Traintuple does not correspond to what is expected")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModelDetails"), keyToJSON(traintupleKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying model details with status %d and message %s", resp.Status, resp.Message)
	payload := outputModelDetails{}
	assert.NoError(t, json.Unmarshal(resp.Payload, &payload))
	assert.NotNil(t, payload.Traintuple, "when querying model tuples, payload should contain one traintuple")

	// query all traintuples related to a traintuple with the same algo
	args = [][]byte{[]byte("queryModels")}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying models with status %d and message %s", resp.Status, resp.Message)
}

func TestQueryTraintupleNotFound(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := shim.NewMockStub("substra", scc)
	registerItem(t, *mockStub, "traintuple")

	// queryTraintuple: normal case
	args := [][]byte{[]byte("queryTraintuple"), keyToJSON(traintupleKey)}
	resp := mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 200, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)

	// queryTraintuple: key does not exist
	notFoundKey := "eedbb7c31f62244c0f34461cc168804227115793d01c270021fe3f7935482eed"
	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(notFoundKey)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)

	// queryTraintuple: key does not exist and use existing other asset type key
	args = [][]byte{[]byte("queryTraintuple"), keyToJSON(algoHash)}
	resp = mockStub.MockInvoke("42", args)
	assert.EqualValuesf(t, 404, resp.Status, "when querying the traintuple - status %d and message %s", resp.Status, resp.Message)
}
