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
	Name                      string   `json:"name"`
	DescriptionStorageAddress string   `json:"descriptionStorageAddress"`
	MetricsStorageAddress     string   `json:"metricsStorageAddress"`
	MetricsHash               string   `json:"metricsHash"`
	Owner                     string   `json:"owner"`
	TestData                  []string `json:"testData"`
	Permissions               string   `json:"permissions"`
}

// Data is one of the element type stored in the ledger
type Data struct {
	Name        string   `json:"name"`
	DataOpener  string   `json:"dataOpener"`
	Owner       string   `json:"owner"`
	Problems    []string `json:problems`
	Permissions string   `json:permissions`
}

// Algo is one of the element type stored in the ledger
type Algo struct {
	Name           string `json:"name"`
	StorageAddress string `json:"storageAddress"`
	Owner          string `json:"owner"`
	Problem        string `json:problem`
	Permissions    string `json:permissions`
}

type TrainTuple struct {
	Problem         map[string][2]string `json:"problem"`
	Algo            map[string]string    `json:"algo"`
	StartModel      map[string]string    `json:"startModel"`
	EndModel        map[string]string    `json:"endModel"`
	TrainData       []string             `json:"trainData"`
	TrainDataOpener string               `json:"trainDataOpener"`
	TrainWorker     string               `json:"trainWorker"`
	TestData        []string             `json:"testData"`
	TestDataOpener  string               `json:"testDataOpener"`
	TestWorker      string               `json:"testWorker"`
	Status          string               `json:"status"`
	Rank            int                  `json:"rank"`
	Perf            float32              `json:"perf"`
	TrainPerf       map[string]float32   `json:"trainPerf"`
	TestPerf        map[string]float32   `json:"testPerf"`
	Log             string               `json:"log"`
	Permissions     string               `json:"permissions"`
	Creator         string               `json:"creator"`
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

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	// err := stub.PutState(args[0], []byte(args[1]))
	// if err != nil {
	// 	return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	// }
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
	case "addProblem":
		result, err = addProblem(stub, args)
	case "addData":
		result, err = addData(stub, args)
	case "addAlgo":
		result, err = addAlgo(stub, args)
	case "query":
		result, err = query(stub, args)
	case "queryProblem":
		result, err = queryAll(stub, args, "problem")
	case "queryData":
		result, err = queryAll(stub, args, "data")
	case "queryAlgo":
		result, err = queryAll(stub, args, "algo")
	default:
		err = fmt.Errorf("function not implemented")
	}
	// Return the result as success payload
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

// addProblema stores a new problem in the ledger.
// If the key exists, it will override the value with the new one
func addProblem(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 7 {
		return nil, fmt.Errorf("incorrect arguments, expecting 7 args: " +
			"description hash, name, description storage address, metrics storage address, " +
			"metrics hash, list of test data, permissions")
	}

	// TODO check input types and check if data exist and are from the same center
	testData := strings.Split(strings.Replace(args[5], " ", "", -1), ",")
	// create problem key
	key := "problem_" + args[0]
	// find associated owner
	// TODO
	owner := "TODO"
	// create data object
	var problem = Problem{
		Name: args[1],
		DescriptionStorageAddress: args[2],
		MetricsStorageAddress:     args[3],
		MetricsHash:               args[4],
		Owner:                     owner,
		TestData:                  testData,
		Permissions:               args[6]}
	problemBytes, _ := json.Marshal(problem)
	err := stub.PutState(key, problemBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add problem with description hash %s", args[0])
	}
	return problemBytes, nil
}

// addData stores a new data in the ledger.
// If the key exists, it will override the value with the new one
func addData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 5 {
		return nil, fmt.Errorf("incorrect arguments, expecting 5 args: " +
			"data hash, name, data opener hash, associated problems, permissions")
	}

	// TODO check input types + check if problems are in the ledger
	problems := strings.Split(strings.Replace(args[3], " ", "", -1), ",")
	// create data key
	key := "data_" + args[0]
	// find associated owner
	// TODO
	owner := "TODO"
	// create data object
	var data = Data{
		Name:        args[1],
		DataOpener:  args[2],
		Owner:       owner,
		Problems:    problems,
		Permissions: args[4]}
	dataBytes, _ := json.Marshal(data)
	err := stub.PutState(key, dataBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add data with hash %s", args[0])
	}
	// create composite keys (one for each associated problem) to find data associated with a problem
	indexName := "data~problem~key"
	for _, problem := range data.Problems {
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
		Name:           args[1],
		StorageAddress: args[2],
		Owner:          owner,
		Problem:        problem,
		Permissions:    args[4]}
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

// addTrainTuple add a Train Tuple in the ledger
// ....
func addTrainTuple(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting the key of the element to query")
	}

	// find associated creator and check permissions (TODO later)
	creator := "TODO"

	problemKey := args[0]
	startModelKey := args[1]
	trainDataKeys := strings.Split(strings.Replace(args[2], " ", "", -1), ",")

	// check if train data exist and are from the same center with the same data opener
	// derive trainDataOpener and trainWorker
	var trainWorker, trainDataOpener string
	for i, dataKey := range trainDataKeys {
		dataBytes, err := stub.GetState(dataKey)
		if err != nil {
			return nil, err
		} else if dataBytes == nil {
			return nil, fmt.Errorf("no data with this key %s", dataKey)
		}
		data := Data{}
		err = json.Unmarshal(dataBytes, &data)
		if i == 0 {
			trainWorker = data.Owner
			trainDataOpener = data.DataOpener
		}
		if data.Owner != trainWorker {
			return nil, fmt.Errorf("data do not come from the same center...")
		}
	}

	// get algo key of start model from previous learnuplets
	// TODO
	compositeIterator, err := stub.GetStateByPartialCompositeKey("trainTuple~endModel~key", []string{"trainTuple", startModelKey})
	if err != nil {
		return nil, err
	}
	defer compositeIterator.Close()
	var i, rank int
	var parentTrainTupleKey string
	var testDataKeys []string
	problem := make(map[string][2]string)
	algo := make(map[string]string)
	startModel := make(map[string]string)
	for i = 0; compositeIterator.HasNext(); i++ {
		compositeKey, err := compositeIterator.Next()
		if err != nil {
			return nil, err
		}
		// get the color and name from color~name composite key
		_, compositeKeyParts, err := stub.SplitCompositeKey(compositeKey.Key)
		if err != nil {
			return nil, err
		}
		parentTrainTupleKey = compositeKeyParts[2]
		fmt.Printf("- found %s with endModel %s\n", parentTrainTupleKey, startModelKey)
	}
	if i > 2 {
		return nil, fmt.Errorf("several models associated with start model hash")
	} else if i == 1 {
		// model derives from a previous TrainTuple
		parentTrainTuple := TrainTuple{}
		trainTupleBytes, err := stub.GetState(parentTrainTupleKey)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(trainTupleBytes, &parentTrainTuple)
		if err != nil {
			return nil, err
		}
		algo = parentTrainTuple.Algo
		startModel = parentTrainTuple.EndModel
		rank = parentTrainTuple.Rank + 1
		problem = parentTrainTuple.Problem
		testDataKeys = parentTrainTuple.TestData
	} else {
		// first time algo is trained
		rank = 0
		// get problem to derive metrics info and test data keys
		problemBytes, err := stub.GetState(problemKey)
		if err != nil {
			return nil, err
		} else if problemBytes == nil {
			return nil, fmt.Errorf("no problem with this key %s", problemKey)
		}
		retrievedProblem := Problem{}
		err = json.Unmarshal(problemBytes, &retrievedProblem)
		if err != nil {
			return nil, err
		}
		testDataKeys = retrievedProblem.TestData
		problem[problemKey] = [2]string{retrievedProblem.MetricsHash, retrievedProblem.MetricsStorageAddress}
		// get algo
		algoBytes, err := stub.GetState(startModelKey)
		if err != nil {
			return nil, err
		} else if algoBytes == nil {
			return nil, fmt.Errorf("no algo with this key %s", startModelKey)
		}
		retrievedAlgo := Algo{}
		err = json.Unmarshal(algoBytes, &retrievedAlgo)
		if err != nil {
			return nil, err
		}
		algo[startModelKey] = retrievedAlgo.StorageAddress
	}

	// get testWorker given test data and test data opener
	// for now, we assume that test data of a problem are all located in the same center
	// and have the same data opener
	testDataBytes, err := stub.GetState(testDataKeys[0])
	if err != nil {
		return nil, err
	}
	testData := Data{}
	err = json.Unmarshal(testDataBytes, &testData)
	if err != nil {
		return nil, err
	}
	testWorker := testData.Owner
	testDataOpener := testData.DataOpener

	// create learnuplet and add it to ledger
	var trainTuple = TrainTuple{
		Problem:         problem,
		Algo:            algo,
		StartModel:      startModel,
		EndModel:        make(map[string]string),
		TrainData:       trainDataKeys,
		TrainDataOpener: trainDataOpener,
		TrainWorker:     trainWorker,
		TestData:        testDataKeys,
		TestDataOpener:  testDataOpener,
		TestWorker:      testWorker,
		Status:          "todo",
		Rank:            rank,
		Perf:            0.0,
		TrainPerf:       make(map[string]float32),
		TestPerf:        make(map[string]float32),
		Log:             "",
		Permissions:     "all",
		Creator:         creator,
	}
	// TODO create key
	key := "todo"
	trainTupleBytes, _ := json.Marshal(trainTuple)
	err = stub.PutState(key, trainTupleBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add trainTuple with startModel %s and problem %s", startModelKey, problemKey)
	}
	return trainTupleBytes, nil
}

// query returns an element of the ledger given its key
// For now, ok for everything. Later returns if the requester has permission to see it
func query(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key string

	if len(args) != 1 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting the key of the element to query")
	}

	key = args[0]
	valBytes, err := stub.GetState(key)
	if err != nil {
		return nil, err
	} else if valBytes == nil {
		return nil, fmt.Errorf("no element with this key %s", key)
	}

	return valBytes, nil
}

// queryAll returns all elements of the ledger given its type
// For now, ok for everything. Later returns if the requester has permission to see it
func queryAll(stub shim.ChaincodeStubInterface, args []string, elementType string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}

	elementsIterator, err := stub.GetStateByRange(elementType+"_", elementType+"a")
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
		err = json.Unmarshal(queryResponse.GetValue(), &element)
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
