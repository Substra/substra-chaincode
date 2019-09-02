package main

import (
	"chaincode/errors"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SubstraChaincode is a Receiver for Chaincode shim functions
type SubstraChaincode struct {
}

// Create a global logger for the chaincode. Its default level is Info
var logger = shim.NewLogger("substra-chaincode")

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
// TODO!!!!
func (t *SubstraChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting nothing...")
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (t *SubstraChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Log all input for potential debug later on.
	logger.Infof("Args received by the chaincode: %#v", stub.GetStringArgs())

	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	db := NewLedgerDB(stub)

	var result interface{}
	var err error
	switch fn {
	case "createComputePlan":
		result, err = createComputePlan(db, args)
	case "createTesttuple":
		result, err = createTesttuple(db, args)
	case "createTraintuple":
		result, err = createTraintuple(db, args)
	case "getObjectiveLeaderboard":
		result, err = getObjectiveLeaderboard(db, args)
	case "logFailTest":
		result, err = logFailTest(db, args)
	case "logFailTrain":
		result, err = logFailTrain(db, args)
	case "logStartTest":
		result, err = logStartTest(db, args)
	case "logStartTrain":
		result, err = logStartTrain(db, args)
	case "logSuccessTest":
		result, err = logSuccessTest(db, args)
	case "logSuccessTrain":
		result, err = logSuccessTrain(db, args)
	case "queryAlgo":
		result, err = queryAlgo(db, args)
	case "queryAlgos":
		result, err = queryAlgos(db, args)
	case "queryDataManager":
		result, err = queryDataManager(db, args)
	case "queryDataManagers":
		result, err = queryDataManagers(db, args)
	case "queryDataSamples":
		result, err = queryDataSamples(db, args)
	case "queryDataset":
		result, err = queryDataset(db, args)
	case "queryFilter":
		result, err = queryFilter(db, args)
	case "queryModelDetails":
		result, err = queryModelDetails(db, args)
	case "queryModels":
		result, err = queryModels(db, args)
	case "queryObjective":
		result, err = queryObjective(db, args)
	case "queryObjectives":
		result, err = queryObjectives(db, args)
	case "queryTesttuple":
		result, err = queryTesttuple(db, args)
	case "queryTesttuples":
		result, err = queryTesttuples(db, args)
	case "queryTraintuple":
		result, err = queryTraintuple(db, args)
	case "queryTraintuples":
		result, err = queryTraintuples(db, args)
	case "registerAlgo":
		result, err = registerAlgo(db, args)
	case "registerDataManager":
		result, err = registerDataManager(db, args)
	case "registerDataSample":
		result, err = registerDataSample(db, args)
	case "registerObjective":
		result, err = registerObjective(db, args)
	case "updateDataManager":
		result, err = updateDataManager(db, args)
	case "updateDataSample":
		result, err = updateDataSample(db, args)
	default:
		err = fmt.Errorf("function not implemented")
	}
	logger.Infof("Response from chaincode: %#v, error: %s", result, err)
	// Return the result as success payload
	if err != nil {
		return formatErrorResponse(err)
	}
	// Marshal to json the smartcontract result
	resp, err := json.Marshal(result)
	if err != nil {
		return formatErrorResponse(fmt.Errorf("could not format response for unknown reason"))
	}

	return shim.Success(resp)
}

func formatErrorResponse(err error) peer.Response {
	e := errors.Wrap(err)
	status := e.HTTPStatusCode()

	errStruct := map[string]interface{}{"error": e.Error()}
	payload, _ := json.Marshal(errStruct)
	// For now we still return both payload and message.
	// We also need to serialize the status into the message until fabrik-sdk-py
	// allow substrabac to access the status
	errStruct["status"] = status
	message, _ := json.Marshal(errStruct)
	return peer.Response{
		Message: string(message),
		Payload: payload,
		Status:  int32(status),
	}
}

// main function starts up the chaincode in the container during instantiate
func main() {
	// TODO use the same level as the shim or an env variable
	logger.SetLevel(shim.LogDebug)
	if err := shim.Start(new(SubstraChaincode)); err != nil {
		fmt.Printf("Error starting SubstraChaincode chaincode: %s", err)
	}
}
