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

	var result interface{}
	var err error
	switch fn {
	case "createComputePlan":
		result, err = createComputePlan(stub, args)
	case "createTesttuple":
		result, err = createTesttuple(stub, args)
	case "createTraintuple":
		result, err = createTraintuple(stub, args)
	case "logFailTest":
		result, err = logFailTest(stub, args)
	case "logFailTrain":
		result, err = logFailTrain(stub, args)
	case "logStartTest":
		result, err = logStartTest(stub, args)
	case "logStartTrain":
		result, err = logStartTrain(stub, args)
	case "logSuccessTest":
		result, err = logSuccessTest(stub, args)
	case "logSuccessTrain":
		result, err = logSuccessTrain(stub, args)
	case "queryAlgo":
		result, err = queryAlgo(stub, args)
	case "queryAlgos":
		result, err = queryAlgos(stub, args)
	case "queryDataManager":
		result, err = queryDataManager(stub, args)
	case "queryDataManagers":
		result, err = queryDataManagers(stub, args)
	case "queryDataSamples":
		result, err = queryDataSamples(stub, args)
	case "queryDataset":
		result, err = queryDataset(stub, args)
	case "queryFilter":
		result, err = queryFilter(stub, args)
	case "queryModelDetails":
		result, err = queryModelDetails(stub, args)
	case "queryModels":
		result, err = queryModels(stub, args)
	case "queryObjective":
		result, err = queryObjective(stub, args)
	case "queryObjectives":
		result, err = queryObjectives(stub, args)
	case "queryTesttuple":
		result, err = queryTesttuple(stub, args)
	case "queryTesttuples":
		result, err = queryTesttuples(stub, args)
	case "queryTraintuple":
		result, err = queryTraintuple(stub, args)
	case "queryTraintuples":
		result, err = queryTraintuples(stub, args)
	case "registerAlgo":
		result, err = registerAlgo(stub, args)
	case "registerDataManager":
		result, err = registerDataManager(stub, args)
	case "registerDataSample":
		result, err = registerDataSample(stub, args)
	case "registerObjective":
		result, err = registerObjective(stub, args)
	case "updateDataManager":
		result, err = updateDataManager(stub, args)
	case "updateDataSample":
		result, err = updateDataSample(stub, args)
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
