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

// setDataSample is a method checking the validity of inputDataSample to be registered in the ledger
// and returning corresponding dataSample hashes, associated dataManagers, testOnly and errors
func setDataSample(stub shim.ChaincodeStubInterface, inp inputDataSample) (dataSampleHashes []string, dataSample DataSample, err error) {
	// validate input dataSample
	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid dataSample inputs %s", err.Error())
		return
	}
	// Get dataSample keys (=hashes)
	dataSampleHashes = strings.Split(strings.Replace(inp.Hashes, " ", "", -1), ",")
	// check validity of dataSampleHashes
	if err = checkHashes(dataSampleHashes); err != nil {
		return
	}
	// check dataSample is not already in the ledger
	if existingKeys := checkExist(stub, dataSampleHashes); existingKeys != nil {
		err = fmt.Errorf("dataSample with these hashes already exist - %s", existingKeys)
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

	dataSample = DataSample{
		DataManagerKeys: dataManagerKeys,
		TestOnly:        testOnly,
		Owner:           owner}

	return
}

// validateUpdateDataSample is a method checking the validity of elements sent to update
// one or more dataSamplef
func validateUpdateDataSample(stub shim.ChaincodeStubInterface, inp inputUpdateDataSample) (dataSampleHashes []string, dataManagerKeys []string, err error) {

	// TODO return full dataSample

	// validate input to updatedataSample
	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid inputs to update dataSample %s", err.Error())
		return
	}
	// Get dataSample keys (=hashes)
	dataSampleHashes = strings.Split(strings.Replace(inp.Hashes, " ", "", -1), ",")
	// check validity of dataSampleHashes
	if err = checkHashes(dataSampleHashes); err != nil {
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
	// create composite keys (one for each associated objective) to find dataSample associated with a objective
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

// registerDataSample stores new dataSample in the ledger (one or more).
func registerDataSample(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputDataSample{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputDataSample
	inp := inputDataSample{}
	stringToInputStruct(args, &inp)
	// check validity of input args
	dataSampleHashes, dataSample, err := setDataSample(stub, inp)
	if err != nil {
		return nil, err
	}

	// Serialized DataSample object
	dataSampleBytes, _ := json.Marshal(dataSample)

	// store dataSample in the ledger
	var dataSampleKeys string
	suffix := ", "
	for _, dataSampleHash := range dataSampleHashes {
		dataSampleKeys = dataSampleKeys + "\"" + dataSampleHash + "\"" + suffix
		if err = stub.PutState(dataSampleHash, dataSampleBytes); err != nil {
			return nil, fmt.Errorf("failed to add dataSample with hash %s", dataSampleHash)
		}
		for _, dataManagerKey := range dataSample.DataManagerKeys {
			// create composite keys to find all dataSample associated with a dataManager and both test and train dataSample
			if err = createCompositeKey(stub, "dataSample~dataManager~key", []string{"dataSample", dataManagerKey, dataSampleHash}); err != nil {
				return nil, err
			}
			// create composite keys to find all dataSample associated with a dataManager and only test or train dataSample
			if err = createCompositeKey(stub, "dataSample~dataManager~testOnly~key", []string{"dataSample", dataManagerKey, strconv.FormatBool(dataSample.TestOnly), dataSampleHash}); err != nil {
				return nil, err
			}
		}
	}
	// return added dataSample keys
	dataSampleKeys = "{\"keys\": [" + strings.TrimSuffix(dataSampleKeys, suffix) + "]}"
	return []byte(dataSampleKeys), nil
}

// updateDataSample associates one or more dataManagerKeys to one or more dataSample
func updateDataSample(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	inp := inputUpdateDataSample{}
	expectedArgs := getFieldNames(&inp)
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputUpdateDataSample
	stringToInputStruct(args, &inp)
	// check validity of input args
	dataSampleHashes, dataManagerKeys, err := validateUpdateDataSample(stub, inp)
	if err != nil {
		return nil, err
	}
	// store dataSample in the ledger
	var dataSampleKeys string
	suffix := ", "
	for _, dataSampleHash := range dataSampleHashes {
		dataSampleKeys = dataSampleKeys + "\"" + dataSampleHash + "\"" + suffix
		dataSample := DataSample{}
		if err = getElementStruct(stub, dataSampleHash, &dataSample); err != nil {
			err = fmt.Errorf("could not retrieve dataSample with key %s - %s", dataSampleHash, err.Error())
			return nil, err
		}
		if err = checkDataSampleOwner(stub, dataSample); err != nil {
			return nil, err
		}
		for _, dataManagerKey := range dataManagerKeys {
			if !stringInSlice(dataManagerKey, dataSample.DataManagerKeys) {
				dataSample.DataManagerKeys = append(dataSample.DataManagerKeys, dataManagerKey)
			}
		}
		dataSampleBytes, _ := json.Marshal(dataSample)
		if err = stub.PutState(dataSampleHash, dataSampleBytes); err != nil {
			err = fmt.Errorf("could not put in ledger dataSample with key %s - %s", dataSampleHash, err.Error())
			return nil, err
		}

	}
	// return updated dataSample keys
	dataSampleKeys = "{\"keys\": [" + strings.TrimSuffix(dataSampleKeys, suffix) + "]}"
	return []byte(dataSampleKeys), nil
}

// updateDataManager associates a objectiveKey to an existing dataManager
func updateDataManager(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	inp := inputUpdateDataManager{}
	expectedArgs := getFieldNames(&inp)
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputDataSample
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

// queryDataset returns info about a dataManager and all related dataSample
func queryDataset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 || len(args[0]) != 64 {
		return nil, fmt.Errorf("incorrect arguments, expecting key, received: %s", args[0])
	}
	key := args[0]
	var dataManager DataManager
	if err := getElementStruct(stub, key, &dataManager); err != nil {
		return nil, err
	}

	// get related train dataSample
	trainDataSampleKeys, err := getDataset(stub, key, false)
	if err != nil {
		return nil, err
	}

	// get related test dataSample
	testDataSampleKeys, err := getDataset(stub, key, true)
	if err != nil {
		return nil, err
	}

	var out outputDataset
	out.Fill(key, dataManager, trainDataSampleKeys, testDataSampleKeys)
	return json.Marshal(out)
}

// -----------------------------------------------------------------
// -------------------- DataSample / DataManager utils -----------------------
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

//  checkDataSampleOwner checks if the transaction requester is the owner of the dataSample
func checkDataSampleOwner(stub shim.ChaincodeStubInterface, dataSample DataSample) (err error) {
	txRequester, err := getTxCreator(stub)
	if err != nil {
		return
	}
	if txRequester != dataSample.Owner {
		err = fmt.Errorf("%s is not the dataSample's owner", txRequester)
	}
	return
}

// checkSameDataManager checks if dataSample in a slice exist and are from the same dataManager.
// If yes, returns two boolean indicating if dataSample are testOnly and trainOnly
func checkSameDataManager(stub shim.ChaincodeStubInterface, dataManagerKey string, dataSampleKeys []string) (testOnly bool, trainOnly bool, err error) {
	testOnly = true
	trainOnly = true
	for _, dataSampleKey := range dataSampleKeys {
		dataSample := DataSample{}
		if err = getElementStruct(stub, dataSampleKey, &dataSample); err != nil {
			err = fmt.Errorf("could not retrieve dataSample with key %s - %s", dataSampleKey, err.Error())
			return
		}
		if !stringInSlice(dataManagerKey, dataSample.DataManagerKeys) {
			err = fmt.Errorf("dataSample do not belong to the same dataManager")
			return
		}
		testOnly = testOnly && dataSample.TestOnly
		trainOnly = trainOnly && !dataSample.TestOnly
	}
	return
}

// getDataset returns all dataSample keys associated to a dataManager
func getDataset(stub shim.ChaincodeStubInterface, dataManagerKey string, testOnly bool) ([]string, error) {
	indexName := "dataSample~dataManager~testOnly~key"
	attributes := []string{"dataSample", dataManagerKey, strconv.FormatBool(testOnly)}
	dataSampleKeys, err := getKeysFromComposite(stub, indexName, attributes)
	if err != nil {
		return nil, err
	}
	return dataSampleKeys, nil
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
