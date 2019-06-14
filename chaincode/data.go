package main

import (
	"chaincode/errors"
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
		err = errors.BadRequest(err, "invalid dataManager inputs")
		return "", "", err
	}
	// check dataManager is not already in the ledger
	dataManagerKey := inp.OpenerHash
	if elementBytes, _ := stub.GetState(dataManagerKey); elementBytes != nil {
		err := errors.Conflict("dataManager with this opener already exists")
		return dataManagerKey, "", err
	}
	// check validity of associated objective
	if len(inp.ObjectiveKey) > 0 {
		if _, err := getElementBytes(stub, inp.ObjectiveKey); err != nil {
			err = errors.Internal(err, "error checking associated objective")
			return "", "", err
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
		err = errors.BadRequest(err, "invalid dataSample inputs")
		return
	}
	// Get dataSample keys (=hashes)
	dataSampleHashes = strings.Split(strings.Replace(inp.Hashes, " ", "", -1), ",")
	// check validity of dataSampleHashes
	if err = checkHashes(dataSampleHashes); err != nil {
		err = errors.BadRequest(err)
		return
	}
	// check dataSample is not already in the ledger
	if existingKeys := checkExist(stub, dataSampleHashes); existingKeys != nil {
		err = errors.Conflict("dataSample with these hashes already exist - %s", existingKeys)
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
		err = errors.BadRequest(err, "invalid inputs to update dataSample")
		return
	}
	// Get dataSample keys (=hashes)
	dataSampleHashes = strings.Split(strings.Replace(inp.Hashes, " ", "", -1), ",")
	// check validity of dataSampleHashes
	if err = checkHashes(dataSampleHashes); err != nil {
		err = errors.BadRequest(err)
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
func registerDataManager(stub shim.ChaincodeStubInterface, args []string) (resp map[string]string, err error) {
	inp := inputDataManager{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	// check validity of input args and convert it to a DataManager
	dataManager := DataManager{}
	dataManagerKey, objectiveKey, err := dataManager.Set(stub, inp)
	if err != nil {
		return
	}
	// submit to ledger
	dataManagerBytes, _ := json.Marshal(dataManager)
	err = stub.PutState(dataManagerKey, dataManagerBytes)
	if err != nil {
		err = errors.Internal(err, "failed to add dataManager with opener hash %s", inp.OpenerHash)
		return
	}
	// create composite keys (one for each associated objective) to find dataSample associated with a objective
	indexName := "dataManager~objective~key"
	err = createCompositeKey(stub, indexName, []string{"dataManager", objectiveKey, dataManagerKey})
	if err != nil {
		return
	}
	// create composite key to find dataManager associated with a owner
	err = createCompositeKey(stub, "dataManager~owner~key", []string{"dataManager", dataManager.Owner, dataManagerKey})
	if err != nil {
		return
	}
	return map[string]string{"key": dataManagerKey}, nil
}

// registerDataSample stores new dataSample in the ledger (one or more).
func registerDataSample(stub shim.ChaincodeStubInterface, args []string) (dataSampleKeys map[string][]string, err error) {
	// convert input strings args to input struct inputDataSample
	inp := inputDataSample{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	// check validity of input args
	dataSampleHashes, dataSample, err := setDataSample(stub, inp)
	if err != nil {
		return
	}

	// Serialized DataSample object
	dataSampleBytes, _ := json.Marshal(dataSample)

	// store dataSample in the ledger
	for _, dataSampleHash := range dataSampleHashes {
		if err = stub.PutState(dataSampleHash, dataSampleBytes); err != nil {
			err = fmt.Errorf("failed to add dataSample with hash %s", dataSampleHash)
			return
		}
		for _, dataManagerKey := range dataSample.DataManagerKeys {
			// create composite keys to find all dataSample associated with a dataManager and both test and train dataSample
			if err = createCompositeKey(stub, "dataSample~dataManager~key", []string{"dataSample", dataManagerKey, dataSampleHash}); err != nil {
				return
			}
			// create composite keys to find all dataSample associated with a dataManager and only test or train dataSample
			if err = createCompositeKey(stub, "dataSample~dataManager~testOnly~key", []string{"dataSample", dataManagerKey, strconv.FormatBool(dataSample.TestOnly), dataSampleHash}); err != nil {
				return
			}
		}
	}
	// return added dataSample keys
	dataSampleKeys = map[string][]string{"keys": dataSampleHashes}
	return
}

// updateDataSample associates one or more dataManagerKeys to one or more dataSample
func updateDataSample(stub shim.ChaincodeStubInterface, args []string) (resp map[string]string, err error) {
	inp := inputUpdateDataSample{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	// check validity of input args
	dataSampleHashes, dataManagerKeys, err := validateUpdateDataSample(stub, inp)
	if err != nil {
		return
	}
	// store dataSample in the ledger
	var dataSampleKeys string
	suffix := ", "
	for _, dataSampleHash := range dataSampleHashes {
		dataSampleKeys = dataSampleKeys + "\"" + dataSampleHash + "\"" + suffix
		dataSample := DataSample{}
		if err = getElementStruct(stub, dataSampleHash, &dataSample); err != nil {
			err = fmt.Errorf("could not retrieve dataSample with key %s - %s", dataSampleHash, err.Error())
			return
		}
		if err = checkDataSampleOwner(stub, dataSample); err != nil {
			return
		}
		for _, dataManagerKey := range dataManagerKeys {
			if !stringInSlice(dataManagerKey, dataSample.DataManagerKeys) {
				// check data manager is not already associated with this data
				dataSample.DataManagerKeys = append(dataSample.DataManagerKeys, dataManagerKey)
				// create composite keys to find all dataSample associated with a dataManager and both test and train dataSample
				if err = createCompositeKey(stub, "dataSample~dataManager~key", []string{"dataSample", dataManagerKey, dataSampleHash}); err != nil {
					return
				}
				// create composite keys to find all dataSample associated with a dataManager and only test or train dataSample
				if err = createCompositeKey(stub, "dataSample~dataManager~testOnly~key", []string{"dataSample", dataManagerKey, strconv.FormatBool(dataSample.TestOnly), dataSampleHash}); err != nil {
					return
				}
			}
		}
		dataSampleBytes, _ := json.Marshal(dataSample)
		if err = stub.PutState(dataSampleHash, dataSampleBytes); err != nil {
			err = fmt.Errorf("could not put in ledger dataSample with key %s - %s", dataSampleHash, err.Error())
			return
		}

	}
	// return updated dataSample keys
	dataSampleKeys = "{\"keys\": [" + strings.TrimSuffix(dataSampleKeys, suffix) + "]}"
	return map[string]string{"key": dataSampleKeys}, nil
}

// updateDataManager associates a objectiveKey to an existing dataManager
func updateDataManager(stub shim.ChaincodeStubInterface, args []string) (resp map[string]string, err error) {
	inp := inputUpdateDataManager{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid update dataManager inputs %s", err.Error())
		return
	}
	// update dataManager.ObjectiveKey
	if err = addObjectiveDataManager(stub, inp.DataManagerKey, inp.ObjectiveKey); err != nil {
		return
	}
	return map[string]string{"key": inp.DataManagerKey}, nil
}

// queryObjective returns a objective of the ledger given its key
func queryDataManager(stub shim.ChaincodeStubInterface, args []string) (out outputDataManager, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	var dataManager DataManager
	if err = getElementStruct(stub, inp.Key, &dataManager); err != nil {
		return
	}
	out.Fill(inp.Key, dataManager)
	return
}

// queryObjectives returns all objectives of the ledger
func queryDataManagers(stub shim.ChaincodeStubInterface, args []string) (outDataManagers []outputDataManager, err error) {
	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}
	var indexName = "dataManager~owner~key"
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{"dataManager"})
	if err != nil {
		err = fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
		return
	}
	for _, key := range elementsKeys {
		var objective DataManager
		if err = getElementStruct(stub, key, &objective); err != nil {
			return
		}
		var out outputDataManager
		out.Fill(key, objective)
		outDataManagers = append(outDataManagers, out)
	}
	return
}

// queryDataset returns info about a dataManager and all related dataSample
func queryDataset(stub shim.ChaincodeStubInterface, args []string) (out outputDataset, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	var dataManager DataManager
	if err = getElementStruct(stub, inp.Key, &dataManager); err != nil {
		return
	}

	// get related train dataSample
	trainDataSampleKeys, err := getDataset(stub, inp.Key, false)
	if err != nil {
		return
	}

	// get related test dataSample
	testDataSampleKeys, err := getDataset(stub, inp.Key, true)
	if err != nil {
		return
	}

	out.Fill(inp.Key, dataManager, trainDataSampleKeys, testDataSampleKeys)
	return
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
			err = errors.BadRequest(err, "could not retrieve dataManager with key %s", dataManagerKey)
			return
		}
		// check transaction requester is the dataManager owner
		if txCreator != dataManager.Owner {
			err = errors.Forbidden("%s is not the owner of the dataManager %s", txCreator, dataManagerKey)
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
		err = errors.Forbidden("%s is not the dataSample's owner", txRequester)
		return
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
		err = errors.BadRequest(err, "could not retrieve dataManager with key %s -", dataManagerKey)
		return
	}
	worker = dataManager.Owner
	return
}
