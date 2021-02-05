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
	"encoding/json"
	"strings"
	"sync"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// State is a in-memory representation of the db state
type State struct {
	items map[string]([]byte)
}

// LedgerDB to access the chaincode database during the lifetime of a SmartContract
type LedgerDB struct {
	cc               shim.ChaincodeStubInterface
	event            *Event
	transactionState State
	mutex            *sync.RWMutex
}

// NewLedgerDB create a new db to access the chaincode during a SmartContract
func NewLedgerDB(stub shim.ChaincodeStubInterface) *LedgerDB {
	return &LedgerDB{
		cc: stub,
		transactionState: State{
			items: make(map[string]([]byte)),
		},
		mutex: &sync.RWMutex{},
	}
}

// ----------------------------------------------
// Low-level functions to handle asset structs
// ----------------------------------------------

// getTransactionState returns a copy of an object that has been updated or created during the transaction
func (db *LedgerDB) getTransactionState(key string) ([]byte, bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	transactionState, ok := db.transactionState.items[key]
	if !ok {
		return nil, false
	}
	state := make([]byte, len(transactionState))
	copy(state, transactionState)
	return state, true
}

// putTransactionState stores an object during a transaction lifetime
func (db *LedgerDB) putTransactionState(key string, state []byte) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.transactionState.items[key] = state
}

// Get retrieves an object stored in the chaincode db and set the input object value
func (db *LedgerDB) Get(key string, object interface{}) error {
	var buff []byte
	var err error

	buff, ok := db.getTransactionState(key)
	if !ok {
		buff, err = db.cc.GetState(key)
		if err != nil || buff == nil {
			return errors.NotFound(err, "no asset for key %s", key)
		}
		db.putTransactionState(key, buff)
	}

	return json.Unmarshal(buff, &object)
}

// KeyExists checks if a key is stored in the chaincode db
func (db *LedgerDB) KeyExists(key string) (bool, error) {
	buff, err := db.cc.GetState(key)
	return buff != nil, err
}

// Put stores an object in the chaincode db, if the object already exists it is replaced
func (db *LedgerDB) Put(key string, object interface{}) error {
	buff, _ := json.Marshal(object)

	if err := db.cc.PutState(key, buff); err != nil {
		return err
	}
	// TransactionState is updated to ensure that even if the data is not committed, a further
	// call to get this struct will returned the updated one (and not the original one).
	// This is currently required when setting the statuses of the traintuple children.
	db.putTransactionState(key, buff)

	return nil
}

// Add stores an object in the chaincode db, it fails if the object already exists
func (db *LedgerDB) Add(key string, object interface{}) error {
	ok, err := db.KeyExists(key)
	if err != nil {
		return err
	}
	if ok {
		return errors.Conflict("struct already exists (tkey: %s)", key).WithKey(key)
	}
	return db.Put(key, object)
}

// ----------------------------------------------
// Low-level functions to handle indexes
// ----------------------------------------------

// CreateIndex adds a new composite key to the chaincode db
func (db *LedgerDB) CreateIndex(index string, attributes []string) error {
	compositeKey, err := db.cc.CreateCompositeKey(index, attributes)
	if err != nil {
		return errors.Internal("cannot create index %s: %s", index, err.Error())
	}
	value := []byte{0x00}
	if err = db.cc.PutState(compositeKey, value); err != nil {
		return errors.Internal("cannot create index %s: %s", index, err.Error())
	}
	return nil
}

// DeleteIndex deletes a composite key in the chaincode db
func (db *LedgerDB) DeleteIndex(index string, attributes []string) error {
	compositeKey, err := db.cc.CreateCompositeKey(index, attributes)
	if err != nil {
		return err
	}
	return db.cc.DelState(compositeKey)
}

// UpdateIndex updates an existing composite key in the chaincode db
func (db *LedgerDB) UpdateIndex(index string, oldAttributes []string, newAttribues []string) error {
	if err := db.DeleteIndex(index, oldAttributes); err != nil {
		return err
	}
	return db.CreateIndex(index, newAttribues)
}

// GetIndexKeys returns keys matching composite key values from the chaincode db
func (db *LedgerDB) GetIndexKeys(index string, attributes []string) ([]string, error) {
	keys := make([]string, 0)
	iterator, err := db.cc.GetStateByPartialCompositeKey(index, attributes)
	if err != nil {
		return nil, errors.Internal("get index %s failed: %s", index, err.Error())
	}
	defer iterator.Close()
	for iterator.HasNext() {
		compositeKey, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		_, keyParts, err := db.cc.SplitCompositeKey(compositeKey.Key)
		if err != nil {
			return nil, errors.Internal("get index %s failed: cannot split key %s: %s", index, compositeKey.Key, err.Error())
		}
		keys = append(keys, keyParts[len(keyParts)-1])
	}
	return keys, nil
}

// GetIndexKeysWithPagination returns keys matching composite key values from the chaincode db
func (db *LedgerDB) GetIndexKeysWithPagination(index string, attributes []string, pageSize int32, bookmark string) ([]string, string, error) {
	keys := make([]string, 0)

	if bookmark != "" {
		// Transform bookmark from JSON-friendly format to CouchDB format
		bookmark = strings.Replace(bookmark, "/", "\x00", -1)
		bookmark = strings.Replace(bookmark, "#", "\\u0000", -1)
		bookmark = strings.Replace(bookmark, "END", "\U0010ffff", -1)
	}

	iterator, metadata, err := db.cc.GetStateByPartialCompositeKeyWithPagination(index, attributes, pageSize, bookmark)
	if err != nil {
		return nil, "", errors.Internal("get index %s failed: %s", index, err.Error())
	}
	defer iterator.Close()
	for iterator.HasNext() {
		compositeKey, err := iterator.Next()
		if err != nil {
			return nil, "", err
		}
		_, keyParts, err := db.cc.SplitCompositeKey(compositeKey.Key)
		if err != nil {
			return nil, "", errors.Internal("get index %s failed: cannot split key %s: %s", index, compositeKey.Key, err.Error())
		}
		keys = append(keys, keyParts[len(keyParts)-1])
	}

	if metadata != nil {
		// Transform bookmark from CouchDB format to JSON-friendly format
		bookmark = strings.Replace(metadata.Bookmark, "\x00", "/", -1)
		bookmark = strings.Replace(bookmark, "\\u0000", "#", -1)
		bookmark = strings.Replace(bookmark, "\U0010ffff", "END", -1)
	}

	return keys, bookmark, nil
}

// ----------------------------------------------
// High-level functions
// ----------------------------------------------

// GetAssetType fetch a object in the chaincode db and return its type.
// It fails if it does not exists or if it does not have an `AssetType` field.
func (db *LedgerDB) GetAssetType(key string) (AssetType, error) {
	asset := struct {
		AssetType AssetType `json:"asset_type"`
	}{}
	if err := db.Get(key, &asset); err != nil {
		return asset.AssetType, err
	}
	return asset.AssetType, nil
}

// GetGenericTuple fetches a GenericTuple (Traintuple, CompositeTraintuple or AggregateTuple)
// from the chaincode db
func (db *LedgerDB) GetGenericTuple(key string) (GenericTuple, error) {
	asset := GenericTuple{}
	err := db.Get(key, &asset)
	if err != nil {
		return asset, err
	}
	asset.Status, err = determineTupleStatus(db, asset.Status, asset.ComputePlanKey)
	return asset, nil
}

// GetStatusUpdater fetches a GenericTuple (Traintuple, CompositeTraintuple or AggregateTuple)
// from the chaincode db
func (db *LedgerDB) GetStatusUpdater(key string) (StatusUpdater, error) {
	tType, err := db.GetAssetType(key)
	if err != nil {
		return nil, err
	}
	var asset StatusUpdater
	switch tType {
	case TraintupleType:
		t, err := db.GetTraintuple(key)
		if err != nil {
			return nil, err
		}
		asset = &t
	case CompositeTraintupleType:
		t, err := db.GetCompositeTraintuple(key)
		if err != nil {
			return nil, err
		}
		asset = &t
	case AggregatetupleType:
		t, err := db.GetAggregatetuple(key)
		if err != nil {
			return nil, err
		}
		asset = &t
	case TesttupleType:
		t, err := db.GetTesttuple(key)
		if err != nil {
			return nil, err
		}
		asset = &t
	}
	return asset, nil
}

// GetAlgo fetches an Algo from the ledger using its unique key
func (db *LedgerDB) GetAlgo(key string) (Algo, error) {
	algo := Algo{}
	if err := db.Get(key, &algo); err != nil {
		return algo, err
	}
	if algo.AssetType != AlgoType {
		return algo, errors.NotFound("algo %s not found", key)
	}
	return algo, nil
}

// GetCompositeAlgo fetches a CompositeAlgo from the ledger using its unique key
func (db *LedgerDB) GetCompositeAlgo(key string) (CompositeAlgo, error) {
	algo := CompositeAlgo{}
	if err := db.Get(key, &algo); err != nil {
		return algo, err
	}
	if algo.AssetType != CompositeAlgoType {
		return algo, errors.NotFound("algo %s not found", key)
	}
	return algo, nil
}

// GetAggregateAlgo fetches a AggregateAlgo from the ledger using its unique key
func (db *LedgerDB) GetAggregateAlgo(key string) (AggregateAlgo, error) {
	algo := AggregateAlgo{}
	if err := db.Get(key, &algo); err != nil {
		return algo, err
	}
	if algo.AssetType != AggregateAlgoType {
		return algo, errors.NotFound("algo %s not found", key)
	}
	return algo, nil
}

// GetObjective fetches an Objective from the ledger using its unique key
func (db *LedgerDB) GetObjective(key string) (Objective, error) {
	objective := Objective{}
	if err := db.Get(key, &objective); err != nil {
		return objective, err
	}
	if objective.AssetType != ObjectiveType {
		return objective, errors.NotFound("objective %s not found", key)
	}
	return objective, nil
}

// GetDataManager fetches a DataManager from the ledger using its unique key
func (db *LedgerDB) GetDataManager(key string) (DataManager, error) {
	dataManager := DataManager{}
	if err := db.Get(key, &dataManager); err != nil {
		return dataManager, err
	}
	if dataManager.AssetType != DataManagerType {
		return dataManager, errors.NotFound("dataManager %s not found", key)
	}
	return dataManager, nil
}

// GetDataSample fetches a DataSample from the ledger using its unique key
func (db *LedgerDB) GetDataSample(key string) (DataSample, error) {
	dataSample := DataSample{}
	if err := db.Get(key, &dataSample); err != nil {
		return dataSample, err
	}
	if dataSample.AssetType != DataSampleType {
		return dataSample, errors.NotFound("dataSample %s not found", key)
	}
	return dataSample, nil
}

// GetTraintuple fetches a Traintuple from the ledger using its unique key
func (db *LedgerDB) GetTraintuple(key string) (Traintuple, error) {
	traintuple := Traintuple{}
	err := db.Get(key, &traintuple)
	if err != nil {
		return traintuple, err
	}
	if traintuple.AssetType != TraintupleType {
		return traintuple, errors.NotFound("traintuple %s not found", key)
	}
	traintuple.Status, err = determineTupleStatus(db, traintuple.Status, traintuple.ComputePlanKey)
	return traintuple, err
}

// GetCompositeTraintuple fetches a CompositeTraintuple from the ledger using its unique key
func (db *LedgerDB) GetCompositeTraintuple(key string) (CompositeTraintuple, error) {
	traintuple := CompositeTraintuple{}
	err := db.Get(key, &traintuple)
	if err != nil {
		return traintuple, err
	}
	if traintuple.AssetType != CompositeTraintupleType {
		return traintuple, errors.NotFound("composite traintuple %s not found", key)
	}
	traintuple.Status, err = determineTupleStatus(db, traintuple.Status, traintuple.ComputePlanKey)
	return traintuple, err
}

// GetAggregatetuple fetches a Aggregatetuple from the ledger using its unique key
func (db *LedgerDB) GetAggregatetuple(key string) (Aggregatetuple, error) {
	aggregatetuple := Aggregatetuple{}
	err := db.Get(key, &aggregatetuple)
	if err != nil {
		return aggregatetuple, err
	}
	if aggregatetuple.AssetType != AggregatetupleType {
		return aggregatetuple, errors.NotFound("aggregatetuple %s not found", key)
	}
	aggregatetuple.Status, err = determineTupleStatus(db, aggregatetuple.Status, aggregatetuple.ComputePlanKey)
	return aggregatetuple, err
}

// GetComputePlan fetches a ComputePlan from the ledger using its unique ID
func (db *LedgerDB) GetComputePlan(key string) (ComputePlan, error) {
	computePlan := ComputePlan{}
	if err := db.Get(key, &computePlan); err != nil {
		return computePlan, err
	}
	if computePlan.AssetType != ComputePlanType {
		return computePlan, errors.NotFound("compute plan %s not found", key)
	}
	if err := db.Get(computePlan.StateKey, &(computePlan.State)); err != nil {
		return computePlan, err
	}
	return computePlan, nil
}

// GetCPWorkerState returns the state for a given compute plan and worker
func (db *LedgerDB) GetCPWorkerState(wStateKey string) (*ComputePlanWorkerState, error) {
	wState := ComputePlanWorkerState{}
	err := db.Get(wStateKey, &wState)
	if err != nil {
		return nil, err
	}
	return &wState, nil
}

// GetOutModelKeyChecksumAddress retrieves an out-Model from a tuple key.
// In case of CompositeTraintuple it return its trunk model
// Return an error if the tupleKey was not found.
func (db *LedgerDB) GetOutModelKeyChecksumAddress(tupleKey string, allowedAssetTypes []AssetType) (*KeyChecksumAddress, error) {
	for _, assetType := range allowedAssetTypes {
		switch assetType {
		case CompositeTraintupleType:
			tuple, err := db.GetCompositeTraintuple(tupleKey)
			if err != nil {
				continue
			}
			return tuple.OutTrunkModel.OutModel, nil
		case TraintupleType:
			tuple, err := db.GetTraintuple(tupleKey)
			if err == nil {
				return tuple.OutModel, nil
			}
		case AggregatetupleType:
			tuple, err := db.GetAggregatetuple(tupleKey)
			if err == nil {
				return tuple.OutModel, nil
			}
		default:
			return nil, errors.Internal("GetOutModelKeyChecksumAddress: Unsupported asset type %s", assetType)
		}
	}

	return nil, errors.NotFound(
		"GetOutModelKeyChecksumAddress: Could not find tuple with key \"%s\". Allowed types: %v.",
		tupleKey,
		allowedAssetTypes)
}

// GetOutHeadModelKeyChecksum retrieves an out-Head-Model from a composite traintuple key.
// Return an error if the compositeTraintupleKey was not found.
func (db *LedgerDB) GetOutHeadModelKeyChecksum(compositeTraintupleKey string) (*KeyChecksum, error) {
	tuple, err := db.GetCompositeTraintuple(compositeTraintupleKey)
	if err != nil {
		return nil, err
	}
	return tuple.OutHeadModel.OutModel, nil
}

// GetTesttuple fetches a Testtuple from the ledger using its unique key
func (db *LedgerDB) GetTesttuple(key string) (Testtuple, error) {
	testtuple := Testtuple{}
	err := db.Get(key, &testtuple)
	if err != nil {
		return testtuple, err
	}
	if testtuple.AssetType != TesttupleType {
		return testtuple, errors.NotFound("testtuple %s not found", key)
	}
	testtuple.Status, err = determineTupleStatus(db, testtuple.Status, testtuple.ComputePlanKey)
	return testtuple, err
}

// GetNode fetches a Node from the ledger based on its unique key
func (db *LedgerDB) GetNode(key string) (Node, error) {
	node := Node{}

	err := db.Get(key, &node)
	if err != nil {
		return node, err
	}

	return node, nil
}

// ----------------------------------------------
// High-level functions for events
// ----------------------------------------------

// SendEvent sends an event with updated tuples if there is any
// Only one event can be sent per transaction
func (db *LedgerDB) SendEvent() error {
	if db.event == nil {
		return nil
	}
	payload, err := json.Marshal(*(db.event))
	if err != nil {
		return err
	}
	err = db.cc.SetEvent("chaincode-updates", payload)
	if err != nil {
		return err
	}
	return nil
}

// AddTupleEvent add the output tuple matching the tupleKey to the event struct
func (db *LedgerDB) AddTupleEvent(tupleKey string) error {
	// We take advantage of the fact that Testtuples have the fields "AssetType"
	// and "Status": we use db.GetGenericTuple to get the value for these fields
	// even though Testtuples aren't technically GenericTuples.
	genericTuple, err := db.GetGenericTuple(tupleKey)
	if err != nil {
		return err
	}
	if genericTuple.Status != StatusTodo {
		return nil
	}
	if db.event == nil {
		db.event = &Event{}
	}
	switch genericTuple.AssetType {
	case TraintupleType:
		tuple, err := db.GetTraintuple(tupleKey)
		if err != nil {
			return err
		}
		out := outputTraintuple{}
		out.Fill(db, tuple)
		db.event.Traintuples = append(db.event.Traintuples, out)
	case CompositeTraintupleType:
		tuple, err := db.GetCompositeTraintuple(tupleKey)
		if err != nil {
			return err
		}
		out := outputCompositeTraintuple{}
		out.Fill(db, tuple)
		db.event.CompositeTraintuples = append(db.event.CompositeTraintuples, out)
	case AggregatetupleType:
		tuple, err := db.GetAggregatetuple(tupleKey)
		if err != nil {
			return err
		}
		out := outputAggregatetuple{}
		out.Fill(db, tuple)
		db.event.Aggregatetuples = append(db.event.Aggregatetuples, out)
	case TesttupleType:
		tuple, err := db.GetTesttuple(tupleKey)
		if err != nil {
			return err
		}
		out := outputTesttuple{}
		out.Fill(db, tuple)
		db.event.Testtuples = append(db.event.Testtuples, out)
	}
	return nil
}

// AddComputePlanEvent add the compute plan matching the ID to the event struct
func (db *LedgerDB) AddComputePlanEvent(ComputePlanKey, status string, ModelsToDelete []string) error {
	if db.event == nil {
		db.event = &Event{}
	}
	cp := eventComputePlan{
		ComputePlanKey: ComputePlanKey,
		Status:         status,
	}
	algokeys, err := db.GetIndexKeys("algo~computeplankey~key", []string{"algo", ComputePlanKey})
	if err != nil {
		return err
	}
	cp.AlgoKeys = algokeys
	cp.ModelsToDelete = ModelsToDelete
	db.event.ComputePlans = append(db.event.ComputePlans, cp)
	return nil
}
