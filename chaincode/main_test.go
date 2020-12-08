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

	peer "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"
)

const objectiveKey = "5c1d9cd1-c2c1-082d-de09-21b56d11030c"
const objectiveDescriptionChecksum = "5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const objectiveMetricsChecksum = "4a1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379"
const objectiveMetricsStorageAddress = "https://toto/objective/222/metrics"
const dataManagerKey = "da1bb7c3-1f62-244c-0f3a-761cc1688042"
const dataManagerOpenerChecksum = "da1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const trainDataSampleKey1 = "aa1bb7c3-1f62-244c-0f3a-761cc1688042"
const trainDataSampleKey2 = "aa2bb7c3-1f62-244c-0f3a-761cc1688042"
const testDataSampleKey1 = "bb1bb7c3-1f62-244c-0f3a-761cc1688042"
const testDataSampleKey2 = "bb2bb7c3-1f62-244c-0f3a-761cc1688042"
const algoKey = "fd1bb7c3-1f62-244c-0f3a-761cc1688042"
const algoKey2 = "cccbb7c3-1f62-244c-0f3a-761cc1688042"
const algoChecksum = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcc"
const algoStorageAddress = "https://toto/algo/222/algo"
const algoName = "hog + svm"
const compositeAlgoKey = "cccbb7c3-1f62-244c-0f3a-761cc1688042"
const compositeAlgoChecksum = "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482dcd"
const compositeAlgoStorageAddress = "https://toto/compositeAlgo/222/algo"
const compositeAlgoName = "hog + svm composite"
const aggregateAlgoKey = "dddbb7c3-1f62-244c-0f3a-761cc1688042"
const aggregateAlgoChecksum = "dddbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482ddd"
const aggregateAlgoStorageAddress = "https://toto/aggregateAlgo/222/algo"
const aggregateAlgoName = "hog + svm aggregate"
const modelKey = "eedbb7c3-1f62-244c-0f3a-761cc1688042"
const modelChecksum = "eedbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482eed"
const modelAddress = "https://substrabac/model/toto"
const headModelKey = modelKey
const headModelChecksum = modelChecksum
const trunkModelKey = "ccdbb7c3-1f62-244c-0f3a-761cc1688042"
const trunkModelChecksum = "ccdbb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482ecc"
const trunkModelAddress = "https://substrabac/model/titi"
const traintupleKey = "b0289ab8-3a71-f01e-2b72-0259a6452244"
const traintupleKey2 = "bbb89ab8-3a71-f01e-2b72-0259a6452244"
const compositeTraintupleKey = "0c0d3956-26b4-878e-76d7-ba8bb6fb152e"
const aggregatetupleKey = "71527661-50f6-26d3-fa86-1bf6387e3896"
const testtupleKey = "dadada11-50f6-26d3-fa86-1bf6387e3896"
const testtupleKey2 = "bbbada11-50f6-26d3-fa86-1bf6387e3896"
const testtupleKey3 = "cccada11-50f6-26d3-fa86-1bf6387e3896"
const tag = "a tag is simply a string"
const computePlanKey = "00000000-50f6-26d3-fa86-1bf6387e3896"
const computePlanKey2 = "11111111-50f6-26d3-fa86-1bf6387e3896"
const computePlanTraintupleKey1 = "11000000-50f6-26d3-fa86-1bf6387e3896"
const computePlanTraintupleKey2 = "22000000-50f6-26d3-fa86-1bf6387e3896"
const computePlanTraintupleKey3 = "33000000-50f6-26d3-fa86-1bf6387e3896"
const computePlanCompositeTraintupleKey1 = "11000011-50f6-26d3-fa86-1bf6387e3896"
const computePlanCompositeTraintupleKey2 = "22000011-50f6-26d3-fa86-1bf6387e3896"
const computePlanCompositeTraintupleKey3 = "33000011-50f6-26d3-fa86-1bf6387e3896"
const computePlanCompositeTraintupleKey4 = "44000011-50f6-26d3-fa86-1bf6387e3896"
const computePlanAggregatetupleKey1 = "11000022-50f6-26d3-fa86-1bf6387e3896"
const computePlanAggregatetupleKey2 = "22000022-50f6-26d3-fa86-1bf6387e3896"
const computePlanTesttupleKey1 = "11000033-50f6-26d3-fa86-1bf6387e3896"
const computePlanTesttupleKey2 = "22000033-50f6-26d3-fa86-1bf6387e3896"
const computePlanTesttupleKey3 = "33000033-50f6-26d3-fa86-1bf6387e3896"
const computePlanTesttupleKey4 = "44000033-50f6-26d3-fa86-1bf6387e3896"
const computePlanTesttupleKey5 = "55000033-50f6-26d3-fa86-1bf6387e3896"

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
	return [][]byte{[]byte("queryAlgo"), keyToJSON(key)}
}

func assetToArgs(asset interface{}) []string {
	return []string{string(assetToJSON(asset))}
}

func keyToArgs(key string) []string {
	return []string{string(keyToJSON(key))}
}

func assetToJSON(asset interface{}) []byte {
	assetjson, _ := json.Marshal(asset)
	return assetjson
}

func keyToJSON(key string) []byte {
	return assetToJSON(inputKey{Key: key})
}

func registerItem(t *testing.T, mockStub MockStub, itemType string) (peer.Response, interface{}) {

	// 1. add dataManager
	inpDataManager := inputDataManager{}
	args := inpDataManager.createDefault()
	resp := mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding dataManager with status %d and message %s", resp.Status, resp.Message)
	if itemType == "dataManager" {
		return resp, inpDataManager
	}
	// 2. add test dataSample
	inpDataSample := inputDataSample{
		Keys:            []string{testDataSampleKey1, testDataSampleKey2},
		DataManagerKeys: []string{dataManagerKey},
		TestOnly:        "true",
	}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding test dataSample with status %d and message %s", resp.Status, resp.Message)
	if itemType == "testDataset" {
		return resp, inpDataSample
	}
	// 3. add objective
	inpObjective := inputObjective{}
	args = inpObjective.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding objective with status %d and message %s", resp.Status, resp.Message)
	if itemType == "objective" {
		return resp, inpObjective
	}
	// 4. Add train dataSample
	inpDataSample = inputDataSample{}
	args = inpDataSample.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding train dataSample with status %d and message %s", resp.Status, resp.Message)
	if itemType == "trainDataset" {
		return resp, inpDataSample
	}
	// 5. Add algo
	inpAlgo := inputAlgo{}
	args = inpAlgo.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "algo" {
		return resp, inpAlgo
	}
	// 6. Add composite algo
	inpCompositeAlgo := inputCompositeAlgo{}
	args = inpCompositeAlgo.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding composite algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "compositeAlgo" {
		return resp, inpCompositeAlgo
	}
	// 7. Add aggregate algo
	inpAggregateAlgo := inputAggregateAlgo{}
	args = inpAggregateAlgo.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding aggregate algo with status %d and message %s", resp.Status, resp.Message)
	if itemType == "aggregateAlgo" {
		return resp, inpAggregateAlgo
	}
	// 8. Add traintuple
	inpTraintuple := inputTraintuple{}
	args = inpTraintuple.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding traintuple with status %d and message %s", resp.Status, resp.Message)

	if itemType == "traintuple" {
		return resp, inpTraintuple
	}
	// 9. Add composite traintuple
	inpCompositeTraintuple := inputCompositeTraintuple{}
	args = inpCompositeTraintuple.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding composite traintuple with status %d and message %s", resp.Status, resp.Message)
	if itemType == "compositeTraintuple" {
		return resp, inpCompositeTraintuple
	}
	// 10. Add aggregate tuple
	inpAggregatetuple := inputAggregatetuple{}
	args = inpAggregatetuple.createDefault()
	resp = mockStub.MockInvoke(args)
	require.EqualValuesf(t, 200, resp.Status, "when adding aggregate tuple with status %d and message %s", resp.Status, resp.Message)
	if itemType == "aggregatetuple" {
		return resp, inpAggregatetuple
	}

	return resp, inpAggregatetuple
}

func registerRandomCompositeAlgo(t *testing.T, mockStub *MockStub) (key string, err error) {
	key = RandomUUID()
	inpAlgo := inputCompositeAlgo{inputAlgo{Key: key}}
	args := inpAlgo.createDefault()
	resp := mockStub.MockInvoke(args)
	if resp.Status != 200 {
		err = fmt.Errorf("failed to register random algo: %s", resp.Message)
		return
	}
	return
}

func registerTraintuple(t *testing.T, mockStub *MockStub, assetType AssetType) (key string, err error) {

	// 1. Generate and register random algo
	// 2. Generate and register traintuple using that algo

	randomAlgoKey := RandomUUID()
	randomTraintupleKey := RandomUUID()

	switch assetType {
	case CompositeTraintupleType:
		randomAlgoKey, err = registerRandomCompositeAlgo(t, mockStub)
		if err != nil {
			return
		}
		inpTraintuple := inputCompositeTraintuple{Key: randomTraintupleKey, AlgoKey: randomAlgoKey}
		inpTraintuple.fillDefaults()
		args := inpTraintuple.getArgs()
		resp := mockStub.MockInvoke(args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register traintuple: %s", resp.Message)
			return
		}
		var _key struct{ Key string }
		json.Unmarshal(resp.Payload, &_key)
		return _key.Key, nil
	case TraintupleType:
		inpAlgo := inputAlgo{Key: randomAlgoKey}
		args := inpAlgo.createDefault()
		resp := mockStub.MockInvoke(args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register random algo: %s", resp.Message)
			return
		}
		inpTraintuple := inputTraintuple{Key: randomTraintupleKey, AlgoKey: randomAlgoKey}
		args = inpTraintuple.createDefault()
		resp = mockStub.MockInvoke(args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register traintuple: %s", resp.Message)
			return
		}
		var _key struct{ Key string }
		json.Unmarshal(resp.Payload, &_key)
		return _key.Key, nil
	case AggregatetupleType:
		inpAlgo := inputAggregateAlgo{inputAlgo{Key: randomAlgoKey}}
		args := inpAlgo.createDefault()
		resp := mockStub.MockInvoke(args)
		if resp.Status != 200 {
			err = fmt.Errorf("failed to register random algo: %s", resp.Message)
			return
		}
		inpTraintuple := inputAggregatetuple{Key: randomTraintupleKey, AlgoKey: randomAlgoKey}
		args = inpTraintuple.createDefault()
		resp = mockStub.MockInvoke(args)
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
	logger.SetLevel(logrus.PanicLevel)
	os.Exit(m.Run())
}

func initializeMockStateDB(t *testing.T, stub *MockStub) {
	stub.MockTransactionStart(mockTxID)
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
			resp := mockStub.MockInvoke(args)

			expectedBookmark := ""
			if contractName == "queryModels" {
				// special case for queryModels: we return a map instead of a string
				expectedBookmarkBytes, _ := json.Marshal(queryModelsBookmarks{})
				expectedBookmark = string(expectedBookmarkBytes)
			}

			expectedResult := map[string]interface{}{
				"results":  make([]string, 0),
				"bookmark": expectedBookmark}

			expectedPayload, _ := json.Marshal(expectedResult)
			assert.Equal(t, expectedPayload, resp.Payload, "payload is not an empty list")
		})
	}
}

func RandomUUID() string {
	uuid, err := GetNewUUID()
	if err != nil {
		panic("GetNewUUID failed")
	}
	return uuid
}
