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
	"chaincode/errors"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Create a global logger for the chaincode. Its default level is Info
var logger = shim.NewLogger("substra-chaincode")

// Chaincode is a Receiver for Chaincode shim functions
type Chaincode struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
// TODO!!!!
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting nothing...")
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Log all input for potential debug later on.
	logger.Infof("Args received by the chaincode: %#v", stub.GetStringArgs())

	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	db := NewLedgerDB(stub)
	body, err := argsToBody(args)
	if err != nil {
		return formatErrorResponse(err)
	}
	return Invoke(db, fn, body)
}

func argsToBody(args []string) (string, error) {
	var body string
	var err error
	switch len(args) {
	case 0:
	case 1:
		body = args[0]
	default:
		err = errors.BadRequest("arguments should only contains 1 json string, received: %s", args)
	}
	return body, err
}

// Invoke is called per transaction on the chaincode.
func Invoke(db *LedgerDB, fn string, body string) peer.Response {
	start := time.Now()

	// Seed with a timestamp from the channel header so the chaincode's output
	// stay determinist for each transaction. It's necessary because endorsers
	// will compare their own output to the proposal.
	// TODO: Pass this timestamp through the db/context
	timestamp, err := db.cc.GetTxTimestamp()
	if err != nil {
		return formatErrorResponse(err)
	}
	seedTime := time.Unix(timestamp.GetSeconds(), int64(timestamp.GetNanos()))
	rand.Seed(seedTime.UnixNano())

	var result interface{}
	switch fn {
	case "createComputePlan":
		result, err = createComputePlan(db, body)
	case "createTesttuple":
		result, err = createTesttuple(db, body)
	case "createTraintuple":
		result, err = createTraintuple(db, body)
	case "createCompositeTraintuple":
		result, err = createCompositeTraintuple(db, body)
	case "createAggregatetuple":
		result, err = createAggregatetuple(db, body)
	case "cancelComputePlan":
		result, err = cancelComputePlan(db, body)
	case "logFailTest":
		result, err = logFailTest(db, body)
	case "logFailTrain":
		result, err = logFailTrain(db, body)
	case "logFailCompositeTrain":
		result, err = logFailCompositeTrain(db, body)
	case "logFailAggregate":
		result, err = logFailAggregate(db, body)
	case "logStartTest":
		result, err = logStartTest(db, body)
	case "logStartTrain":
		result, err = logStartTrain(db, body)
	case "logStartCompositeTrain":
		result, err = logStartCompositeTrain(db, body)
	case "logStartAggregate":
		result, err = logStartAggregate(db, body)
	case "logSuccessTest":
		result, err = logSuccessTest(db, body)
	case "logSuccessTrain":
		result, err = logSuccessTrain(db, body)
	case "logSuccessCompositeTrain":
		result, err = logSuccessCompositeTrain(db, body)
	case "logSuccessAggregate":
		result, err = logSuccessAggregate(db, body)
	case "queryAlgo":
		result, err = queryAlgo(db, body)
	case "queryAlgos":
		result, err = queryAlgos(db, body)
	case "queryCompositeAlgo":
		result, err = queryCompositeAlgo(db, body)
	case "queryCompositeAlgos":
		result, err = queryCompositeAlgos(db, body)
	case "queryAggregateAlgo":
		result, err = queryAggregateAlgo(db, body)
	case "queryAggregateAlgos":
		result, err = queryAggregateAlgos(db, body)
	case "queryDataManager":
		result, err = queryDataManager(db, body)
	case "queryDataManagers":
		result, err = queryDataManagers(db, body)
	case "queryDataSamples":
		result, err = queryDataSamples(db, body)
	case "queryDataset":
		result, err = queryDataset(db, body)
	case "queryFilter":
		result, err = queryFilter(db, body)
	case "queryModelDetails":
		result, err = queryModelDetails(db, body)
	case "queryModelPermissions":
		result, err = queryModelPermissions(db, body)
	case "queryModels":
		result, err = queryModels(db, body)
	case "queryObjective":
		result, err = queryObjective(db, body)
	case "queryObjectiveLeaderboard":
		result, err = queryObjectiveLeaderboard(db, body)
	case "queryObjectives":
		result, err = queryObjectives(db, body)
	case "queryTesttuple":
		result, err = queryTesttuple(db, body)
	case "queryTesttuples":
		result, err = queryTesttuples(db, body)
	case "queryTraintuple":
		result, err = queryTraintuple(db, body)
	case "queryCompositeTraintuple":
		result, err = queryCompositeTraintuple(db, body)
	case "queryAggregatetuple":
		result, err = queryAggregatetuple(db, body)
	case "queryTraintuples":
		result, err = queryTraintuples(db, body)
	case "queryCompositeTraintuples":
		result, err = queryCompositeTraintuples(db, body)
	case "queryAggregatetuples":
		result, err = queryAggregatetuples(db, body)
	case "queryComputePlan":
		result, err = queryComputePlan(db, body)
	case "queryComputePlans":
		result, err = queryComputePlans(db, body)
	case "registerAlgo":
		result, err = registerAlgo(db, body)
	case "registerCompositeAlgo":
		result, err = registerCompositeAlgo(db, body)
	case "registerAggregateAlgo":
		result, err = registerAggregateAlgo(db, body)
	case "registerDataManager":
		result, err = registerDataManager(db, body)
	case "registerDataSample":
		result, err = registerDataSample(db, body)
	case "registerObjective":
		result, err = registerObjective(db, body)
	case "updateComputePlan":
		result, err = updateComputePlan(db, body)
	case "updateDataManager":
		result, err = updateDataManager(db, body)
	case "updateDataSample":
		result, err = updateDataSample(db, body)
	case "registerNode":
		result, err = registerNode(db, body)
	case "queryNodes":
		result, err = queryNodes(db, body)
	default:
		err = errors.BadRequest("function \"%s\" not implemented", fn)
	}
	// Invoke duration
	duration := int(time.Since(start).Nanoseconds()) / 1e6
	logger.Infof("Response from chaincode (in %dms): %#v, error: %s", duration, result, err)
	// Return the result as success payload
	if err != nil {
		return formatErrorResponse(err)
	}
	// Send event if there is any. It's done in one batch since we can only send
	// one event per call
	err = db.SendEvent()
	if err != nil {
		return formatErrorResponse(errors.Internal("could not send event: %s", err.Error()))
	}
	// Marshal to json the smartcontract result
	resp, err := json.Marshal(result)
	if err != nil {
		return formatErrorResponse(errors.Internal("could not format response: %s", err.Error()))
	}
	return shim.Success(resp)
}

func formatErrorResponse(err error) peer.Response {
	e := errors.Wrap(err)
	status := e.HTTPStatusCode()

	errStruct := map[string]interface{}{
		"error": e.Error(),
		// Serialize status in the message until fabric-sdk-py allows subtrabac to
		// access the status
		"status": status,
	}
	for k, v := range e.GetContext() {
		errStruct[k] = v
	}

	payload, _ := json.Marshal(errStruct)
	return peer.Response{
		Message: string(payload),
		Payload: payload,
		Status:  int32(status),
	}
}
