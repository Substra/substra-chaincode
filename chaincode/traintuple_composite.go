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
)

// ------------------------------------------
// Methods on receivers composite traintuple
// ------------------------------------------

// SetFromInput is a method of the receiver CompositeTraintuple.
// It uses the inputCompositeTraintuple to check and set the traintuple's parameters
// which don't depend on previous traintuples values :
//  - AssetType
//  - Creator & permissions
//  - Tag
//  - AlgoKey & ObjectiveKey
//  - Dataset
func (traintuple *CompositeTraintuple) SetFromInput(db *LedgerDB, inp inputCompositeTraintuple) error {

	creator, err := GetTxCreator(db.cc)
	if err != nil {
		return err
	}
	traintuple.AssetType = CompositeTraintupleType
	traintuple.Creator = creator
	traintuple.Tag = inp.Tag
	algo, err := db.GetCompositeAlgo(inp.AlgoKey)
	if err != nil {
		return errors.BadRequest(err, "could not retrieve Composite algo with key %s", inp.AlgoKey)
	}
	if !algo.Permissions.CanProcess(algo.Owner, creator) {
		return errors.Forbidden("not authorized to process algo %s", inp.AlgoKey)
	}
	traintuple.AlgoKey = inp.AlgoKey

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
		return errors.Forbidden("not authorized to process dataManager %s", inp.DataManagerKey)
	}

	// fill traintuple.Dataset from dataManager and dataSample
	traintuple.Dataset = &Dataset{
		DataManagerKey: inp.DataManagerKey,
		DataSampleKeys: inp.DataSampleKeys,
	}
	traintuple.Dataset.Worker, err = getDataManagerOwner(db, traintuple.Dataset.DataManagerKey)

	// permissions (head): worker only where the data belong
	workerOnly := Permission{
		Public:        false,
		AuthorizedIDs: []string{traintuple.Dataset.Worker}}
	traintuple.OutHeadModel.Permissions = Permissions{Process: workerOnly, Download: workerOnly}

	// permissions (trunk): dictated by input
	permissions, err := NewPermissions(db, inp.OutTrunkModelPermissions)
	if err != nil {
		return err
	}
	traintuple.OutTrunkModel.Permissions = permissions

	return err
}

// SetFromParents set the status of the traintuple depending on its "parents",
// i.e. the traintuples from which it received the outModels as inModels.
// Also it's InModelKeys are set.
// TODO: rename to SetInModels
func (traintuple *CompositeTraintuple) SetFromParents(db *LedgerDB, inp inputCompositeTraintuple) error {
	traintuple.Status = StatusTodo
	if inp.InHeadModelKey == "" || inp.InTrunkModelKey == "" {
		return nil
	}

	// [Head]
	// It can only be a composite traintuple's head out model
	traintuple.InHeadModel = inp.InHeadModelKey
	head, err := db.GetGenericTuple(inp.InHeadModelKey)
	if err != nil {
		return err
	}
	if !typeInSlice(head.AssetType, []AssetType{CompositeTraintupleType}) {
		return errors.BadRequest(
			"tuple type %s from key %s is not supported as head InModel",
			head.AssetType,
			inp.InHeadModelKey)
	}
	// [Trunk]
	// It can be either:
	// - a traintuple's out model
	// - a composite traintuple's trunk out model
	// - an aggregate tuple's out model
	traintuple.InTrunkModel = inp.InTrunkModelKey
	trunk, err := db.GetGenericTuple(inp.InTrunkModelKey)
	if err != nil {
		return err
	}
	if !typeInSlice(trunk.AssetType, []AssetType{TraintupleType, CompositeTraintupleType, AggregatetupleType}) {
		return errors.BadRequest(
			"tuple type %s from key %s is not supported as trunk InModel",
			trunk.AssetType,
			inp.InTrunkModelKey)
	}
	traintuple.Status = determineStatusFromInModels([]string{head.Status, trunk.Status})
	return nil
}

// GetKey return the key of the traintuple depending on its key parameters.
func (traintuple *CompositeTraintuple) GetKey() string {
	hashKeys := []string{
		traintuple.Creator,
		traintuple.AlgoKey,
		traintuple.Dataset.DataManagerKey,
		traintuple.InHeadModel,
		traintuple.InTrunkModel}
	hashKeys = append(hashKeys, traintuple.Dataset.DataSampleKeys...)
	return HashForKey("compositeTraintuple", hashKeys...)
}

// AddToComputePlan set the traintuple's parameters that determines if it's part of on ComputePlan and how.
// It uses the inputCompositeTraintuple values as follow:
//  - If neither ComputePlanID nor rank is set it returns immediately
//  - If rank is 0 and ComputePlanID empty, it's start a new one using this traintuple key
//  - If rank and ComputePlanID are set, it checks if there are coherent with previous ones and set it.
// Use checkComputePlanAvailability to ensure the compute plan exists and no other tuple is registered with the same worker/rank
func (traintuple *CompositeTraintuple) AddToComputePlan(db *LedgerDB, inp inputCompositeTraintuple, traintupleKey string, checkComputePlanAvailability bool) error {
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
		computePlan := ComputePlan{Status: traintuple.Status, CompositeTraintupleKeys: []string{traintupleKey}}
		traintuple.ComputePlanID, err = computePlan.Create(db)
		if err != nil {
			return err
		}
		return nil
	}
	traintuple.ComputePlanID = inp.ComputePlanID
	computePlan, err := db.GetComputePlan(inp.ComputePlanID)
	if err != nil {
		return err
	}
	computePlan.CompositeTraintupleKeys = append(computePlan.CompositeTraintupleKeys, traintupleKey)
	computePlan.TupleCount++
	computePlan.CheckNewTupleStatus(traintuple.Status)
	err = computePlan.Save(db, traintuple.ComputePlanID)
	if err != nil {
		return err
	}

	if !checkComputePlanAvailability {
		return nil
	}
	var ttKeys []string
	ttKeys, err = db.GetIndexKeys("computePlan~computeplanid~worker~rank~key", []string{"computePlan", inp.ComputePlanID, traintuple.Dataset.Worker, inp.Rank})
	if err != nil {
		return err
	} else if len(ttKeys) > 0 {
		err = errors.BadRequest("ComputePlanID %s with worker %s rank %d already exists", inp.ComputePlanID, traintuple.Dataset.Worker, traintuple.Rank)
		return err
	}
	return nil
}

// Save will put in the legder interface both the traintuple with its key
// and all the associated composite keys
func (traintuple *CompositeTraintuple) Save(db *LedgerDB, traintupleKey string) error {

	// store in ledger
	if err := db.Add(traintupleKey, traintuple); err != nil {
		return err
	}

	// create composite keys
	if err := db.CreateIndex("compositeTraintuple~algo~key", []string{"compositeTraintuple", traintuple.AlgoKey, traintupleKey}); err != nil {
		return err
	}
	if err := db.CreateIndex("compositeTraintuple~worker~status~key", []string{"compositeTraintuple", traintuple.Dataset.Worker, traintuple.Status, traintupleKey}); err != nil {
		return err
	}
	// TODO: Do we create an index for head/trunk inModel or do we concider that
	// they are classic inModels ?
	if err := db.CreateIndex("compositeTraintuple~inModel~key", []string{"compositeTraintuple", traintuple.InHeadModel, traintupleKey}); err != nil {
		return err
	}
	if err := db.CreateIndex("compositeTraintuple~inModel~key", []string{"compositeTraintuple", traintuple.InTrunkModel, traintupleKey}); err != nil {
		return err
	}
	if traintuple.ComputePlanID != "" {
		if err := db.CreateIndex("computePlan~computeplanid~worker~rank~key", []string{"computePlan", traintuple.ComputePlanID, traintuple.Dataset.Worker, strconv.Itoa(traintuple.Rank), traintupleKey}); err != nil {
			return err
		}
	}
	if traintuple.Tag != "" {
		err := db.CreateIndex("compositeTraintuple~tag~key", []string{"compositeTraintuple", traintuple.Tag, traintupleKey})
		if err != nil {
			return err
		}
	}
	return nil
}

// -------------------------------------------------
// Smart contracts related to composite traintuples
// -------------------------------------------------

// createCompositeTraintuple is the wrapper for the substra smartcontract createCompositeTraintuple
func createCompositeTraintuple(db *LedgerDB, args []string) (map[string]string, error) {
	inp := inputCompositeTraintuple{}
	err := AssetFromJSON(args, &inp)
	if err != nil {
		return nil, err
	}

	key, err := createCompositeTraintupleInternal(db, inp, true)
	if err != nil {
		return nil, err
	}

	return map[string]string{"key": key}, nil
}

// createCompositeTraintupleInternal adds a CompositeTraintuple in the ledger
func createCompositeTraintupleInternal(db *LedgerDB, inp inputCompositeTraintuple, checkComputePlanAvailability bool) (string, error) {
	traintuple := CompositeTraintuple{}
	err := traintuple.SetFromInput(db, inp)
	if err != nil {
		return "", err
	}
	err = traintuple.SetFromParents(db, inp)
	if err != nil {
		return "", err
	}

	traintupleKey := traintuple.GetKey()
	// Test if the key (ergo the traintuple) already exists
	tupleExists, err := db.KeyExists(traintupleKey)
	if err != nil {
		return "", err
	}
	if tupleExists {
		return "", errors.Conflict("composite traintuple already exists").WithKey(traintupleKey)
	}

	err = traintuple.AddToComputePlan(db, inp, traintupleKey, checkComputePlanAvailability)
	if err != nil {
		return "", err
	}

	err = traintuple.Save(db, traintupleKey)
	if err != nil {
		return "", err
	}
	err = db.AddTupleEvent(traintupleKey)
	if err != nil {
		return "", err
	}
	return traintupleKey, nil
}

// logStartCompositeTrain modifies a traintuple by changing its status from todo to doing
func logStartCompositeTrain(db *LedgerDB, args []string) (o outputCompositeTraintuple, err error) {
	status := StatusDoing
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get compositeTraintuple, check validity of the update
	compositeTraintuple, err := db.GetCompositeTraintuple(inp.Key)
	if err != nil {
		return
	}

	if err = validateTupleOwner(db, compositeTraintuple.Dataset.Worker); err != nil {
		return
	}
	if err = compositeTraintuple.commitStatusUpdate(db, inp.Key, status); err != nil {
		return
	}
	err = o.Fill(db, compositeTraintuple, inp.Key)
	if err != nil {
		return
	}
	err = UpdateComputePlan(db, compositeTraintuple.ComputePlanID, compositeTraintuple.Status)
	return
}

// logSuccessCompositeTrain modifies a traintuple by changing its status from doing to done
// reports logs and associated performances
func logSuccessCompositeTrain(db *LedgerDB, args []string) (o outputCompositeTraintuple, err error) {
	status := StatusDone
	inp := inputLogSuccessCompositeTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	compositeTraintupleKey := inp.Key

	// get, update and commit traintuple
	compositeTraintuple, err := db.GetCompositeTraintuple(compositeTraintupleKey)
	if err != nil {
		return
	}

	compositeTraintuple.OutHeadModel.OutModel = &Hash{
		Hash:           inp.OutHeadModel.Hash}

	compositeTraintuple.OutTrunkModel.OutModel = &HashDress{
		Hash:           inp.OutTrunkModel.Hash,
		StorageAddress: inp.OutTrunkModel.StorageAddress}
	compositeTraintuple.Log += inp.Log

	err = db.CreateIndex("tuple~modelHash~key", []string{"tuple", inp.OutHeadModel.Hash, compositeTraintupleKey})
	if err != nil {
		return
	}
	err = db.CreateIndex("tuple~modelHash~key", []string{"tuple", inp.OutTrunkModel.Hash, compositeTraintupleKey})
	if err != nil {
		return
	}
	if err = validateTupleOwner(db, compositeTraintuple.Dataset.Worker); err != nil {
		return
	}
	if err = compositeTraintuple.commitStatusUpdate(db, compositeTraintupleKey, status); err != nil {
		return
	}

	err = UpdateTraintupleChildren(db, compositeTraintupleKey, compositeTraintuple.Status, []string{})
	if err != nil {
		return
	}

	err = UpdateTesttupleChildren(db, compositeTraintupleKey, compositeTraintuple.Status)
	if err != nil {
		return
	}

	err = o.Fill(db, compositeTraintuple, inp.Key)
	if err != nil {
		return
	}

	err = UpdateComputePlan(db, compositeTraintuple.ComputePlanID, compositeTraintuple.Status)
	return
}

// logFailCompositeTrain modifies a traintuple by changing its status to fail and reports associated logs
func logFailCompositeTrain(db *LedgerDB, args []string) (o outputCompositeTraintuple, err error) {
	status := StatusFailed
	inp := inputLogFailTrain{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}

	// get, update and commit traintuple
	compositeTraintuple, err := db.GetCompositeTraintuple(inp.Key)
	if err != nil {
		return
	}

	compositeTraintuple.Log += inp.Log

	if err = validateTupleOwner(db, compositeTraintuple.Dataset.Worker); err != nil {
		return
	}
	if err = compositeTraintuple.commitStatusUpdate(db, inp.Key, status); err != nil {
		return
	}

	err = o.Fill(db, compositeTraintuple, inp.Key)
	if err != nil {
		return
	}
	err = UpdateComputePlan(db, compositeTraintuple.ComputePlanID, compositeTraintuple.Status)
	if err != nil {
		return
	}
	// Do not propagate failure if we are in a compute plan
	if compositeTraintuple.ComputePlanID != "" {
		return
	}
	// update depending tuples
	err = UpdateTesttupleChildren(db, inp.Key, compositeTraintuple.Status)
	if err != nil {
		return
	}

	err = UpdateTraintupleChildren(db, inp.Key, compositeTraintuple.Status, []string{})
	return
}

// queryCompositeTraintuple returns info about a composite traintuple given its key
func queryCompositeTraintuple(db *LedgerDB, args []string) (outputTraintuple outputCompositeTraintuple, err error) {
	inp := inputKey{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	traintuple, err := db.GetCompositeTraintuple(inp.Key)
	if err != nil {
		return
	}
	if traintuple.AssetType != CompositeTraintupleType {
		err = errors.NotFound("no element with key %s", inp.Key)
		return
	}
	outputTraintuple.Fill(db, traintuple, inp.Key)
	return
}

// queryCompositeTraintuples returns all composite traintuples
func queryCompositeTraintuples(db *LedgerDB, args []string) ([]outputCompositeTraintuple, error) {
	outTraintuples := []outputCompositeTraintuple{}

	if len(args) != 0 {
		err := errors.BadRequest("incorrect number of arguments, expecting nothing")
		return outTraintuples, err
	}
	elementsKeys, err := db.GetIndexKeys("compositeTraintuple~algo~key", []string{"compositeTraintuple"})
	if err != nil {
		return outTraintuples, err
	}

	nb := getLimitedNbSliceElements(elementsKeys)
	for _, key := range elementsKeys[:nb] {
		outputTraintuple, err := getOutputCompositeTraintuple(db, key)
		if err != nil {
			return outTraintuples, err
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return outTraintuples, nil
}

// ----------------------------------------------------------
// Utils for smartcontracts related to composite traintuples
// ----------------------------------------------------------

// UpdateCompositeTraintupleChild updates the status of a waiting trainuple, given the new parent traintuple status
func UpdateCompositeTraintupleChild(db *LedgerDB, parentTraintupleKey string, childTraintupleKey string, traintupleStatus string) (childStatus string, err error) {
	// get and update traintuple
	childTraintuple, err := db.GetCompositeTraintuple(childTraintupleKey)
	if err != nil {
		return
	}

	childStatus = childTraintuple.Status

	// get traintuple new status
	var newStatus string
	if traintupleStatus == StatusFailed {
		newStatus = StatusFailed
	} else if traintupleStatus == StatusDone {
		ready, _err := childTraintuple.isReady(db, parentTraintupleKey)
		if _err != nil {
			err = _err
			return
		}
		if ready {
			newStatus = StatusTodo
		}
	}

	// commit new status
	if newStatus == "" {
		return
	}
	if err = childTraintuple.commitStatusUpdate(db, childTraintupleKey, newStatus); err != nil {
		return
	}

	// update return value after status update
	childStatus = childTraintuple.Status

	err = db.AddTupleEvent(childTraintupleKey)
	return
}

// getOutputCompositeTraintuple takes as input a traintuple key and returns the outputCompositeTraintuple
func getOutputCompositeTraintuple(db *LedgerDB, traintupleKey string) (outTraintuple outputCompositeTraintuple, err error) {
	traintuple, err := db.GetCompositeTraintuple(traintupleKey)
	if err != nil {
		return
	}
	outTraintuple.Fill(db, traintuple, traintupleKey)
	return
}

// getOutputCompositeTraintuples takes as input a list of keys and returns a paylaod containing a list of associated retrieved elements
func getOutputCompositeTraintuples(db *LedgerDB, traintupleKeys []string) (outTraintuples []outputCompositeTraintuple, err error) {
	for _, key := range traintupleKeys {
		var outputTraintuple outputCompositeTraintuple
		outputTraintuple, err = getOutputCompositeTraintuple(db, key)
		if err != nil {
			return
		}
		outTraintuples = append(outTraintuples, outputTraintuple)
	}
	return
}

// validateNewStatus verifies that the new status is consistent with the tuple current status
func (traintuple *CompositeTraintuple) validateNewStatus(db *LedgerDB, status string) error {
	// check validity of worker and change of status
	if err := checkUpdateTuple(db, traintuple.Dataset.Worker, traintuple.Status, status); err != nil {
		return err
	}
	return nil
}

func (traintuple *CompositeTraintuple) isReady(db *LedgerDB, newDoneTraintupleKey string) (ready bool, err error) {
	return IsReady(db, []string{traintuple.InHeadModel, traintuple.InTrunkModel}, newDoneTraintupleKey)
}

// commitStatusUpdate update the traintuple status in the ledger
func (traintuple *CompositeTraintuple) commitStatusUpdate(db *LedgerDB, traintupleKey string, newStatus string) error {
	if traintuple.Status == newStatus {
		return nil
	}

	// do not update if previous status is already Done, Failed, Todo, Doing
	if StatusAborted == newStatus && traintuple.Status != StatusWaiting {
		return nil
	}

	if err := traintuple.validateNewStatus(db, newStatus); err != nil {
		return errors.Internal("update traintuple %s failed: %s", traintupleKey, err.Error())
	}

	oldStatus := traintuple.Status
	traintuple.Status = newStatus
	if err := db.Put(traintupleKey, traintuple); err != nil {
		return errors.Internal("failed to update traintuple %s - %s", traintupleKey, err.Error())
	}

	// update associated composite keys
	indexName := "compositeTraintuple~worker~status~key"
	oldAttributes := []string{"compositeTraintuple", traintuple.Dataset.Worker, oldStatus, traintupleKey}
	newAttributes := []string{"compositeTraintuple", traintuple.Dataset.Worker, traintuple.Status, traintupleKey}
	if err := db.UpdateIndex(indexName, oldAttributes, newAttributes); err != nil {
		return err
	}
	logger.Infof("traintuple %s status updated: %s (from=%s)", traintupleKey, newStatus, oldStatus)
	return nil
}
