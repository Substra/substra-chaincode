package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SubstraChaincode is a Receiver for Chaincode shim functions
type SubstraChaincode struct {
}

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
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result []byte
	var err error
	switch fn {
	case "registerChallenge":
		result, err = registerChallenge(stub, args)
	case "registerDataset":
		result, err = registerDataset(stub, args)
	case "updateDataset":
		result, err = updateDataset(stub, args)
	case "registerData":
		result, err = registerData(stub, args)
	case "updateData":
		result, err = updateData(stub, args)
	case "registerAlgo":
		result, err = registerAlgo(stub, args)
	case "createTraintuple":
		result, err = createTraintuple(stub, args)
	case "createTesttuple":
		result, err = createTesttuple(stub, args)
	case "logStartTrain":
		result, err = logStartTrain(stub, args)
	case "logStartTest":
		result, err = logStartTest(stub, args)
	case "logSuccessTrain":
		result, err = logSuccessTrain(stub, args)
	case "logSuccessTest":
		result, err = logSuccessTest(stub, args)
	case "logFailTrain":
		result, err = logFailTrain(stub, args)
	case "logFailTest":
		result, err = logFailTest(stub, args)
	case "query":
		result, err = query(stub, args)
	case "queryTraintuple":
		result, err = queryTraintuple(stub, args)
	case "queryChallenges":
		result, err = queryAll(stub, args, "challenge")
	case "queryAlgos":
		result, err = queryAll(stub, args, "algo")
	case "queryTraintuples":
		result, err = queryTraintuples(stub, args)
	case "queryTesttuples":
		result, err = queryAll(stub, args, "testtuple")
	case "queryDatasets":
		result, err = queryAll(stub, args, "dataset")
	case "queryFilter":
		result, err = queryFilter(stub, args)
	case "queryModelDetails":
		result, err = queryModelDetails(stub, args)
	case "queryModels":
		result, err = queryModels(stub, args)
	case "queryDatasetData":
		result, err = queryDatasetData(stub, args)
	default:
		err = fmt.Errorf("function not implemented")
	}
	// Return the result as success payload
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SubstraChaincode)); err != nil {
		fmt.Printf("Error starting SubstraChaincode chaincode: %s", err)
	}
}
