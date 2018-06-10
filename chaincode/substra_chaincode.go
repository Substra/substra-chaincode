package main

import (
	"fmt"
	"strings"

	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

// SubstraChaincode is a Receiver for Chaincode shim functions
type SubstraChaincode struct {
}

// Problem is one of the element type stored in the ledger
type Problem struct {
	name                      string   `json:"name"`
	descriptionStorageAddress string   `json:"descriptionStorageAddress"`
	metricsStorageAddress     string   `json:"metricsStorageAddress"`
	metricsHash               string   `json:"metricsHash"`
	owner                     string   `json:"owner"`
	testData                  []string `json:"testData"`
	permissions               string   `json:"permissions"`
}

// Data is one of the element type stored in the ledger
type Data struct {
	name        string   `json:"name"`
	dataOpener  string   `json:"dataOpener"`
	owner       string   `json:"owner"`
	problems    []string `json:problems`
	permissions string   `json:permissions`
}

// Algo is one of the element type stored in the ledger
type Algo struct {
	name           string `json:"name"`
	storageAddress string `json:"storageAddress"`
	owner          string `json:"owner"`
	problem        string `json:problem`
	permissions    string `json:permissions`
}

type TrainTuple struct {
	problem         map[string][2]string `json:"problem"`
	algo            map[string]string    `json:"algo"`
	startModel      map[string]string    `json:"startModel"`
	endModel        map[string]string    `json:"endModel"`
	trainData       []string             `json:"trainData"`
	trainDataOpener string               `json:"trainDataOpener"`
	trainWorker     string               `json:"trainWorker"`
	testWorker      string               `json:"testWorker"`
	status          string               `json:"status"`
	rank            int                  `json:"rank"`
	perf            float32              `json:"perf"`
	trainPerf       map[string]float32   `json:"trainPerf"`
	testPerf        map[string]float32   `json:"testPerf"`
	log             string               `json:"log"`
	permissions     string               `json:"permissions"`
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
// TODO!!!!
func (t *SubstraChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
// TODO
func (t *SubstraChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result []byte
	var err error
	switch fn {
	case "addData":
		result, err = addData(stub, args)
	case "addAlgo":
		result, err = addAlgo(stub, args)
	case "query":
		result, err = query(stub, args)
	case "queryData":
		result, err = queryData(stub, args)
	default:
		err = fmt.Errorf("function not implemented")
	}
	// Return the result as success payload
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

// addData stores a new data in the ledger.
// If the key exists, it will override the value with the new one
// TODO check if args 0 or args 1
func addData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("incorrect arguments, expecting 5 args: " +
			"data hash, name, data opener hash, associated problems, permissions")
	}

	// TODO check input types
	problems := strings.Split(strings.Replace(args[3], " ", "", -1), ",")
	// create data key
	key := "data_" + args[0]
	// find associated owner
	// TODO
	owner := "TODO"
	// create data object
	var data = Data{
		name:        args[1],
		dataOpener:  args[2],
		owner:       owner,
		problems:    problems,
		permissions: args[4]}
	dataBytes, _ := json.Marshal(data)
	err := stub.PutState(key, dataBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add data with hash %s", args[0])
	}
	// create composite keys (one for each associated problem) to find data associated with a problem
	indexName := "data~problem~key"
	for _, problem := range data.problems {
		compositeKey, err := stub.CreateCompositeKey(indexName, []string{"data", problem, key})
		if err != nil {
			return nil, err
		}
		value := []byte{0x00}
		err = stub.PutState(compositeKey, value)
		if err != nil {
			return nil, fmt.Errorf("failed to add composite key for data with hash %s", args[0])
		}
	}
	// create composite key to find data of a same dataset (= same data opener)
	return dataBytes, nil
}

// addAlgo stores a new algo in the ledger.
// If the key exists, it will override the value with the new one
// TODO check if args 0 or args 1
func addAlgo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("incorrect arguments, expecting 5 args: " +
			"algo hash, name, storage address, associated problem, permissions")
	}

	// TODO check input types
	problem := args[3]
	// create data key
	key := "algo_" + args[0]
	// find associated owner
	// TODO
	owner := "TODO"
	// create data object
	var algo = Algo{
		name:           args[1],
		storageAddress: args[2],
		owner:          owner,
		problem:        problem,
		permissions:    args[4]}
	algoBytes, _ := json.Marshal(algo)
	err := stub.PutState(key, algoBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add algo with hash %s", args[0])
	}
	// create composite key
	indexName := "algo~problem~key"
	compositeKey, err := stub.CreateCompositeKey(indexName, []string{"algo", problem, key})
	if err != nil {
		return nil, err
	}
	value := []byte{0x00}
	err = stub.PutState(compositeKey, value)
	if err != nil {
		return nil, fmt.Errorf("failed to add composite key for algo with hash %s", args[0])
	}
	return algoBytes, nil
}

// query returns an element of the ledger given its key
// For now, ok for everything. Later returns if the requester has permission to see it
func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key string

	if len(args) != 1 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting the key of the element to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	} else if valAsbytes == nil {
		return nil, fmt.Errorf("no element with this key %s", key)
	}

	return valAsbytes, nil
}

// query returns an element of the ledger given its key
// For now, ok for everything. Later returns if the requester has permission to see it
func queryData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}

	elementsIterator, err := stub.GetStateByRange("data_", "dataa")
	if err != nil {
		return nil, err
	}
	var elements []map[string]interface{}
	for elementsIterator.HasNext() {
		queryResponse, err := elementsIterator.Next()
		if err != nil {
			return nil, err
		}
		var element map[string]interface{}
		err = json.Unmarshal(queryResponse.GetValue(), &elements)
		if err != nil {
			return nil, err
		}
		element["key"] = queryResponse.GetKey()
		elements = append(elements, element)
	}
	payload, err := json.Marshal(elements)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SubstraChaincode)); err != nil {
		fmt.Printf("Error starting SubstraChaincode chaincode: %s", err)
	}
}
