package main

import (
	"chaincode/errors"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"encoding/json"
)

// List of the possible tuple's status
const (
	StatusDoing   = "doing"
	StatusTodo    = "todo"
	StatusWaiting = "waiting"
	StatusFailed  = "failed"
	StatusDone    = "done"
)

// -------------------------------------------------------------------------------------------
// Methods on receivers traintuple
// -------------------------------------------------------------------------------------------

// SetFromInput is a method of the receiver Traintuple.
// It uses the inputTraintuple to check and set the traintuple's parameters
// which don't depend on previous traintuples values :
//  - AssetType
//  - Creator & permissions
//  - Tag
//  - AlgoKey & ObjectiveKey
//  - Dataset
func (traintuple *Traintuple) SetFromInput(db LedgerDB, inp inputTraintuple) error {

	// TODO later: check permissions
	// find associated creator and check permissions (TODO later)
	creator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	traintuple.AssetType = TraintupleType
	traintuple.Creator = creator
	traintuple.Tag = inp.Tag
	algo, err := db.GetAlgo(inp.AlgoKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve algo with key %s", inp.AlgoKey)
	}
	if !algo.Permissions.CanProcess(algo.Owner, creator) {
		return errors.Forbidden("not authorized to process this algo %s", inp.AlgoKey)
	}
	traintuple.AlgoKey = inp.AlgoKey

	// check objective exists
	objective, err := db.GetObjective(inp.ObjectiveKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve objective with key %s", inp.ObjectiveKey)
	}
	if !objective.Permissions.CanProcess(objective.Owner, creator) {
		return errors.Forbidden("not authorized to process this objective %s", inp.ObjectiveKey)
	}
	traintuple.ObjectiveKey = inp.ObjectiveKey

	// check if DataSampleKeys are from the same dataManager and if they are not test only dataSample
	_, trainOnly, err := checkSameDataManager(db, inp.DataManagerKey, inp.DataSampleKeys)
	if err != nil {
		return err
	}
	if !trainOnly {
		return errors.BadRequest("not possible to create a traintuple with test only data")
	}

	dataManager, err := db.GetDataManager(inp.DataManagerKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve dataManager with key %s", inp.DataManagerKey)
	}
	if !dataManager.Permissions.CanProcess(dataManager.Owner, creator) {
		return errors.Forbidden("not authorized to process this dataManager %s", inp.DataManagerKey)
	}

	traintuple.Permissions = MergePermissions(dataManager.Permissions, algo.Permissions)

	// fill traintuple.Dataset from dataManager and dataSample
	traintuple.Dataset = &Dataset{
		DataManagerKey: inp.DataManagerKey,
		DataSampleKeys: inp.DataSampleKeys,
	}
	traintuple.Dataset.Worker, err = getDataManagerOwner(db, traintuple.Dataset.DataManagerKey)
	return err
}

// SetFromParents set the status of the traintuple depending on its "parents",
// i.e. the traintuples from which it received the outModels as inModels.
// Also it's InModelKeys are set.
func (traintuple *Traintuple) SetFromParents(db LedgerDB, inModels []string) error {
	status := StatusTodo
	parentTraintupleKeys := inModels
	for _, parentTraintupleKey := range parentTraintupleKeys {
		parentTraintuple, err := db.GetTraintuple(parentTraintupleKey)
		if err != nil {
			err = errors.BadRequest(err, "could not retrieve parent traintuple with key %s %d", parentTraintupleKeys, len(parentTraintupleKeys))
			return err
		}
		// set traintuple to waiting if one of the parent traintuples is not done
		if parentTraintuple.OutModel == nil {
			status = StatusWaiting
		}
		traintuple.InModelKeys = append(traintuple.InModelKeys, parentTraintupleKey)
	}
	traintuple.Status = status
	return nil
}

// GetKey return the key of the traintuple depending on its key parameters.
func (traintuple *Traintuple) GetKey() string {
	hashKeys := []string{traintuple.Creator, traintuple.AlgoKey, traintuple.Dataset.DataManagerKey}
	hashKeys = append(hashKeys, traintuple.Dataset.DataSampleKeys...)
	hashKeys = append(hashKeys, traintuple.InModelKeys...)
	return HashForKey("traintuple", hashKeys...)

}

// AddToComputePlan set the traintuple's parameters that determines if it's part of on ComputePlan and how.
// It uses the inputTraintuple values as follow:
//  - If neither ComputePlanID nor rank is set it returns immediately
//  - If rank is 0 and ComputePlanID empty, it's start a new one using this traintuple key
//  - If rank and ComputePlanID are set, it checks if there are coherent with previous ones and set it.
func (traintuple *Traintuple) AddToComputePlan(db LedgerDB, inp inputTraintuple, traintupleKey string) error {
	// check ComputePlanID and Rank and set it when required
	var err error
	if inp.Rank == "" {
		if inp.ComputePlanID != "" {
			return errors.BadRequest("invalid inputs, a ComputePlan should have a rank")
		}
		return nil
	}
	traintuple.Rank, err = strconv.Atoi(inp.Rank)
	if err != nil {
		return err
	}
	if inp.ComputePlanID == "" {
		if traintuple.Rank != 0 {
			err = errors.BadRequest("invalid inputs, a new ComputePlan should have a rank 0")
			return err
		}
		traintuple.ComputePlanID = traintupleKey
		return nil
	}
	var ttKeys []string
	ttKeys, err = db.GetIndexKeys("traintuple~computeplanid~worker~rank~key", []string{"traintuple", inp.ComputePlanID})
	if err != nil {
		return err
	}
	if len(ttKeys) == 0 {
		return errors.BadRequest("cannot find the ComputePlanID %s", inp.ComputePlanID)
	}
	for _, ttKey := range ttKeys {
		FLTraintuple, err := db.GetTraintuple(ttKey)
		if err != nil {
			return err
		}
		if FLTraintuple.AlgoKey != inp.AlgoKey {
			return errors.BadRequest("previous traintuple for ComputePlanID %s does not have the same algo key %s", inp.ComputePlanID, inp.AlgoKey)
		}
	}

	ttKeys, err = db.GetIndexKeys("traintuple~computeplanid~worker~rank~key", []string{"traintuple", inp.ComputePlanID, traintuple.Dataset.Worker, inp.Rank})
	if err != nil {
		return err
	} else if len(ttKeys) > 0 {
		err = errors.BadRequest("ComputePlanID %s with worker %s rank %d already exists", inp.ComputePlanID, traintuple.Dataset.Worker, traintuple.Rank)
		return err
	}

	traintuple.ComputePlanID = inp.ComputePlanID

	return nil
}

// Save will put in the legder interface both the traintuple with its key
// and all the associated composite keys
func (traintuple *Traintuple) Save(db LedgerDB, traintupleKey string) error {

	// store in ledger
	if err := db.Add(traintupleKey, traintuple); err != nil {
		return err
	}

	// create composite keys
	if err := db.CreateIndex("traintuple~algo~key", []string{"traintuple", traintuple.AlgoKey, traintupleKey}); err != nil {
		return err
	}
	if err := db.CreateIndex("traintuple~worker~status~key", []string{"traintuple", traintuple.Dataset.Worker, traintuple.Status, traintupleKey}); err != nil {
		return err
	}
	for _, inModelKey := range traintuple.InModelKeys {
		if err := db.CreateIndex("traintuple~inModel~key", []string{"traintuple", inModelKey, traintupleKey}); err != nil {
			return err
		}
	}
	if traintuple.ComputePlanID != "" {
		if err := db.CreateIndex("traintuple~computeplanid~worker~rank~key", []string{"traintuple", traintuple.ComputePlanID, traintuple.Dataset.Worker, strconv.Itoa(traintuple.Rank), traintupleKey}); err != nil {
			return err
		}
	}
	if traintuple.Tag != "" {
		err := db.CreateIndex("traintuple~tag~key", []string{"traintuple", traintuple.Tag, traintupleKey})
		if err != nil {
			return err
		}
	}
	return nil
}

// -------------------------------------------------------------------------------------------
// Methods on receivers testtuple
// -------------------------------------------------------------------------------------------

// SetFromInput is a method of the receiver Testtuple.
// It uses the inputTesttuple to check and set the testtuple's parameters
// which don't depend on previous testtuples values :
//  - AssetType
//  - Creator & permissions
//  - Tag
//  - AlgoKey & ObjectiveKey
//  - Dataset
//  - Certified
func (testtuple *Testtuple) SetFromInput(db LedgerDB, inp inputTesttuple) error {

	// TODO later: check permissions
	// find associated creator and check permissions (TODO later)
	creator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	testtuple.Creator = creator
	testtuple.Permissions = "all"
	testtuple.Tag = inp.Tag
	testtuple.AssetType = TesttupleType

	// Get test dataset from objective
	objective, err := db.GetObjective(testtuple.ObjectiveKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve objective with key %s", testtuple.ObjectiveKey)
	}
	var objectiveDataManagerKey string
	var objectiveDataSampleKeys []string
	if objective.TestDataset != nil {
		objectiveDataManagerKey = objective.TestDataset.DataManagerKey
		objectiveDataSampleKeys = objective.TestDataset.DataSampleKeys
	}
	// For now we need to sort it but in fine it should be save sorted
	// TODO
	sort.Strings(objectiveDataSampleKeys)

	var dataManagerKey string
	var dataSampleKeys []string
	if len(inp.DataManagerKey) > 0 && len(inp.DataSampleKeys) > 0 {
		// non-certified testtuple
		// test dataset are specified by the user
		dataSampleKeys = inp.DataSampleKeys
		_, _, err = checkSameDataManager(db, inp.DataManagerKey, dataSampleKeys)
		if err != nil {
			return err
		}
		dataManagerKey = inp.DataManagerKey
		sort.Strings(dataSampleKeys)
		testtuple.Certified = objectiveDataManagerKey == dataManagerKey && reflect.DeepEqual(objectiveDataSampleKeys, dataSampleKeys)
	} else if len(inp.DataManagerKey) > 0 || len(inp.DataSampleKeys) > 0 {
		return errors.BadRequest("invalid input: dataManagerKey and dataSampleKey should be provided together")
	} else if objective.TestDataset != nil {
		dataSampleKeys = objectiveDataSampleKeys
		dataManagerKey = objectiveDataManagerKey
		testtuple.Certified = true
	} else {
		return errors.BadRequest("can not create a certified testtuple, no data associated with objective %s", testtuple.ObjectiveKey)
	}
	// retrieve dataManager owner
	dataManager, err := db.GetDataManager(dataManagerKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve dataManager with key %s", dataManagerKey)
	}
	testtuple.Dataset = &TtDataset{
		Worker:         dataManager.Owner,
		DataSampleKeys: dataSampleKeys,
		OpenerHash:     dataManagerKey,
	}
	return nil
}

// SetFromTraintuple set the parameters of the testuple depending on traintuple
// it depends on. It sets:
//  - AlgoKey
//  - ObjectiveKey
//  - Model
//  - Status
func (testtuple *Testtuple) SetFromTraintuple(db LedgerDB, traintupleKey string) error {

	// check associated traintuple
	traintuple, err := db.GetTraintuple(traintupleKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve traintuple with key %s", traintupleKey)
	}
	testtuple.ObjectiveKey = traintuple.ObjectiveKey
	testtuple.AlgoKey = traintuple.AlgoKey
	testtuple.Model = &Model{
		TraintupleKey: traintupleKey,
	}
	if traintuple.OutModel != nil {
		testtuple.Model.Hash = traintuple.OutModel.Hash
		testtuple.Model.StorageAddress = traintuple.OutModel.StorageAddress
	}

	switch status := traintuple.Status; status {
	case StatusDone:
		testtuple.Status = StatusTodo
	case StatusFailed:
		return errors.BadRequest(
			"could not register this testtuple, the traintuple %s has a failed status",
			traintupleKey)
	default:
		testtuple.Status = StatusWaiting
	}
	return nil
}

// GetKey return the key of the testuple depending on its key parameters.
func (testtuple *Testtuple) GetKey() string {
	// create testtuple key and check if it already exists
	hashKeys := []string{
		testtuple.Model.TraintupleKey,
		testtuple.Dataset.OpenerHash,
		testtuple.Creator,
	}
	hashKeys = append(hashKeys, testtuple.Dataset.DataSampleKeys...)
	return HashForKey("testtuple", hashKeys...)
}

// Save will put in the legder interface both the testtuple with its key
// and all the associated composite keys
func (testtuple *Testtuple) Save(db LedgerDB, testtupleKey string) error {
	var err error
	if err = db.Add(testtupleKey, testtuple); err != nil {
		return err
	}

	// create composite keys
	if err = db.CreateIndex("testtuple~objective~certified~key", []string{"testtuple", testtuple.ObjectiveKey, strconv.FormatBool(testtuple.Certified), testtupleKey}); err != nil {
		return err
	}
	if err = db.CreateIndex("testtuple~algo~key", []string{"testtuple", testtuple.AlgoKey, testtupleKey}); err != nil {
		return err
	}
	if err = db.CreateIndex("testtuple~worker~status~key", []string{"testtuple", testtuple.Dataset.Worker, testtuple.Status, testtupleKey}); err != nil {
		return err
	}
	if err = db.CreateIndex("testtuple~traintuple~certified~key", []string{"testtuple", testtuple.Model.TraintupleKey, strconv.FormatBool(testtuple.Certified), testtupleKey}); err != nil {
		return err
	}
	if testtuple.Tag != "" {
		err = db.CreateIndex("testtuple~tag~key", []string{"traintuple", testtuple.Tag, testtupleKey})
		if err != nil {
			return err
		}
	}
	return nil
}

// -------------------------------------------------------------------------------------------
// Smart contracts related to traintuples and testuples
// args  [][]byte or []string, it is not possible to input a string looking like a json
// -------------------------------------------------------------------------------------------

// createComputePlan is the wrapper for the substra smartcontract CreateComputePlan
func createComputePlan(db LedgerDB, args []string) (resp outputComputePlan, err error) {
	inp := inputComputePlan{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintupleKeysByID := map[string]string{}
	resp.TraintupleKeys = []string{}
	var traintuplesTodo []outputTraintuple
	for i, computeTraintuple := range inp.Traintuples {
		inpTraintuple := inputTraintuple{}
		inpTraintuple.AlgoKey = inp.AlgoKey
		inpTraintuple.ObjectiveKey = inp.ObjectiveKey
		inpTraintuple.DataManagerKey = computeTraintuple.DataManagerKey
		inpTraintuple.DataSampleKeys = computeTraintuple.DataSampleKeys
		inpTraintuple.Tag = computeTraintuple.Tag
		inpTraintuple.Rank = strconv.Itoa(i)

		traintuple := Traintuple{}
		err := traintuple.SetFromInput(db, inpTraintuple)
		if err != nil {
			return resp, err
		}

		// Set the inModels by matching the id to traintuples key previously
		// encontered in this compute plan
		for _, InModelID := range computeTraintuple.InModelsIDs {
			inModelKey, ok := traintupleKeysByID[InModelID]
			if !ok {
				return resp, errors.BadRequest("traintuple ID %s: model ID %s not found, check traintuple list order", computeTraintuple.ID, InModelID)
			}
			traintuple.InModelKeys = append(traintuple.InModelKeys, inModelKey)
		}

		traintupleKey := traintuple.GetKey()

		// Set the ComputePlanID
		if i == 0 {
			traintuple.ComputePlanID = traintupleKey
			resp.ComputePlanID = traintuple.ComputePlanID
		} else {
			traintuple.ComputePlanID = resp.ComputePlanID
		}

		// Set status: if it has parents it's waiting
		// if not it's todo and it has to be included in the event
		if len(computeTraintuple.InModelsIDs) > 0 {
			traintuple.Status = StatusWaiting
		} else {
			traintuple.Status = StatusTodo
			out := outputTraintuple{}
			err = out.Fill(db, traintuple, traintupleKey)
			if err != nil {
				return resp, err
			}
			traintuplesTodo = append(traintuplesTodo, out)
		}

		err = traintuple.Save(db, traintupleKey)
		if err != nil {
			return resp, err
		}
		traintupleKeysByID[computeTraintuple.ID] = traintupleKey
		resp.TraintupleKeys = append(resp.TraintupleKeys, traintupleKey)
	}

	resp.TesttupleKeys = []string{}
	for index, computeTesttuple := range inp.Testtuples {
		traintupleKey, ok := traintupleKeysByID[computeTesttuple.TraintupleID]
		if !ok {
			return resp, errors.BadRequest("testtuple index %s: traintuple ID %s not found", index, computeTesttuple.TraintupleID)
		}
		testtuple := Testtuple{}
		testtuple.Model = &Model{TraintupleKey: traintupleKey}
		testtuple.ObjectiveKey = inp.ObjectiveKey
		testtuple.AlgoKey = inp.AlgoKey

		inputTesttuple := inputTesttuple{}
		inputTesttuple.DataManagerKey = computeTesttuple.DataManagerKey
		inputTesttuple.DataSampleKeys = computeTesttuple.DataSampleKeys
		inputTesttuple.Tag = computeTesttuple.Tag
		err = testtuple.SetFromInput(db, inputTesttuple)
		if err != nil {
			return resp, err
		}
		testtuple.Status = StatusWaiting
		testtupleKey := testtuple.GetKey()
		err = testtuple.Save(db, testtupleKey)
		if err != nil {
			return resp, err
		}
		resp.TesttupleKeys = append(resp.TesttupleKeys, testtupleKey)
	}

	event := TuplesEvent{}
	event.SetTraintuples(traintuplesTodo...)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// createTraintuple adds a Traintuple in the ledger
func createTraintuple(db LedgerDB, args []string) (map[string]string, error) {
	inp := inputTraintuple{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return nil, err
	}

	traintuple := Traintuple{}
	err = traintuple.SetFromInput(db, inp)
	if err != nil {
		return nil, err
	}
	err = traintuple.SetFromParents(db, inp.InModels)
	if err != nil {
		return nil, err
	}
	traintupleKey := traintuple.GetKey()
	err = traintuple.AddToComputePlan(db, inp, traintupleKey)
	if err != nil {
		return nil, err
	}
	err = traintuple.Save(db, traintupleKey)
	if err != nil {
		return nil, err
	}
	out := outputTraintuple{}
	err = out.Fill(db, traintuple, traintupleKey)
	if err != nil {
		return nil, err
	}

	event := TuplesEvent{}
	event.SetTraintuples(out)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return nil, err
	}

	return map[string]string{"key": traintupleKey}, nil
}

// createTesttuple adds a Testtuple in the ledger
func createTesttuple(db LedgerDB, args []string) (map[string]string, error) {
	inp := inputTesttuple{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return nil, err
	}

	// check validity of input arg and set testtuple
	testtuple := Testtuple{}
	err = testtuple.SetFromTraintuple(db, inp.TraintupleKey)
	if err != nil {
		return nil, err
	}
	err = testtuple.SetFromInput(db, inp)
	if err != nil {
		return nil, err
	}
	testtupleKey := testtuple.GetKey()
	err = testtuple.Save(db, testtupleKey)
	if err != nil {
		return nil, err
	}
	out := outputTesttuple{}
	err = out.Fill(db, testtupleKey, testtuple)
	if err != nil {
		return nil, err
	}

	event := TuplesEvent{}
	event.SetTesttuples(out)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return nil, err
	}

	return map[string]string{"key": testtupleKey}, nil
}

// logStartTrain modifies a traintuple by changing its status from todo to doing
func logStartTrain(db LedgerDB, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get traintuple, check validity of the update
	traintuple, err := db.GetTraintuple(inp.Key)
	if err != nil {
		return
	}
	if err = validateTupleOwner(db, traintuple.Dataset.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, inp.Key, StatusDoing); err != nil {
		return
	}
	outputTraintuple.Fill(db, traintuple, inp.Key)
	return
}

// logStartTest modifies a testtuple by changing its status from todo to doing
func logStartTest(db LedgerDB, args []string) (outputTesttuple outputTesttuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get testtuple, check validity of the update, and update its status
	testtuple, err := db.GetTesttuple(inp.Key)
	if err != nil {
		return
	}
	if err = validateTupleOwner(db, testtuple.Dataset.Worker); err != nil {
		return
	}
	if err = testtuple.commitStatusUpdate(db, inp.Key, StatusDoing); err != nil {
		return
	}
	err = outputTesttuple.Fill(db, inp.Key, testtuple)
	if err != nil {
		return
	}
	return
}

// logSuccessTrain modifies a traintuple by changing its status from doing to done
// reports logs and associated performances
func logSuccessTrain(db LedgerDB, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputLogSuccessTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintupleKey := inp.Key

	// get, update and commit traintuple
	traintuple, err := db.GetTraintuple(traintupleKey)
	if err != nil {
		return
	}
	traintuple.Perf = inp.Perf
	traintuple.OutModel = &HashDress{
		Hash:           inp.OutModel.Hash,
		StorageAddress: inp.OutModel.StorageAddress}
	traintuple.Log += inp.Log

	if err = validateTupleOwner(db, traintuple.Dataset.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, traintupleKey, StatusDone); err != nil {
		return
	}

	// update depending tuples
	traintuplesEvent, err := traintuple.updateTraintupleChildren(db, traintupleKey)
	if err != nil {
		return
	}

	testtuplesEvent, err := traintuple.updateTesttupleChildren(db, traintupleKey)
	if err != nil {
		return
	}

	outputTraintuple.Fill(db, traintuple, inp.Key)

	event := TuplesEvent{}
	event.SetTraintuples(traintuplesEvent...)
	event.SetTesttuples(testtuplesEvent...)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return
	}

	return
}

// logSuccessTest modifies a testtuple by changing its status to done, reports perf and logs
func logSuccessTest(db LedgerDB, args []string) (outputTesttuple outputTesttuple, err error) {
	inp := inputLogSuccessTest{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	testtuple, err := db.GetTesttuple(inp.Key)
	if err != nil {
		return
	}

	testtuple.Dataset.Perf = inp.Perf
	testtuple.Log += inp.Log

	if err = validateTupleOwner(db, testtuple.Dataset.Worker); err != nil {
		return
	}
	if err = testtuple.commitStatusUpdate(db, inp.Key, StatusDone); err != nil {
		return
	}
	err = outputTesttuple.Fill(db, inp.Key, testtuple)
	return
}

// logFailTrain modifies a traintuple by changing its status to fail and reports associated logs
func logFailTrain(db LedgerDB, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputLogFailTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get, update and commit traintuple
	traintuple, err := db.GetTraintuple(inp.Key)
	if err != nil {
		return
	}
	traintuple.Log += inp.Log

	if err = validateTupleOwner(db, traintuple.Dataset.Worker); err != nil {
		return
	}
	if err = traintuple.commitStatusUpdate(db, inp.Key, StatusFailed); err != nil {
		return
	}

	outputTraintuple.Fill(db, traintuple, inp.Key)

	// update depending tuples
	testtuplesEvent, err := traintuple.updateTesttupleChildren(db, inp.Key)
	if err != nil {
		return
	}

	traintuplesEvent, err := traintuple.updateTraintupleChildren(db, inp.Key)
	if err != nil {
		return
	}

	event := TuplesEvent{}
	event.SetTraintuples(traintuplesEvent...)
	event.SetTesttuples(testtuplesEvent...)
	err = SendTuplesEvent(db.cc, event)
	if err != nil {
		return
	}

	return
}

// logFailTest modifies a testtuple by changing its status to fail and reports associated logs
func logFailTest(db LedgerDB, args []string) (outputTesttuple outputTesttuple, err error) {
	inp := inputLogFailTest{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get, update and commit testtuple
	testtuple, err := db.GetTesttuple(inp.Key)
	if err != nil {
		return
	}

	testtuple.Log += inp.Log

	if err = validateTupleOwner(db, testtuple.Dataset.Worker); err != nil {
		return
	}
	if err = testtuple.commitStatusUpdate(db, inp.Key, StatusFailed); err != nil {
		return
	}
	err = outputTesttuple.Fill(db, inp.Key, testtuple)
	return
}

// queryTraintuple returns info about a traintuple given its key
func queryTraintuple(db LedgerDB, args []string) (outputTraintuple outputTraintuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintuple, err := db.GetTraintuple(inp.Key)
	if err != nil {
		return
	}
	if traintuple.AssetType != TraintupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	outputTraintuple.Fill(db, traintuple, inp.Key)
	return
}

// queryTraintuples returns all traintuples
func queryTraintuples(db LedgerDB, args []string) ([]outputTraintuple, error) {
	outTraintuples := []outputTraintuple{}

	if len(args) != 0 {
		err := errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outTraintuples, err
	}
	elementsKeys, err := db.GetIndexKeys("traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return outTraintuples, err
	}
	for _, key := range elementsKeys {
		outputTraintuple, err := getOutputTraintuple(db, key)
		if err != nil {
			return outTraintuples, err
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return outTraintuples, nil
}

// queryTesttuple returns a testtuple of the ledger given its key
func queryTesttuple(db LedgerDB, args []string) (out outputTesttuple, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	testtuple, err := db.GetTesttuple(inp.Key)
	if err != nil {
		return
	}
	if testtuple.AssetType != TesttupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	err = out.Fill(db, inp.Key, testtuple)
	return
}

// queryTesttuples returns all testtuples of the ledger
func queryTesttuples(db LedgerDB, args []string) ([]outputTesttuple, error) {
	outTesttuples := []outputTesttuple{}

	if len(args) != 0 {
		err := errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outTesttuples, err
	}
	elementsKeys, err := db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple"})
	if err != nil {
		return outTesttuples, err
	}
	for _, key := range elementsKeys {
		var out outputTesttuple
		out, err = getOutputTesttuple(db, key)
		if err != nil {
			return outTesttuples, err
		}
		outTesttuples = append(outTesttuples, out)
	}
	return outTesttuples, nil
}

// queryModelDetails returns info about the testtuple and algo related to a traintuple
func queryModelDetails(db LedgerDB, args []string) (outModelDetails outputModelDetails, err error) {
	inp := inputHash{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get associated traintuple
	outModelDetails.Traintuple, err = getOutputTraintuple(db, inp.Key)
	if err != nil {
		return
	}

	// get certified and non-certified testtuples related to traintuple
	testtupleKeys, err := db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", inp.Key})
	if err != nil {
		return
	}
	for _, testtupleKey := range testtupleKeys {
		// get testtuple and serialize it
		var outputTesttuple outputTesttuple
		outputTesttuple, err = getOutputTesttuple(db, testtupleKey)
		if err != nil {
			return
		}

		if outputTesttuple.Certified {
			outModelDetails.Testtuple = outputTesttuple
		} else {
			outModelDetails.NonCertifiedTesttuples = append(outModelDetails.NonCertifiedTesttuples, outputTesttuple)
		}
	}
	return
}

// queryModels returns all traintuples and associated testuples
func queryModels(db LedgerDB, args []string) (outModels []outputModel, err error) {
	outModels = []outputModel{}

	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting nothing")
		return
	}

	traintupleKeys, err := db.GetIndexKeys("traintuple~algo~key", []string{"traintuple"})
	if err != nil {
		return
	}
	for _, traintupleKey := range traintupleKeys {
		var outputModel outputModel

		// get traintuple
		outputModel.Traintuple, err = getOutputTraintuple(db, traintupleKey)
		if err != nil {
			return
		}

		// get associated testtuple
		var testtupleKeys []string
		testtupleKeys, err = db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", traintupleKey, "true"})
		if err != nil {
			return
		}
		if len(testtupleKeys) == 1 {
			// get testtuple and serialize it
			testtupleKey := testtupleKeys[0]
			outputModel.Testtuple, err = getOutputTesttuple(db, testtupleKey)
			if err != nil {
				return
			}
		}
		outModels = append(outModels, outputModel)
	}
	return
}

// --------------------------------------------------------------
// Utils for smartcontracts related to traintuples and testtuples
// --------------------------------------------------------------

// getOutputTraintuple takes as input a traintuple key and returns the outputTraintuple
func getOutputTraintuple(db LedgerDB, traintupleKey string) (outTraintuple outputTraintuple, err error) {
	traintuple, err := db.GetTraintuple(traintupleKey)
	if err != nil {
		return
	}
	outTraintuple.Fill(db, traintuple, traintupleKey)
	return
}

// getOutputTraintuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputTraintuples(db LedgerDB, traintupleKeys []string) (outTraintuples []outputTraintuple, err error) {
	for _, key := range traintupleKeys {
		var outputTraintuple outputTraintuple
		outputTraintuple, err = getOutputTraintuple(db, key)
		if err != nil {
			return
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return
}

// getOutputTesttuple takes as input a testtuple key and returns the outputTesttuple
func getOutputTesttuple(db LedgerDB, testtupleKey string) (outTesttuple outputTesttuple, err error) {
	testtuple, err := db.GetTesttuple(testtupleKey)
	if err != nil {
		return
	}
	err = outTesttuple.Fill(db, testtupleKey, testtuple)
	return
}

// getOutputTesttuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputTesttuples(db LedgerDB, testtupleKeys []string) (outTesttuples []outputTesttuple, err error) {
	for _, key := range testtupleKeys {
		var outputTesttuple outputTesttuple
		outputTesttuple, err = getOutputTesttuple(db, key)
		if err != nil {
			return
		}
		outTesttuples = append(outTesttuples, outputTesttuple)
	}
	return
}

// checkLog checks the validity of logs
func checkLog(log string) (err error) {
	maxLength := 200
	if length := len(log); length > maxLength {
		err = fmt.Errorf("too long log, is %d and should be %d ", length, maxLength)
	}
	return
}

func validateTupleOwner(db LedgerDB, worker string) error {
	txCreator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	if txCreator != worker {
		return fmt.Errorf("%s is not allowed to update tuple (%s)", txCreator, worker)
	}
	return nil
}

// check validity of traintuple update: consistent status and agent submitting the transaction
func checkUpdateTuple(db LedgerDB, worker string, oldStatus string, newStatus string) error {
	statusPossibilities := map[string]string{
		StatusWaiting: StatusTodo,
		StatusTodo:    StatusDoing,
		StatusDoing:   StatusDone}
	if statusPossibilities[oldStatus] != newStatus && newStatus != StatusFailed {
		return errors.BadRequest("cannot change status from %s to %s", oldStatus, newStatus)
	}
	return nil
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (traintuple *Traintuple) validateNewStatus(db LedgerDB, status string) error {
	// check validity of worker and change of status
	if err := checkUpdateTuple(db, traintuple.Dataset.Worker, traintuple.Status, status); err != nil {
		return err
	}
	return nil
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (testtuple *Testtuple) validateNewStatus(db LedgerDB, status string) error {
	// check validity of worker and change of status
	if err := checkUpdateTuple(db, testtuple.Dataset.Worker, testtuple.Status, status); err != nil {
		return err
	}
	return nil
}

// updateTraintupleChildren updates the status of waiting trainuples  InModels of traintuples once they have been trained (succesfully or failed)
func (traintuple *Traintuple) updateTraintupleChildren(db LedgerDB, traintupleKey string) ([]outputTraintuple, error) {

	// tuples to be sent in event
	otuples := []outputTraintuple{}

	// get traintuples having as inModels the input traintuple
	indexName := "traintuple~inModel~key"
	childTraintupleKeys, err := db.GetIndexKeys(indexName, []string{"traintuple", traintupleKey})
	if err != nil {
		return otuples, fmt.Errorf("error while getting associated traintuples to update their inModel")
	}
	for _, childTraintupleKey := range childTraintupleKeys {
		// get and update traintuple
		childTraintuple, err := db.GetTraintuple(childTraintupleKey)
		if err != nil {
			return otuples, err
		}

		// remove associated composite key
		if err := db.DeleteIndex("traintuple~inModel~key", []string{"traintuple", traintupleKey, childTraintupleKey}); err != nil {
			return otuples, err
		}

		// traintuple is already failed, don't update it
		if childTraintuple.Status == StatusFailed {
			continue
		}

		if childTraintuple.Status != StatusWaiting {
			return otuples, fmt.Errorf("traintuple %s has invalid status : '%s' instead of waiting", childTraintupleKey, childTraintuple.Status)
		}

		// get traintuple new status
		var newStatus string
		if traintuple.Status == StatusFailed {
			newStatus = StatusFailed
		} else if traintuple.Status == StatusDone {
			ready, err := childTraintuple.isReady(db, traintupleKey)
			if err != nil {
				return otuples, err
			}
			if ready {
				newStatus = StatusTodo
			}
		}

		// commit new status
		if newStatus == "" {
			continue
		}
		if err := childTraintuple.commitStatusUpdate(db, childTraintupleKey, newStatus); err != nil {
			return otuples, err
		}
		if newStatus == StatusTodo {
			out := outputTraintuple{}
			err = out.Fill(db, childTraintuple, childTraintupleKey)
			if err != nil {
				return otuples, err
			}
			otuples = append(otuples, out)
		}
	}
	return otuples, nil
}

// isReady checks if inModels of a traintuple have been trained, except the newDoneTraintupleKey (since the transaction is not commited)
// and updates the traintuple status if necessary
func (traintuple *Traintuple) isReady(db LedgerDB, newDoneTraintupleKey string) (ready bool, err error) {
	for _, key := range traintuple.InModelKeys {
		// don't check newly done traintuple
		if key == newDoneTraintupleKey {
			continue
		}
		tt, err := db.GetTraintuple(key)
		if err != nil {
			return false, err
		}
		if tt.Status != StatusDone {
			return false, nil
		}
	}
	return true, nil
}

// commitStatusUpdate update the traintuple status in the ledger
func (traintuple *Traintuple) commitStatusUpdate(db LedgerDB, traintupleKey string, newStatus string) error {
	if traintuple.Status == newStatus {
		return fmt.Errorf("cannot update traintuple %s - status already %s", traintupleKey, newStatus)
	}

	if err := traintuple.validateNewStatus(db, newStatus); err != nil {
		return fmt.Errorf("update traintuple %s failed: %s", traintupleKey, err.Error())
	}

	oldStatus := traintuple.Status
	traintuple.Status = newStatus
	if err := db.Put(traintupleKey, traintuple); err != nil {
		return fmt.Errorf("failed to update traintuple %s - %s", traintupleKey, err.Error())
	}

	// update associated composite keys
	indexName := "traintuple~worker~status~key"
	oldAttributes := []string{"traintuple", traintuple.Dataset.Worker, oldStatus, traintupleKey}
	newAttributes := []string{"traintuple", traintuple.Dataset.Worker, traintuple.Status, traintupleKey}
	if err := db.UpdateIndex(indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	logger.Infof("traintuple %s status updated: %s (from=%s)", traintupleKey, newStatus, oldStatus)
	return nil
}

// updateTesttupleChildren update testtuples status associated with a done or failed traintuple
func (traintuple *Traintuple) updateTesttupleChildren(db LedgerDB, traintupleKey string) ([]outputTesttuple, error) {

	otuples := []outputTesttuple{}

	var newStatus string
	if traintuple.Status == StatusFailed {
		newStatus = StatusFailed
	} else if traintuple.Status == StatusDone {
		newStatus = StatusTodo
	} else {
		return otuples, nil
	}

	indexName := "testtuple~traintuple~certified~key"
	// get testtuple associated with this traintuple and updates its status
	testtupleKeys, err := db.GetIndexKeys(indexName, []string{"testtuple", traintupleKey})
	if err != nil {
		return otuples, err
	}
	for _, testtupleKey := range testtupleKeys {
		// get and update testtuple
		testtuple, err := db.GetTesttuple(testtupleKey)
		if err != nil {
			return otuples, err
		}
		testtuple.Model = &Model{
			TraintupleKey: traintupleKey,
		}

		if newStatus == StatusTodo {
			testtuple.Model.Hash = traintuple.OutModel.Hash
			testtuple.Model.StorageAddress = traintuple.OutModel.StorageAddress
		}

		if err := testtuple.commitStatusUpdate(db, testtupleKey, newStatus); err != nil {
			return otuples, err
		}

		if newStatus == StatusTodo {
			out := outputTesttuple{}
			err = out.Fill(db, testtupleKey, testtuple)
			if err != nil {
				return nil, err
			}
			otuples = append(otuples, out)
		}
	}
	return otuples, nil
}

// commitStatusUpdate update the testtuple status in the ledger
func (testtuple *Testtuple) commitStatusUpdate(db LedgerDB, testtupleKey string, newStatus string) error {
	if err := testtuple.validateNewStatus(db, newStatus); err != nil {
		return fmt.Errorf("update testtuple %s failed: %s", testtupleKey, err.Error())
	}

	oldStatus := testtuple.Status
	testtuple.Status = newStatus

	if err := db.Put(testtupleKey, testtuple); err != nil {
		return fmt.Errorf("failed to update testtuple status to %s with key %s", newStatus, testtupleKey)
	}

	// update associated composite key
	indexName := "testtuple~worker~status~key"
	oldAttributes := []string{"testtuple", testtuple.Dataset.Worker, oldStatus, testtupleKey}
	newAttributes := []string{"testtuple", testtuple.Dataset.Worker, testtuple.Status, testtupleKey}
	if err := db.UpdateIndex(indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	logger.Infof("testtuple %s status updated: %s (from=%s)", testtupleKey, newStatus, oldStatus)
	return nil
}

// getTraintuplesPayload takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getTraintuplesPayload(db LedgerDB, traintupleKeys []string) ([]map[string]interface{}, error) {

	var elements []map[string]interface{}
	for _, key := range traintupleKeys {
		var element map[string]interface{}
		outputTraintuple, err := getOutputTraintuple(db, key)
		if err != nil {
			return nil, err
		}
		oo, err := json.Marshal(outputTraintuple)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(oo, &element)
		element["key"] = key
		elements = append(elements, element)
	}
	return elements, nil
}

// HashForKey to generate key for an asset
func HashForKey(objectType string, hashElements ...string) string {
	toHash := objectType
	sort.Strings(hashElements)
	for _, element := range hashElements {
		toHash += "," + element
	}
	sum := sha256.Sum256([]byte(toHash))
	return hex.EncodeToString(sum[:])
}
