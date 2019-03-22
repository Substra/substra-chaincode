package main

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gopkg.in/go-playground/validator.v9"
)

// Set is a method of the receiver DataManager. It checks the validity of inputDataManager and uses its fields to set the DataManager
// Returns the dataManagerKey and associated objectiveKeys
func (dataManager *DataManager) Set(stub shim.ChaincodeStubInterface, inp inputDataManager) (string, string, error) {
	// check validity of submitted fields
	validate := validator.New()
	if err := validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid dataManager inputs %s", err.Error())
		return "", "", err
	}
	// check dataManager is not already in the ledger
	dataManagerKey := inp.OpenerHash
	if elementBytes, _ := stub.GetState(dataManagerKey); elementBytes != nil {
		return dataManagerKey, "", fmt.Errorf("dataManager with this opener already exists")
	}
	// check validity of associated objective
	if len(inp.ObjectiveKey) > 0 {
		if _, err := getElementBytes(stub, inp.ObjectiveKey); err != nil {
			err = fmt.Errorf("error checking associated objective %s", err.Error())
			return "", "", nil
		}
		dataManager.ObjectiveKey = inp.ObjectiveKey
	}
	dataManager.Name = inp.Name
	dataManager.OpenerStorageAddress = inp.OpenerStorageAddress
	dataManager.Type = inp.Type
	dataManager.Description = &HashDress{
		Hash:           inp.DescriptionHash,
		StorageAddress: inp.DescriptionStorageAddress,
	}
	owner, err := getTxCreator(stub)
	if err != nil {
		return "", "", err
	}
	dataManager.Owner = owner
	dataManager.Permissions = inp.Permissions
	return dataManagerKey, dataManager.ObjectiveKey, nil
}

// setData is a method checking the validity of inputData to be registered in the ledger
// and returning corresponding data hashes, associated dataManagers, testOnly and errors
func setData(stub shim.ChaincodeStubInterface, inp inputData) (dataHashes []string, data Data, err error) {
	// validate input data
	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid data inputs %s", err.Error())
		return
	}
	// Get data keys (=hashes)
	dataHashes = strings.Split(strings.Replace(inp.Hashes, " ", "", -1), ",")
	// check validity of dataHashes
	if err = checkHashes(dataHashes); err != nil {
		return
	}
	// check data is not already in the ledger
	if existingKeys := checkExist(stub, dataHashes); existingKeys != nil {
		err = fmt.Errorf("data with these hashes already exist - %s", existingKeys)
		return
	}

	// get transaction owner
	owner, err := getTxCreator(stub)
	if err != nil {
		return
	}
	// check if associated dataManager(s) exists
	var dataManagerKeys []string
	if len(inp.DataManagerKeys) > 0 {
		dataManagerKeys = strings.Split(strings.Replace(inp.DataManagerKeys, " ", "", -1), ",")
		if err = checkDataManagerOwner(stub, dataManagerKeys); err != nil {
			return
		}
	}
	// convert input testOnly to boolean
	testOnly, err := strconv.ParseBool(inp.TestOnly)

	data = Data{
		DataManagerKeys: dataManagerKeys,
		TestOnly:    testOnly,
		Owner:       owner}

	return
}

// validateUpdateData is a method checking the validity of elements sent to update
// one or more dataf
func validateUpdateData(stub shim.ChaincodeStubInterface, inp inputUpdateData) (dataHashes []string, dataManagerKeys []string, err error) {

	// TODO return full data

	// validate input to updatedata
	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid inputs to update data %s", err.Error())
		return
	}
	// Get data keys (=hashes)
	dataHashes = strings.Split(strings.Replace(inp.Hashes, " ", "", -1), ",")
	// check validity of dataHashes
	if err = checkHashes(dataHashes); err != nil {
		return
	}
	// check dataManagers exist and are owned by the transaction requester
	dataManagerKeys = strings.Split(strings.Replace(inp.DataManagerKeys, " ", "", -1), ",")
	if err = checkDataManagerOwner(stub, dataManagerKeys); err != nil {
		return
	}
	return
}

// -----------------------------------------------------------------
// ----------------------- Smart Contracts  ------------------------
// -----------------------------------------------------------------

// registerDataManager stores a new dataManager in the ledger.
func registerDataManager(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputDataManager{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}
	// convert input strings args to input struct inputDataManager
	inp := inputDataManager{}
	stringToInputStruct(args, &inp)
	// check validity of input args and convert it to a DataManager
	dataManager := DataManager{}
	dataManagerKey, objectiveKey, err := dataManager.Set(stub, inp)
	if err != nil {
		return nil, err
	}
	// submit to ledger
	dataManagerBytes, _ := json.Marshal(dataManager)
	err = stub.PutState(dataManagerKey, dataManagerBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add dataManager with opener hash %s, error is %s", inp.OpenerHash, err.Error())
	}
	// create composite keys (one for each associated objective) to find data associated with a objective
	indexName := "dataManager~objective~key"
	err = createCompositeKey(stub, indexName, []string{"dataManager", objectiveKey, dataManagerKey})
	if err != nil {
		return nil, err
	}
	// create composite key to find dataManager associated with a owner
	err = createCompositeKey(stub, "dataManager~owner~key", []string{"dataManager", dataManager.Owner, dataManagerKey})
	if err != nil {
		return nil, err
	}
	return []byte(dataManagerKey), nil
}

// registerData stores new data in the ledger (one or more).
func registerData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputData{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputData
	inp := inputData{}
	stringToInputStruct(args, &inp)
	// check validity of input args
	dataHashes, data, err := setData(stub, inp)
	if err != nil {
		return nil, err
	}

	// Serialized Data object
	dataBytes, _ := json.Marshal(data)

	// store data in the ledger
	var dataKeys string
	suffix := ", "
	for _, dataHash := range dataHashes {
		dataKeys = dataKeys + "\"" + dataHash + "\"" + suffix
		if err = stub.PutState(dataHash, dataBytes); err != nil {
			return nil, fmt.Errorf("failed to add data with hash %s", dataHash)
		}
		for _, dataManagerKey := range data.DataManagerKeys {
			// create composite keys to find all data associated with a dataManager and both test and train data
			if err = createCompositeKey(stub, "data~dataManager~key", []string{"data", dataManagerKey, dataHash}); err != nil {
				return nil, err
			}
			// create composite keys to find all data associated with a dataManager and only test or train data
			if err = createCompositeKey(stub, "data~dataManager~testOnly~key", []string{"data", dataManagerKey, strconv.FormatBool(data.TestOnly), dataHash}); err != nil {
				return nil, err
			}
		}
	}
	// return added data keys
	dataKeys = "{\"keys\": [" + strings.TrimSuffix(dataKeys, suffix) + "]}"
	return []byte(dataKeys), nil
}

// updateData associates one or more dataManagerKeys to one or more data
func updateData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	inp := inputUpdateData{}
	expectedArgs := getFieldNames(&inp)
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputUpdateData
	stringToInputStruct(args, &inp)
	// check validity of input args
	dataHashes, dataManagerKeys, err := validateUpdateData(stub, inp)
	if err != nil {
		return nil, err
	}
	// store data in the ledger
	var dataKeys string
	suffix := ", "
	for _, dataHash := range dataHashes {
		dataKeys = dataKeys + "\"" + dataHash + "\"" + suffix
		data := Data{}
		if err = getElementStruct(stub, dataHash, &data); err != nil {
			err = fmt.Errorf("could not retrieve data with key %s - %s", dataHash, err.Error())
			return nil, err
		}
		if err = checkDataOwner(stub, data); err != nil {
			return nil, err
		}
		for _, dataManagerKey := range dataManagerKeys {
			if !stringInSlice(dataManagerKey, data.DataManagerKeys) {
				data.DataManagerKeys = append(data.DataManagerKeys, dataManagerKey)
			}
		}
		dataBytes, _ := json.Marshal(data)
		if err = stub.PutState(dataHash, dataBytes); err != nil {
			err = fmt.Errorf("could not put in ledger data with key %s - %s", dataHash, err.Error())
			return nil, err
		}

	}
	// return updated data keys
	dataKeys = "{\"keys\": [" + strings.TrimSuffix(dataKeys, suffix) + "]}"
	return []byte(dataKeys), nil
}

// updateDataManager associates a objectiveKey to an existing dataManager
func updateDataManager(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	inp := inputUpdateDataManager{}
	expectedArgs := getFieldNames(&inp)
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputData
	stringToInputStruct(args, &inp)
	validate := validator.New()
	if err := validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid update dataManager inputs %s", err.Error())
		return nil, err
	}
	// update dataManager.ObjectiveKey
	if err := addObjectiveDataManager(stub, inp.DataManagerKey, inp.ObjectiveKey); err != nil {
		return nil, err
	}
	return []byte(inp.DataManagerKey), nil
}

// queryObjective returns a objective of the ledger given its key
func queryDataManager(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 || len(args[0]) != 64 {
		return nil, fmt.Errorf("incorrect arguments, expecting key, received: %s", args[0])
	}
	key := args[0]
	var dataManager DataManager
	if err := getElementStruct(stub, key, &dataManager); err != nil {
		return nil, err
	}
	var out outputDataManager
	out.Fill(key, dataManager)
	return json.Marshal(out)
}

// queryObjectives returns all objectives of the ledger
func queryDataManagers(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("incorrect number of arguments, expecting nothing")
	}
	var indexName = "dataManager~owner~key"
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{"dataManager"})
	if err != nil {
		return nil, fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
	}
	var outDataManagers []outputDataManager
	for _, key := range elementsKeys {
		var objective DataManager
		if err := getElementStruct(stub, key, &objective); err != nil {
			return nil, err
		}
		var out outputDataManager
		out.Fill(key, objective)
		outDataManagers = append(outDataManagers, out)
	}
	return json.Marshal(outDataManagers)
}

// queryDataset returns info about a dataManager and all related data
func queryDataset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"dataManager key"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	dataManagerKey := args[0]

	// get dataManager info
	var mPayload map[string]interface{}
	if err := getElementStruct(stub, dataManagerKey, &mPayload); err != nil {
		return nil, err
	}
	mPayload["key"] = dataManagerKey
	// get related train data
	trainDataKeys, err := getDataset(stub, dataManagerKey, false)
	if err != nil {
		return nil, err
	}
	mPayload["trainDataKeys"] = trainDataKeys
	// get related test data
	testDataKeys, err := getDataset(stub, dataManagerKey, true)
	if err != nil {
		return nil, err
	}
	mPayload["testDataKeys"] = testDataKeys
	// Marshal payload
	payload, err := json.Marshal(mPayload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// -----------------------------------------------------------------
// -------------------- Data / DataManager utils -----------------------
// -----------------------------------------------------------------

// check

// checkDataManagerOwner checks if the transaction requester is the owner of dataManager
// specified by their keys in a slice
func checkDataManagerOwner(stub shim.ChaincodeStubInterface, dataManagerKeys []string) (err error) {
	// get transaction requester
	txCreator, err := getTxCreator(stub)
	if err != nil {
		return
	}
	for _, dataManagerKey := range dataManagerKeys {
		dataManager := DataManager{}
		if err = getElementStruct(stub, dataManagerKey, &dataManager); err != nil {
			err = fmt.Errorf("could not retrieve dataManager with key %s - %s", dataManagerKey, err.Error())
			return
		}
		// check transaction requester is the dataManager owner
		if txCreator != dataManager.Owner {
			err = fmt.Errorf("%s is not the owner of the dataManager %s", txCreator, dataManagerKey)
			return
		}
	}
	return
}

//  checkDataOwner checks if the transaction requester is the owner of the data
func checkDataOwner(stub shim.ChaincodeStubInterface, data Data) (err error) {
	txRequester, err := getTxCreator(stub)
	if err != nil {
		return
	}
	if txRequester != data.Owner {
		err = fmt.Errorf("%s is not the data's owner", txRequester)
	}
	return
}

// checkSameDataManager checks if data in a slice exist and are from the same dataManager.
// If yes, returns two boolean indicating if data are testOnly and trainOnly
func checkSameDataManager(stub shim.ChaincodeStubInterface, dataManagerKey string, dataKeys []string) (testOnly bool, trainOnly bool, err error) {
	testOnly = true
	trainOnly = true
	for _, dataKey := range dataKeys {
		data := Data{}
		if err = getElementStruct(stub, dataKey, &data); err != nil {
			err = fmt.Errorf("could not retrieve data with key %s - %s", dataKey, err.Error())
			return
		}
		if !stringInSlice(dataManagerKey, data.DataManagerKeys) {
			err = fmt.Errorf("data do not belong to the same dataManager")
			return
		}
		testOnly = testOnly && data.TestOnly
		trainOnly = trainOnly && !data.TestOnly
	}
	return
}

// getDataset returns all data keys associated to a dataManager
func getDataset(stub shim.ChaincodeStubInterface, dataManagerKey string, testOnly bool) ([]string, error) {
	indexName := "data~dataManager~testOnly~key"
	attributes := []string{"data", dataManagerKey, strconv.FormatBool(testOnly)}
	dataKeys, err := getKeysFromComposite(stub, indexName, attributes)
	if err != nil {
		return nil, err
	}
	return dataKeys, nil
}

// getDataManagerOwner returns the owner of a dataManager given its key
func getDataManagerOwner(stub shim.ChaincodeStubInterface, dataManagerKey string) (worker string, err error) {

	dataManager := DataManager{}
	if err = getElementStruct(stub, dataManagerKey, &dataManager); err != nil {
		err = fmt.Errorf("could not retrieve dataManager with key %s - %s", dataManagerKey, err.Error())
		return
	}
	worker = dataManager.Owner
	return
}
