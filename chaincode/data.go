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

package main

import (
	"chaincode/errors"
	"strconv"
	"strings"
)

// Set is a method of the receiver DataManager. It uses inputDataManager fields to set the DataManager
// Returns the dataManagerKey and associated objectiveKeys
func (dataManager *DataManager) Set(db *LedgerDB, inp inputDataManager) (string, string, error) {
	dataManagerKey := inp.OpenerHash
	dataManager.ObjectiveKey = inp.ObjectiveKey
	dataManager.AssetType = DataManagerType
	dataManager.Name = inp.Name
	dataManager.OpenerStorageAddress = inp.OpenerStorageAddress
	dataManager.Type = inp.Type
	dataManager.Metadata = initMapOutput(inp.Metadata)
	dataManager.Description = &HashDress{
		Hash:           inp.DescriptionHash,
		StorageAddress: inp.DescriptionStorageAddress,
	}
	owner, err := GetTxCreator(db.cc)
	if err != nil {
		return "", "", err
	}
	dataManager.Owner = owner

	permissions, err := NewPermissions(db, inp.Permissions)
	if err != nil {
		return "", "", err
	}

	dataManager.Permissions = permissions
	return dataManagerKey, dataManager.ObjectiveKey, nil
}

// setDataSample is a method checking the validity of inputDataSample to be registered in the ledger
// and returning corresponding dataSample hashes, associated dataManagers, testOnly and errors
func setDataSample(db *LedgerDB, inp inputDataSample) (dataSampleHashes []string, dataSample DataSample, err error) {
	dataSampleHashes = inp.Hashes
	// check dataSample is not already in the ledger
	if existingKeys := checkDataSamplesExist(db, dataSampleHashes); existingKeys != nil {
		err = errors.Conflict("data samples with keys %s already exist", existingKeys).WithKeys(existingKeys)
		return
	}

	// get transaction owner
	owner, err := GetTxCreator(db.cc)
	if err != nil {
		return
	}
	// check if associated dataManager(s) exists
	var dataManagerKeys []string
	if len(inp.DataManagerKeys) > 0 {
		dataManagerKeys = inp.DataManagerKeys
		if err = checkDataManagerOwner(db, dataManagerKeys); err != nil {
			return
		}
	}
	// convert input testOnly to boolean
	testOnly, err := strconv.ParseBool(inp.TestOnly)

	dataSample = DataSample{
		AssetType:       DataSampleType,
		DataManagerKeys: dataManagerKeys,
		TestOnly:        testOnly,
		Owner:           owner}

	return
}

// validateUpdateDataSample is a method checking the validity of elements sent to update
// one or more dataSamplef
func validateUpdateDataSample(db *LedgerDB, inp inputUpdateDataSample) (dataSampleHashes []string, dataManagerKeys []string, err error) {
	// TODO return full dataSample
	// check dataManagers exist and are owned by the transaction requester
	if err = checkDataManagerOwner(db, inp.DataManagerKeys); err != nil {
		return
	}
	return inp.Hashes, inp.DataManagerKeys, nil
}

// -----------------------------------------------------------------
// ----------------------- Smart Contracts  ------------------------
// -----------------------------------------------------------------

// registerDataManager stores a new dataManager in the ledger.
func registerDataManager(db *LedgerDB, args []string) (resp outputKey, err error) {
	inp := inputDataManager{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	// check validity of input args and convert it to a DataManager
	if len(inp.ObjectiveKey) > 0 {
		if _, err := db.GetObjective(inp.ObjectiveKey); err != nil {
			err = errors.BadRequest(err, "error checking associated objective")
			return resp, err
		}
	}
	dataManager := DataManager{}
	dataManagerKey, objectiveKey, err := dataManager.Set(db, inp)
	if err != nil {
		return
	}
	// submit to ledger
	err = db.Add(dataManagerKey, dataManager)
	if err != nil {
		return
	}
	// create composite keys (one for each associated objective) to find dataSample associated with a objective
	indexName := "dataManager~objective~key"
	err = db.CreateIndex(indexName, []string{"dataManager", objectiveKey, dataManagerKey})
	if err != nil {
		return
	}
	// create composite key to find dataManager associated with a owner
	err = db.CreateIndex("dataManager~owner~key", []string{"dataManager", dataManager.Owner, dataManagerKey})
	if err != nil {
		return
	}
	return outputKey{Key: dataManagerKey}, nil
}

// registerDataSample stores new dataSample in the ledger (one or more).
func registerDataSample(db *LedgerDB, args []string) (dataSampleKeys map[string][]string, err error) {
	// convert input strings args to input struct inputDataSample
	inp := inputDataSample{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	// check validity of input args
	dataSampleHashes, dataSample, err := setDataSample(db, inp)
	if err != nil {
		return
	}

	// store dataSample in the ledger
	for _, dataSampleHash := range dataSampleHashes {
		if err = db.Add(dataSampleHash, dataSample); err != nil {
			return
		}
		for _, dataManagerKey := range dataSample.DataManagerKeys {
			// create composite keys to find all dataSample associated with a dataManager and both test and train dataSample
			if err = db.CreateIndex("dataSample~dataManager~key", []string{"dataSample", dataManagerKey, dataSampleHash}); err != nil {
				return
			}
			// create composite keys to find all dataSample associated with a dataManager and only test or train dataSample
			if err = db.CreateIndex("dataSample~dataManager~testOnly~key", []string{"dataSample", dataManagerKey, strconv.FormatBool(dataSample.TestOnly), dataSampleHash}); err != nil {
				return
			}
		}
	}
	// return added dataSample keys
	dataSampleKeys = map[string][]string{"keys": dataSampleHashes}
	return
}

// updateDataSample associates one or more dataManagerKeys to one or more dataSample
func updateDataSample(db *LedgerDB, args []string) (resp outputKey, err error) {
	inp := inputUpdateDataSample{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	// check validity of input args
	dataSampleHashes, dataManagerKeys, err := validateUpdateDataSample(db, inp)
	if err != nil {
		return
	}
	// store dataSample in the ledger
	var dataSampleKeys string
	suffix := ", "
	for _, dataSampleHash := range dataSampleHashes {
		dataSampleKeys = dataSampleKeys + "\"" + dataSampleHash + "\"" + suffix
		var dataSample DataSample
		dataSample, err = db.GetDataSample(dataSampleHash)
		if err != nil {
			return
		}
		if err = checkDataSampleOwner(db, dataSample); err != nil {
			return
		}
		for _, dataManagerKey := range dataManagerKeys {
			if !stringInSlice(dataManagerKey, dataSample.DataManagerKeys) {
				// check data manager is not already associated with this data
				dataSample.DataManagerKeys = append(dataSample.DataManagerKeys, dataManagerKey)
				// create composite keys to find all dataSample associated with a dataManager and both test and train dataSample
				if err = db.CreateIndex("dataSample~dataManager~key", []string{"dataSample", dataManagerKey, dataSampleHash}); err != nil {
					return
				}
				// create composite keys to find all dataSample associated with a dataManager and only test or train dataSample
				if err = db.CreateIndex("dataSample~dataManager~testOnly~key", []string{"dataSample", dataManagerKey, strconv.FormatBool(dataSample.TestOnly), dataSampleHash}); err != nil {
					return
				}
			}
		}
		if err = db.Put(dataSampleHash, dataSample); err != nil {
			return
		}

	}
	// return updated dataSample keys
	// TODO return a json struct
	dataSampleKeys = "{\"keys\": [" + strings.TrimSuffix(dataSampleKeys, suffix) + "]}"
	return outputKey{Key: dataSampleKeys}, nil
}

// updateDataManager associates a objectiveKey to an existing dataManager
func updateDataManager(db *LedgerDB, args []string) (resp outputKey, err error) {
	inp := inputUpdateDataManager{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// update dataManager.ObjectiveKey
	if err = addObjectiveDataManager(db, inp.DataManagerKey, inp.ObjectiveKey); err != nil {
		return
	}
	return outputKey{Key: inp.DataManagerKey}, nil
}

// queryDataManager returns dataManager and its key
func queryDataManager(db *LedgerDB, args []string) (out outputDataManager, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	dataManager, err := db.GetDataManager(inp.Key)
	if err != nil {
		return
	}
	if dataManager.AssetType != DataManagerType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	out.Fill(inp.Key, dataManager)
	return
}

// queryDataManagers returns all DataManagers of the ledger
func queryDataManagers(db *LedgerDB, args []string) ([]outputDataManager, error) {
	var err error
	outDataManagers := []outputDataManager{}
	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outDataManagers, err
	}
	var indexName = "dataManager~owner~key"
	elementsKeys, err := db.GetIndexKeys(indexName, []string{"dataManager"})
	if err != nil {
		return outDataManagers, err
	}
	for _, key := range elementsKeys {
		dataManager, err := db.GetDataManager(key)
		if err != nil {
			return outDataManagers, err
		}
		var out outputDataManager
		out.Fill(key, dataManager)
		outDataManagers = append(outDataManagers, out)
	}
	return outDataManagers, nil
}

// queryDataset returns info about a dataManager and all related dataSample
func queryDataset(db *LedgerDB, args []string) (outputDataset, error) {
	inp := inputKey{}
	out := outputDataset{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return out, err
	}

	dataManager, err := db.GetDataManager(inp.Key)
	if err != nil {
		return out, err
	}

	// get related train dataSample
	trainDataSampleKeys, err := getDataset(db, inp.Key, false)
	if err != nil {
		return out, err
	}

	// get related test dataSample
	testDataSampleKeys, err := getDataset(db, inp.Key, true)
	if err != nil {
		return out, err
	}

	out.Fill(inp.Key, dataManager, trainDataSampleKeys, testDataSampleKeys)
	return out, nil
}

func queryDataSamples(db *LedgerDB, args []string) ([]outputDataSample, error) {
	outDataSamples := []outputDataSample{}
	if len(args) != 0 {
		err := errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outDataSamples, err
	}
	elementsKeys, err := db.GetIndexKeys("dataSample~dataManager~key", []string{"dataSample"})
	if err != nil {
		return outDataSamples, err
	}
	for _, key := range elementsKeys {
		var dataSample DataSample
		dataSample, err = db.GetDataSample(key)
		if err != nil {
			return outDataSamples, err
		}
		var out outputDataSample
		out.Fill(key, dataSample)
		outDataSamples = append(outDataSamples, out)
	}
	return outDataSamples, nil
}

// -----------------------------------------------------------------
// -------------------- DataSample / DataManager utils -----------------------
// -----------------------------------------------------------------

// check

// checkDataManagerOwner checks if the transaction requester is the owner of dataManager
// specified by their keys in a slice
func checkDataManagerOwner(db *LedgerDB, dataManagerKeys []string) error {
	// get transaction requester
	txCreator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	for _, dataManagerKey := range dataManagerKeys {
		dataManager, err := db.GetDataManager(dataManagerKey)
		if err != nil {
			return errors.BadRequest(err, "could not retrieve dataManager with key %s", dataManagerKey)
		}
		// check transaction requester is the dataManager owner
		if txCreator != dataManager.Owner {
			return errors.Forbidden("%s is not the owner of the dataManager %s", txCreator, dataManagerKey)
		}
	}
	return nil
}

//  checkDataSampleOwner checks if the transaction requester is the owner of the dataSample
func checkDataSampleOwner(db *LedgerDB, dataSample DataSample) error {
	txRequester, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	if txRequester != dataSample.Owner {
		return errors.Forbidden("%s is not the dataSample's owner", txRequester)
	}
	return nil
}

// checkSameDataManager checks if dataSample in a slice exist and are from the same dataManager.
// If yes, returns two boolean indicating if dataSample are testOnly and trainOnly
func checkSameDataManager(db *LedgerDB, dataManagerKey string, dataSampleKeys []string) (bool, bool, error) {
	testOnly := true
	trainOnly := true
	for _, dataSampleKey := range dataSampleKeys {
		dataSample, err := db.GetDataSample(dataSampleKey)
		if err != nil {
			return testOnly, trainOnly, err
		}
		if !stringInSlice(dataManagerKey, dataSample.DataManagerKeys) {
			err = errors.BadRequest("dataSample do not belong to the same dataManager")
			return testOnly, trainOnly, err
		}
		testOnly = testOnly && dataSample.TestOnly
		trainOnly = trainOnly && !dataSample.TestOnly
	}
	return testOnly, trainOnly, nil
}

// getDataset returns all dataSample keys associated to a dataManager
func getDataset(db *LedgerDB, dataManagerKey string, testOnly bool) ([]string, error) {
	indexName := "dataSample~dataManager~testOnly~key"
	attributes := []string{"dataSample", dataManagerKey, strconv.FormatBool(testOnly)}
	dataSampleKeys, err := db.GetIndexKeys(indexName, attributes)
	if err != nil {
		return nil, err
	}
	return dataSampleKeys, nil
}

// getDataManagerOwner returns the owner of a dataManager given its key
func getDataManagerOwner(db *LedgerDB, dataManagerKey string) (string, error) {
	dataManager, err := db.GetDataManager(dataManagerKey)
	if err != nil {
		return "", errors.BadRequest(err, "dataManager %s not found", dataManagerKey)
	}
	return dataManager.Owner, nil
}

// checkDataSamplesExist checks if keys in a slice correspond to existing elements in the ledger
// returns the slice of already existing elements
func checkDataSamplesExist(db *LedgerDB, keys []string) (existingKeys []string) {
	for _, key := range keys {
		if _, err := db.GetDataSample(key); err == nil {
			existingKeys = append(existingKeys, key)
		}
	}
	return
}
