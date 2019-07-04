package main

import (
	"chaincode/errors"
	"fmt"
	"strings"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Set is a method of the receiver Objective. It checks the validity of inputObjective and uses its fields to set the Objective.
// Returns the objectiveKey and the dataManagerKey associated to test dataSample
func (objective *Objective) Set(stub shim.ChaincodeStubInterface, inp inputObjective) (objectiveKey string, dataManagerKey string, err error) {
	dataManagerKey = strings.Split(inp.TestDataset, ":")[0]
	if dataManagerKey != "" {
		var testOnly bool
		dataSampleKeys := strings.Split(strings.Replace(strings.Split(inp.TestDataset, ":")[1], " ", "", -1), ",")
		testOnly, _, err = checkSameDataManager(stub, dataManagerKey, dataSampleKeys)
		if err != nil {
			err = errors.BadRequest(err, "invalid test dataSample")
			return
		} else if !testOnly {
			err = errors.BadRequest("test dataSample are not tagged as testOnly dataSample")
			return
		}
		objective.TestDataset = &Dataset{
			DataManagerKey: dataManagerKey,
			DataSampleKeys: dataSampleKeys,
		}
	} else {
		objective.TestDataset = nil
	}
	objective.AssetType = ObjectiveType
	objective.Name = inp.Name
	objective.DescriptionStorageAddress = inp.DescriptionStorageAddress
	objective.Metrics = &HashDressName{
		Name:           inp.MetricsName,
		Hash:           inp.MetricsHash,
		StorageAddress: inp.MetricsStorageAddress,
	}
	owner, err := getTxCreator(stub)
	if err != nil {
		return
	}
	objective.Owner = owner
	objective.Permissions = inp.Permissions
	objectiveKey = inp.DescriptionHash
	return
}

// -------------------------------------------------------------------------------------------
// Smart contract related to objectivess
// -------------------------------------------------------------------------------------------

// registerObjective stores a new objective in the ledger.
// If the key exists, it will override the value with the new one
func registerObjective(stub shim.ChaincodeStubInterface, args []string) (resp map[string]string, err error) {
	// convert input strings args to input struct inputObjective
	inp := inputObjective{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}

	// check validity of input args and convert it to Objective
	objective := Objective{}
	objectiveKey, dataManagerKey, err := objective.Set(stub, inp)
	if err != nil {
		return
	}
	// check objective is not already in ledger
	if elementBytes, _ := stub.GetState(objectiveKey); elementBytes != nil {
		err = errors.Conflict("this objective already exists (tkey: %s)", objectiveKey)
		return
	}
	// submit to ledger
	objectiveBytes, _ := json.Marshal(objective)
	if err = stub.PutState(objectiveKey, objectiveBytes); err != nil {
		err = errors.E(err, "failed to submit to ledger the objective with key %s", objectiveKey)
		return
	}
	// create composite key
	if err = createCompositeKey(stub, "objective~owner~key", []string{"objective", objective.Owner, objectiveKey}); err != nil {
		return
	}
	// add objective to dataManager
	err = addObjectiveDataManager(stub, dataManagerKey, objectiveKey)
	return map[string]string{"key": objectiveKey}, err
}

// queryObjective returns a objective of the ledger given its key
func queryObjective(stub shim.ChaincodeStubInterface, args []string) (out outputObjective, err error) {
	inp := inputHashe{}
	err = AssetFromJSON(args[0], &inp)
	if err != nil {
		return
	}
	var objective Objective
	if err = getElementStruct(stub, inp.Key, &objective); err != nil {
		return
	}
	if objective.AssetType != ObjectiveType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	out.Fill(inp.Key, objective)
	return
}

// queryObjectives returns all objectives of the ledger
func queryObjectives(stub shim.ChaincodeStubInterface, args []string) (outObjectives []outputObjective, err error) {
	outObjectives = []outputObjective{}
	if len(args) != 0 {
		err = fmt.Errorf("incorrect number of arguments, expecting nothing")
		return
	}
	var indexName = "objective~owner~key"
	elementsKeys, err := getKeysFromComposite(stub, indexName, []string{"objective"})
	if err != nil {
		err = fmt.Errorf("issue getting keys from composite key %s - %s", indexName, err.Error())
		return
	}
	for _, key := range elementsKeys {
		var objective Objective
		if err = getElementStruct(stub, key, &objective); err != nil {
			return
		}
		var out outputObjective
		out.Fill(key, objective)
		outObjectives = append(outObjectives, out)
	}
	return
}

// -------------------------------------------------------------------------------------------
// Utils for objectivess
// -------------------------------------------------------------------------------------------

// addObjectiveDataManager associates a objective to a dataManager, more precisely, it adds the objective key to the dataManager
func addObjectiveDataManager(stub shim.ChaincodeStubInterface, dataManagerKey string, objectiveKey string) error {
	dataManager := DataManager{}
	if err := getElementStruct(stub, dataManagerKey, &dataManager); err != nil {
		return nil
	}
	if dataManager.ObjectiveKey != "" {
		return errors.BadRequest("dataManager is already associated with a objective")
	}
	dataManager.ObjectiveKey = objectiveKey
	dataManagerBytes, _ := json.Marshal(dataManager)
	return stub.PutState(dataManagerKey, dataManagerBytes)
}
