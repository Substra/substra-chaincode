package main

import (
	"fmt"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gopkg.in/go-playground/validator.v9"
)

// Set is a method of the receiver Dataset. It checks the validity of inputDataset and uses its fields to set the Dataset
// Returns the datasetKey and associated challengeKeys
func (dataset *Dataset) Set(stub shim.ChaincodeStubInterface, inp inputDataset) (string, string, error) {
	// check validity of submitted fields
	validate := validator.New()
	if err := validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid dataset inputs %s", err.Error())
		return "", "", err
	}
	// check dataset is not already in the ledger
	datasetKey := inp.OpenerHash
	if elementBytes, _ := stub.GetState(datasetKey); elementBytes != nil {
		return datasetKey, "", fmt.Errorf("dataset with this opener already exists")
	}
	// check validity of associated challenge
	if len(inp.ChallengeKey) > 0 {
		if _, err := getElementBytes(stub, inp.ChallengeKey); err != nil {
			err = fmt.Errorf("error checking associated challenge %s", err.Error())
			return "", "", nil
		}
		dataset.ChallengeKey = inp.ChallengeKey
	}
	dataset.Name = inp.Name
	dataset.OpenerStorageAddress = inp.OpenerStorageAddress
	dataset.Type = inp.Type
	dataset.Description = &HashDress{
		Hash:           inp.DescriptionHash,
		StorageAddress: inp.DescriptionStorageAddress,
	}
	owner, err := getTxCreator(stub)
	if err != nil {
		return "", "", err
	}
	dataset.Owner = owner
	dataset.Permissions = inp.Permissions
	return datasetKey, dataset.ChallengeKey, nil
}

// setData is a method checking the validity of inputData to be registered in the ledger
// and returning corresponding data hashes, associated datasets, testOnly and errors
func setData(stub shim.ChaincodeStubInterface, inp inputData) ([]string, []string, bool, error) {
	// validate input data
	validate := validator.New()
	if err := validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid data inputs %s", err.Error())
		return nil, nil, false, err
	}
	// Get data keys (=hashes)
	dataHashes := strings.Split(strings.Replace(inp.Hashes, " ", "", -1), ",")
	// check validity of dataHashes
	if err := checkHashes(dataHashes); err != nil {
		return nil, nil, false, err
	}
	// check data is not already in the ledger
	if existingKeys := checkExist(stub, dataHashes); existingKeys != nil {
		err := fmt.Errorf("data with these hashes already exist - %s", existingKeys)
		return nil, nil, false, err
	}
	// check if associated dataset(s) exists
	var datasetKeys []string
	if len(inp.DatasetKeys) > 0 {
		datasetKeys = strings.Split(strings.Replace(inp.DatasetKeys, " ", "", -1), ",")
		if err := checkDatasetOwner(stub, datasetKeys); err != nil {
			return nil, nil, false, err
		}
	}
	// convert input testOnly to boolean
	testOnly, err := strconv.ParseBool(inp.TestOnly)
	return dataHashes, datasetKeys, testOnly, err
}

// validateUpdateData is a method checking the validity of elements sent to update
// one or more dataf
func validateUpdateData(stub shim.ChaincodeStubInterface, inp inputUpdateData) (dataHashes []string, datasetKeys []string, err error) {

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
	// check datasets exist and are owned by the transaction requester
	datasetKeys = strings.Split(strings.Replace(inp.DatasetKeys, " ", "", -1), ",")
	if err = checkDatasetOwner(stub, datasetKeys); err != nil {
		return
	}
	return
}

// -----------------------------------------------------------------
// ----------------------- Smart Contracts  ------------------------
// -----------------------------------------------------------------

// registerDataset stores a new dataset in the ledger.
func registerDataset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputDataset{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}
	// convert input strings args to input struct inputDataset
	inp := inputDataset{}
	stringToInputStruct(args, &inp)
	// check validity of input args and convert it to a Dataset
	dataset := Dataset{}
	datasetKey, challengeKey, err := dataset.Set(stub, inp)
	if err != nil {
		return nil, err
	}
	// submit to ledger
	datasetBytes, _ := json.Marshal(dataset)
	err = stub.PutState(datasetKey, datasetBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add dataset with opener hash %s, error is %s", inp.OpenerHash, err.Error())
	}
	// create composite keys (one for each associated challenge) to find data associated with a challenge
	indexName := "dataset~challenge~key"
	err = createCompositeKey(stub, indexName, []string{"dataset", challengeKey, datasetKey})
	if err != nil {
		return nil, err
	}
	// create composite key to find dataset associated with a owner
	err = createCompositeKey(stub, "dataset~owner~key", []string{"dataset", dataset.Owner, datasetKey})
	if err != nil {
		return nil, err
	}
	return []byte(datasetKey), nil
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
	dataHashes, datasetKeys, testOnly, err := setData(stub, inp)
	if err != nil {
		return nil, err
	}

	// create data object
	var data = Data{
		DatasetKeys: datasetKeys,
		TestOnly:    testOnly}
	dataBytes, _ := json.Marshal(data)

	// store data in the ledger
	var dataKeys string
	suffix := ", "
	for _, dataHash := range dataHashes {
		dataKeys = dataKeys + "\"" + dataHash + "\"" + suffix
		if err = stub.PutState(dataHash, dataBytes); err != nil {
			return nil, fmt.Errorf("failed to add data with hash %s", dataHash)
		}
		for _, datasetKey := range datasetKeys {
			// create composite keys to find all data associated with a dataset and only test or train data
			if err = createCompositeKey(stub, "data~dataset~key", []string{"data", datasetKey, dataHash}); err != nil {
				return nil, err
			}
		}
	}
	for _, dataHash := range dataHashes {
		for _, datasetKey := range datasetKeys {
			if err = createCompositeKey(stub, "data~dataset~testOnly~key", []string{"data", datasetKey, strconv.FormatBool(testOnly), dataHash}); err != nil {
				return nil, err
			}
		}
	}
	// return added data keys
	dataKeys = "{\"keys\": [" + strings.TrimSuffix(dataKeys, suffix) + "]}"
	return []byte(dataKeys), nil
}

// updateData associates one or more datasetKeys to one or more data
func updateData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	inp := inputUpdateData{}
	expectedArgs := getFieldNames(&inp)
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputUpdateData
	stringToInputStruct(args, &inp)
	// check validity of input args
	dataHashes, datasetKeys, err := validateUpdateData(stub, inp)
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
		for _, datasetKey := range datasetKeys {
			if !stringInSlice(datasetKey, data.DatasetKeys) {
				data.DatasetKeys = append(data.DatasetKeys, datasetKey)
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

// updateDataset associates a challengeKey to an existing dataset
func updateDataset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	inp := inputUpdateDataset{}
	expectedArgs := getFieldNames(&inp)
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputData
	stringToInputStruct(args, &inp)
	validate := validator.New()
	if err := validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid update dataset inputs %s", err.Error())
		return nil, err
	}
	// update dataset.ChallengeKey
	if err := addChallengeDataset(stub, inp.DatasetKey, inp.ChallengeKey); err != nil {
		return nil, err
	}
	return []byte(inp.DatasetKey), nil
}

// queryDatasetData returns info about a dataset and all related data
func queryDatasetData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := [1]string{"dataset key"}
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs[:], ", "))
	}
	datasetKey := args[0]

	// get dataset info
	var mPayload map[string]interface{}
	if err := getElementStruct(stub, datasetKey, &mPayload); err != nil {
		return nil, err
	}
	mPayload["key"] = datasetKey
	// get related train data
	trainDataKeys, err := getDatasetData(stub, datasetKey, false)
	if err != nil {
		return nil, err
	}
	mPayload["trainDataKeys"] = trainDataKeys
	// get related test data
	testDataKeys, err := getDatasetData(stub, datasetKey, true)
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
// -------------------- Data / Dataset utils -----------------------
// -----------------------------------------------------------------

// check

// checkDatasetOwner checks if the transaction requester is the owner of dataset
// specified by their keys in a slice
func checkDatasetOwner(stub shim.ChaincodeStubInterface, datasetKeys []string) (err error) {
	// get transaction requester
	txCreator, err := getTxCreator(stub)
	if err != nil {
		return
	}
	for _, datasetKey := range datasetKeys {
		dataset := Dataset{}
		if err = getElementStruct(stub, datasetKey, &dataset); err != nil {
			err = fmt.Errorf("could not retrieve dataset with key %s - %s", datasetKey, err.Error())
			return
		}
		// check transaction requester is the dataset owner
		if txCreator != dataset.Owner {
			err = fmt.Errorf("%s is not the owner of the dataset %s", txCreator, datasetKey)
			return
		}
	}
	return
}

// checkSameDataset checks if data in a slice exist and are from the same dataset.
// If yes, returns two boolean indicating if data are testOnly and trainOnly
func checkSameDataset(stub shim.ChaincodeStubInterface, datasetKey string, dataKeys []string) (testOnly bool, trainOnly bool, err error) {
	for i, dataKey := range dataKeys {
		data := Data{}
		if err = getElementStruct(stub, dataKey, &data); err != nil {
			err = fmt.Errorf("could not retrieve data with key %s - %s", dataKey, err.Error())
			return
		}
		if i == 0 {
			testOnly = data.TestOnly
			trainOnly = !testOnly
			continue
		}
		if !stringInSlice(datasetKey, data.DatasetKeys) {
			err = fmt.Errorf("data do not belong to the same dataset")
			return
		}
		if !data.TestOnly && testOnly {
			testOnly = false
		} else if data.TestOnly && trainOnly {
			trainOnly = false
		}
	}
	return
}

// getDatasetData returns all data keys associated to a dataset
func getDatasetData(stub shim.ChaincodeStubInterface, datasetKey string, testOnly bool) ([]string, error) {
	indexName := "data~dataset~testOnly~key"
	attributes := []string{"data", datasetKey, strconv.FormatBool(testOnly)}
	dataKeys, err := getKeysFromComposite(stub, indexName, attributes)
	if err != nil {
		return nil, err
	}
	return dataKeys, nil
}

// getDatasetOwner returns the owner of a dataset given its key
func getDatasetOwner(stub shim.ChaincodeStubInterface, datasetKey string) (worker string, err error) {

	dataset := Dataset{}
	if err = getElementStruct(stub, datasetKey, &dataset); err != nil {
		err = fmt.Errorf("could not retrieve dataset with key %s - %s", datasetKey, err.Error())
		return
	}
	worker = dataset.Owner
	return
}
