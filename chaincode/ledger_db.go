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
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// State is a in-memory representation of the db state
type State struct {
	items map[string]([]byte)
}

// LedgerDB to access the chaincode database during the lifetime of a SmartContract
type LedgerDB struct {
	cc               shim.ChaincodeStubInterface
	transactionState State
	mutex            *sync.RWMutex
}

// NewLedgerDB create a new db to access the chaincode during a SmartContract
func NewLedgerDB(stub shim.ChaincodeStubInterface) LedgerDB {
	return LedgerDB{
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

// gettransactionState returns a copy of an object that has been updated or created during the transaction
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
			return errors.NotFound(err)
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
		return fmt.Errorf("cannot create index %s: %s", index, err.Error())
	}
	value := []byte{0x00}
	if err = db.cc.PutState(compositeKey, value); err != nil {
		return fmt.Errorf("cannot create index %s: %s", index, err.Error())
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
		return nil, fmt.Errorf("get index %s failed: %s", index, err.Error())
	}
	defer iterator.Close()
	for iterator.HasNext() {
		compositeKey, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		_, keyParts, err := db.cc.SplitCompositeKey(compositeKey.Key)
		if err != nil {
			return nil, fmt.Errorf("get index %s failed: cannot split key %s: %s", index, compositeKey.Key, err.Error())
		}
		keys = append(keys, keyParts[len(keyParts)-1])
	}
	return keys, nil
}

// ----------------------------------------------
// High-level functions
// ----------------------------------------------

// GetAssetType fetch a object in the chaincode db and return its type.
// It fails if it does not exists or if it does not have an `AssetType` field.
func (db *LedgerDB) GetAssetType(key string) (AssetType, error) {
	asset := struct {
		AssetType AssetType `json:"assetType"`
	}{}
	if err := db.Get(key, &asset); err != nil {
		return asset.AssetType, err
	}
	return asset.AssetType, nil
}

// GetGenericTraintuple fetches a regular/composite/aggregate tuple in the chaincode db and returns
// its common properties.
func (db *LedgerDB) GetGenericTraintuple(key string) (assetType AssetType, status string, err error) {
	asset := struct {
		Status    string    `json:"status"`
		AssetType AssetType `json:"assetType"`
	}{}
	if err := db.Get(key, &asset); err != nil {
		return asset.AssetType, asset.Status, err
	}
	return asset.AssetType, asset.Status, nil
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
	if err := db.Get(key, &traintuple); err != nil {
		return traintuple, err
	}
	if traintuple.AssetType != TraintupleType {
		return traintuple, errors.NotFound("traintuple %s not found", key)
	}
	return traintuple, nil
}

// GetCompositeTraintuple fetches a CompositeTraintuple from the ledger using its unique key
func (db *LedgerDB) GetCompositeTraintuple(key string) (CompositeTraintuple, error) {
	traintuple := CompositeTraintuple{}
	if err := db.Get(key, &traintuple); err != nil {
		return traintuple, err
	}
	if traintuple.AssetType != CompositeTraintupleType {
		return traintuple, errors.NotFound("composite traintuple %s not found", key)
	}
	return traintuple, nil
}

// GetAggregatetuple fetches a Aggregatetuple from the ledger using its unique key
func (db *LedgerDB) GetAggregatetuple(key string) (Aggregatetuple, error) {
	aggregatetuple := Aggregatetuple{}
	if err := db.Get(key, &aggregatetuple); err != nil {
		return aggregatetuple, err
	}
	if aggregatetuple.AssetType != AggregatetupleType {
		return aggregatetuple, errors.NotFound("aggregatetuple %s not found", key)
	}
	return aggregatetuple, nil
}

// CompositeModelType is a type of model generated by a composite traintuple
type CompositeModelType uint8

// Const representing the types models generated by a composite traintuple
const (
	HeadType CompositeModelType = iota
	TrunkType
)

// GetOutModelHashDress retrieves an out-Model from a traintuple key.
// Return an error if the traintupleKey was not found.
func (db *LedgerDB) GetOutModelHashDress(traintupleKey string, modelType CompositeModelType, allowedAssetTypes []AssetType) (*HashDress, error) {
	for _, assetType := range allowedAssetTypes {
		switch assetType {
		case CompositeTraintupleType:
			tuple, err := db.GetCompositeTraintuple(traintupleKey)
			if err != nil {
				continue
			}
			switch modelType {
			case HeadType:
				return tuple.OutHeadModel.OutModel, nil
			case TrunkType:
				return tuple.OutTrunkModel.OutModel, nil
			default:
				return nil, fmt.Errorf("GetOutModelHashDress: Unsupported composite model type %s", modelType)
			}

		case TraintupleType:
			tuple, err := db.GetTraintuple(traintupleKey)
			if err == nil {
				return tuple.OutModel, nil
			}
		case AggregatetupleType:
			tuple, err := db.GetAggregatetuple(traintupleKey)
			if err == nil {
				return tuple.OutModel, nil
			}
		default:
			return nil, fmt.Errorf("GetOutModelHashDress: Unsupported asset type %s", assetType)
		}
	}

	return nil, errors.NotFound(
		"GetOutModelHashDress: Could not find traintuple %s with key \"%s\". Allowed types: %v.",
		modelType,
		traintupleKey,
		allowedAssetTypes)
}

// GetTesttuple fetches a Testtuple from the ledger using its unique key
func (db *LedgerDB) GetTesttuple(key string) (Testtuple, error) {
	testtuple := Testtuple{}
	if err := db.Get(key, &testtuple); err != nil {
		return testtuple, err
	}
	if testtuple.AssetType != TesttupleType {
		return testtuple, errors.NotFound("testtuple %s not found", key)
	}
	return testtuple, nil
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
