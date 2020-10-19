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
	"io"
	"os"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const objectiveKey = "5c1d9cd1c2c1082dde0921b56d11030c"
const objectiveDescriptionHash = "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const objectiveDescriptionStorageAddress = "https://toto/objective/222/description"
const objectiveMetricsHash = "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const objectiveMetricsStorageAddress = "https://toto/objective/222/metrics"
const dataManagerOpenerHash = "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataSampleHash1 = "aa1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataSampleHash2 = "aa2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const testDataSampleHash1 = "bb1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const testDataSampleHash2 = "bb2bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const algoHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const algoStorageAddress = "https://toto/algo/222/algo"
const algoName = "hog + svm"
const compositeAlgoHash = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcd"
const compositeAlgoStorageAddress = "https://toto/compositeAlgo/222/algo"
const compositeAlgoName = "hog + svm composite"
const aggregateAlgoHash = "dddbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482ddd"
const aggregateAlgoStorageAddress = "https://toto/aggregateAlgo/222/algo"
const aggregateAlgoName = "hog + svm aggregate"
const modelHash = "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"
const modelAddress = "https://substrabac/model/toto"
const headModelHash = modelHash
const trunkModelHash = "ccdbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482ecc"
const trunkModelAddress = "https://substrabac/model/titi"
const worker = "SampleOrg"
const traintupleKey = "ebbf6cdde286539ea9cc34214dce7acb71e72799a676e4845be1b0fea155b35c"
const compositeTraintupleKey = "b6f20e7ed89073d995c50d4b2bff6bb365a05b8d77f10469117d8aad81d83989"
const aggregatetupleKey = "48c17bb556e1a122138d89178d81b22469a0cae260af322de9b391086ad27b2c"
const tag = "a tag is simply a string"

func TestInit(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)

	// resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	resp := mockStub.MockInit("42", [][]byte{[]byte("init")})
	assert.EqualValuesf(t, 200, resp.Status, "init failed with status %d and message %s", resp.Status, resp.Message)
}

func methodToByte(methodName string) [][]byte {
	return [][]byte{[]byte(methodName)}
}

func methodAndAssetToByte(methodName string, asset interface{}) [][]byte {
	return [][]byte{[]byte(methodName), assetToJSON(asset)}
}

func methodAndKeyToByte(key string, asset interface{}) [][]byte {
	return [][]byte{[]byte("queryAlgo"), keyToJSONOld(key)}
}

func assetToArgs(asset interface{}) []string {
	return []string{string(assetToJSON(asset))}
}

func keyToArgs(key string) []string {
	return []string{string(keyToJSONOld(key))}
}

func assetToJSON(asset interface{}) []byte {
	assetjson, _ := json.Marshal(asset)
	return assetjson
}

func keyToJSONOld(key string) []byte {
	return assetToJSON(inputKeyOld{Key: key})
}

func keyToJSON(key string) []byte {
	return assetToJSON(inputKey{Key: key})
}

func registerItem(t *testing.T, mockStub MockStub, itemType string) (peer.Response, interface{}) {
	// 1. add dataManager
	inpDataManager := inputDataManager{}
	args := inpDataManager.createDefault()
	resp := mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding dataManager with status %d and message %s", resp.Status, resp.Message)
	if itemType == "dataManager" {
		return resp, inpDataManager
	}
	// 2. add test dataSample
	inpDataSample := inputDataSample{
		Hashes:          []string{testDataSampleHash1, testDataSampleHash2},
		DataManagerKeys: []string{dataManagerOpenerHash},
		TestOnly:        "true",
	}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding test dataSample with status %d and message %s", resp.Status, resp.Message)
	if itemType == "testDataset" {
		return resp, inpDataSample
	}
	// 3. add objective
	inpObjective := inputObjective{}
	args = inpObjective.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding objective with status %d and message %s", resp.Status, resp.Message)
	if itemType == "objective" {
		return resp, inpObjective
	}
	// 4. Add train dataSample
	inpDataSample = inputDataSample{}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding train dataSample with status %d and message %s", resp.Status, resp.Message)
	if itemType == "trainDataset" {
		return resp, inpDataSample
	}
	// 5. Add algo
	inpAlgo := inputAlgo{}
	args = inpAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "algo" {
		return resp, inpAlgo
	}
	// 6. Add composite algo
	inpCompositeAlgo := inputCompositeAlgo{}
	args = inpCompositeAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding composite algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "compositeAlgo" {
		return resp, inpCompositeAlgo
	}
	// 7. Add aggregate algo
	inpAggregateAlgo := inputAggregateAlgo{}
	args = inpAggregateAlgo.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding aggregate algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "aggregateAlgo" {
		return resp, inpAggregateAlgo
	}
	// 8. Add traintuple
	inpTraintuple := inputTraintuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)
	if itemType == "traintuple" {
		return resp, inpTraintuple
	}
	// 9. Add composite traintuple
	inpCompositeTraintuple := inputCompositeTraintuple{}
	args = inpCompositeTraintuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding composite traintuple with status %d and message %s", resp.Status, resp.Message)
	if itemType == "compositeTraintuple" {
		return resp, inpCompositeTraintuple
	}
	// 10. Add aggregate tuple
	inpAggregatetuple := inputAggregatetuple{}
	args = inpAggregatetuple.createDefault()
	resp = mockStub.MockInvoke("42", args)
	require.EqualValuesf(t, 200, resp.Status, "when adding aggregate tuple with status %d and message %s", resp.Status, resp.Message)
	if itemType == "aggregatetuple" {
		return resp, inpAggregatetuple
	}

	return resp, inpAggregatetuple
}

func registerRandomCompositeAlgo(mockStub *MockStub) (key string, err error) {
	key = GetRandomHash()
	inpAlgo := inputCompositeAlgo{inputAlgo{Hash: key}}
	args := inpAlgo.createDefault()
	resp := mockStub.MockInvoke("42", args)
	if resp.Status != 200 {
		err = fmt.Errorf("failed to register random algo: %s", resp.Message)
		return
	}
	return
}

func registerTraintuple(mockStub *MockStub, assetType AssetType) (key string, err error) {

	// 1. Generate and register random algo
	// 2. Generate and register traintuple using that algo

	switch assetType {
	case CompositeTraintupleType:
		randAlgoKey, _err := registerRandomCompositeAlgo(mockStub)
		if _err != nil {
			err = _err
			return
		}
		inpTraintuple := inputCompositeTraintuple{AlgoKey: randAlgoKey}
		inpTraintuple.fillDefaults()
		args := inpTraintuple.getArgs()
		resp := mockStub.MockInvoke("42", args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register traintuple: %s", resp.Message)
			return
		}
		var _key struct{ Key string }
		json.Unmarshal(resp.Payload, &_key)
		return _key.Key, nil
	case TraintupleType:
		randAlgoKey := GetRandomHash()
		inpAlgo := inputAlgo{Hash: randAlgoKey}
		args := inpAlgo.createDefault()
		resp := mockStub.MockInvoke("42", args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register random algo: %s", resp.Message)
			return
		}
		inpTraintuple := inputTraintuple{AlgoKey: randAlgoKey}
		args = inpTraintuple.createDefault()
		resp = mockStub.MockInvoke("42", args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register traintuple: %s", resp.Message)
			return
		}
		var _key struct{ Key string }
		json.Unmarshal(resp.Payload, &_key)
		return _key.Key, nil
	case AggregatetupleType:
		randAlgoKey := GetRandomHash()
		inpAlgo := inputAggregateAlgo{inputAlgo{Hash: randAlgoKey}}
		args := inpAlgo.createDefault()
		resp := mockStub.MockInvoke("42", args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register random algo: %s", resp.Message)
			return
		}
		inpTraintuple := inputAggregatetuple{AlgoKey: randAlgoKey}
		args = inpTraintuple.createDefault()
		resp = mockStub.MockInvoke("42", args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register traintuple: %s", resp.Message)
			return
		}
		var _key struct{ Key string }
		json.Unmarshal(resp.Payload, &_key)
		return _key.Key, nil
	default:
		err = fmt.Errorf("invalid asset type: %v", assetType)
		return
	}
}

func printResp(buf io.Writer, payload []byte) {
	var toPrint []byte
	if strings.HasPrefix(string(payload), "{") {
		obj := map[string]interface{}{}
		json.Unmarshal(payload, &obj)
		toPrint, _ = json.MarshalIndent(obj, "", " ")
	} else if strings.HasPrefix(string(payload), "[") {
		obj := []map[string]interface{}{}
		json.Unmarshal(payload, &obj)
		toPrint, _ = json.MarshalIndent(obj, "", " ")
	} else {
		toPrint = payload
	}
	fmt.Fprintf(buf, "##### Command output:\n```json\n%s\n```\n", toPrint)
}

func TestMain(m *testing.M) {
	//Raise log level to silence it during tests
	logger.SetLevel(shim.LogCritical)
	os.Exit(m.Run())
}

func initializeMockStateDB(t *testing.T, stub *MockStub) {
	stub.MockTransactionStart("42")
	stub.PutState("key", []byte("value"))
}

func TestQueryEmptyResponse(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	initializeMockStateDB(t, mockStub)

	smartContracts := []string{
		"queryAlgos",
		"queryDataSamples",
		"queryObjectives",
		"queryDataManagers",
		"queryTraintuples",
		"queryTesttuples",
		"queryModels",
	}

	for _, contractName := range smartContracts {
		t.Run(contractName, func(t *testing.T) {
			args := [][]byte{[]byte(contractName)}
			resp := mockStub.MockInvoke("42", args)

			expectedPayload, _ := json.Marshal(make([]string, 0))
			assert.Equal(t, expectedPayload, resp.Payload, "payload is not an empty list")
		})
	}
}
