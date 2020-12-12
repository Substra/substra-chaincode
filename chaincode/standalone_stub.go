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
	"container/list"
	"strings"
	"unicode/utf8"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/hyperledger/fabric-protos-go/msp"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/util"
	"github.com/pkg/errors"
)

const standaloneMockTxID = "fa0f757bc278fdf6a32d00975602eb853e23a86a156781588d99ddef5b80720f"
const standaloneWorker = "SampleOrg"

type substraChaincode interface {
	shim.Chaincode
	_Invoke() ([]byte, *Event, error)
}

// Stub is an implementation of ChaincodeStubInterface for the standalone chaincode.
type Stub struct {
	// The transaction creator
	Creator string

	// arguments the stub was called with
	args [][]byte

	// A pointer back to the chaincode that will invoke this, set by constructor.
	// If a peer calls this stub, the chaincode will be invoked from here.
	cc shim.Chaincode

	// A nice name that can be used for logging
	Name string

	// State keeps name value pairs
	State map[string][]byte

	// Keys stores the list of mapped values in lexical order
	Keys *list.List

	// registered list of other StandaloneStub chaincodes that can be called from this StandaloneStub
	Invokables map[string]*Stub

	// stores a transaction uuid while being Invoked / Deployed
	// TODO if a chaincode uses recursion this may need to be a stack of TxIDs or possibly a reference counting map
	TxID string

	TxTimestamp *timestamp.Timestamp

	// mocked signedProposal
	signedProposal *pb.SignedProposal

	// stores a channel ID of the proposal
	ChannelID string

	PvtState map[string]map[string][]byte

	// stores per-key endorsement policy, first map index is the collection, second map index is the key
	EndorsementPolicies map[string]map[string][]byte

	// channel to store ChaincodeEvents
	ChaincodeEventsChannel chan *pb.ChaincodeEvent

	Decorations map[string][]byte
}

// GetTxID ...
func (stub *Stub) GetTxID() string {
	return stub.TxID
}

// GetChannelID ...
func (stub *Stub) GetChannelID() string {
	return stub.ChannelID
}

// GetArgs ...
func (stub *Stub) GetArgs() [][]byte {
	return stub.args
}

// GetStringArgs ...
func (stub *Stub) GetStringArgs() []string {
	args := stub.GetArgs()
	strargs := make([]string, 0, len(args))
	for _, barg := range args {
		strargs = append(strargs, string(barg))
	}
	return strargs
}

// GetFunctionAndParameters ...
func (stub *Stub) GetFunctionAndParameters() (function string, params []string) {
	allargs := stub.GetStringArgs()
	function = ""
	params = []string{}
	if len(allargs) >= 1 {
		function = allargs[0]
		params = allargs[1:]
	}
	return
}

// MockTransactionStart ...
func (stub *Stub) MockTransactionStart(txid string) {
	stub.TxID = txid
	stub.setSignedProposal(&pb.SignedProposal{})
	stub.setTxTimestamp(util.CreateUtcTimestamp())
}

// MockTransactionEnd ...
func (stub *Stub) MockTransactionEnd(uuid string) {
	stub.signedProposal = nil
	stub.TxID = ""
}

// MockPeerChaincode ...
func (stub *Stub) MockPeerChaincode(invokableChaincodeName string, otherStub *Stub) {
	stub.Invokables[invokableChaincodeName] = otherStub
}

// MockInit ...
func (stub *Stub) MockInit(uuid string, args [][]byte) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	res := stub.cc.Init(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

// MockInvoke ...
func (stub *Stub) MockInvoke(args [][]byte) pb.Response {
	return stub.MockInvokeTxID(standaloneMockTxID, args)
}

// MockInvokeTxID ...
func (stub *Stub) MockInvokeTxID(txid string, args [][]byte) pb.Response {
	stub.args = args
	stub.MockTransactionStart(txid)
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(txid)
	return res
}

// GetDecorations ...
func (stub *Stub) GetDecorations() map[string][]byte {
	return stub.Decorations
}

// MockInvokeWithSignedProposal ...
func (stub *Stub) MockInvokeWithSignedProposal(uuid string, args [][]byte, sp *pb.SignedProposal) pb.Response {
	stub.args = args
	stub.MockTransactionStart(uuid)
	stub.signedProposal = sp
	res := stub.cc.Invoke(stub)
	stub.MockTransactionEnd(uuid)
	return res
}

// GetPrivateData ...
func (stub *Stub) GetPrivateData(collection string, key string) ([]byte, error) {
	m, in := stub.PvtState[collection]

	if !in {
		return nil, nil
	}

	return m[key], nil
}

// PutPrivateData ...
func (stub *Stub) PutPrivateData(collection string, key string, value []byte) error {
	m, in := stub.PvtState[collection]
	if !in {
		stub.PvtState[collection] = make(map[string][]byte)
		m, in = stub.PvtState[collection]
	}

	m[key] = value

	return nil
}

// DelPrivateData ...
func (stub *Stub) DelPrivateData(collection string, key string) error {
	return errors.New("Not Implemented")
}

// GetPrivateDataByRange ...
func (stub *Stub) GetPrivateDataByRange(collection, startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	return nil, errors.New("Not Implemented")
}

// GetPrivateDataByPartialCompositeKey ...
func (stub *Stub) GetPrivateDataByPartialCompositeKey(collection, objectType string, attributes []string) (shim.StateQueryIteratorInterface, error) {
	return nil, errors.New("Not Implemented")
}

// GetPrivateDataQueryResult ...
func (stub *Stub) GetPrivateDataQueryResult(collection, query string) (shim.StateQueryIteratorInterface, error) {
	return nil, errors.New("Not Implemented")
}

// GetState ...
func (stub *Stub) GetState(key string) ([]byte, error) {
	value := stub.State[key]
	return value, nil
}

// PutState ...
func (stub *Stub) PutState(key string, value []byte) error {
	if stub.TxID == "" {
		err := errors.New("cannot PutState without a transactions - call stub.MockTransactionStart()?")
		return err
	}

	// If the value is nil or empty, delete the key
	if len(value) == 0 {
		return stub.DelState(key)
	}

	stub.State[key] = value

	// insert key into ordered list of keys
	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		elemValue := elem.Value.(string)
		comp := strings.Compare(key, elemValue)
		if comp < 0 {
			// key < elem, insert it before elem
			stub.Keys.InsertBefore(key, elem)
			break
		} else if comp == 0 {
			// keys exists, no need to change
			break
		} else { // comp > 0
			// key > elem, keep looking unless this is the end of the list
			if elem.Next() == nil {
				stub.Keys.PushBack(key)
				break
			}
		}
	}

	// special case for empty Keys list
	if stub.Keys.Len() == 0 {
		stub.Keys.PushFront(key)
	}

	return nil
}

// DelState removes the specified `key` and its value from the ledger.
func (stub *Stub) DelState(key string) error {
	delete(stub.State, key)

	for elem := stub.Keys.Front(); elem != nil; elem = elem.Next() {
		if strings.Compare(key, elem.Value.(string)) == 0 {
			stub.Keys.Remove(elem)
		}
	}

	return nil
}

// GetStateByRange ...
func (stub *Stub) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	if err := validateSimpleKeys(startKey, endKey); err != nil {
		return nil, err
	}
	return NewStandaloneStateRange(stub, startKey, endKey, false, OutputPageSize), nil
}

// GetQueryResult function can be invoked by a chaincode to perform a
// rich query against state database.  Only supported by state database implementations
// that support rich query.  The query string is in the syntax of the underlying
// state database. An iterator is returned which can be used to iterate (next) over
// the query result set
func (stub *Stub) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	// Not implemented since the mock engine does not have a query engine.
	// However, a very simple query engine that supports string matching
	// could be implemented to test that the framework supports queries
	return nil, errors.New("not implemented")
}

// GetHistoryForKey function can be invoked by a chaincode to return a history of
// key values across time. GetHistoryForKey is intended to be used for read-only queries.
func (stub *Stub) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	return nil, errors.New("not implemented")
}

//GetStateByPartialCompositeKey function can be invoked by a chaincode to query the
//state based on a given partial composite key. This function returns an
//iterator which can be used to iterate over all composite keys whose prefix
//matches the given partial composite key. This function should be used only for
//a partial composite key. For a full composite key, an iter with empty response
//would be returned.
func (stub *Stub) GetStateByPartialCompositeKey(objectType string, attributes []string) (shim.StateQueryIteratorInterface, error) {
	partialCompositeKey, err := stub.CreateCompositeKey(objectType, attributes)
	if err != nil {
		return nil, err
	}
	return NewStandaloneStateRange(stub, partialCompositeKey, partialCompositeKey+string(utf8.MaxRune), false, OutputPageSize), nil
}

// CreateCompositeKey combines the list of attributes
//to form a composite key.
func (stub *Stub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return newCompositeKey(objectType, attributes)
}

// SplitCompositeKey splits the composite key into attributes
// on which the composite key was formed.
func (stub *Stub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	return splitCompositeKey(compositeKey)
}

// GetStateByRangeWithPagination ...
func (stub *Stub) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	return nil, nil, nil
}

// GetStateByPartialCompositeKeyWithPagination ...
func (stub *Stub) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
	pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	partialCompositeKey, err := stub.CreateCompositeKey(objectType, keys)
	if err != nil {
		return nil, nil, err
	}
	startKey := partialCompositeKey
	if bookmark != "" {
		startKey = bookmark
	}
	iterator := NewStandaloneStateRange(stub, startKey, partialCompositeKey+string(utf8.MaxRune), true, pageSize)
	return iterator, iterator.Metadata, nil
}

// GetQueryResultWithPagination ...
func (stub *Stub) GetQueryResultWithPagination(query string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	return nil, nil, nil
}

// InvokeChaincode calls a peered chaincode.
// E.g. stub1.InvokeChaincode("stub2Hash", funcArgs, channel)
// Before calling this make sure to create another StandaloneStub stub2, call stub2.MockInit(uuid, func, args)
// and register it with stub1 by calling stub1.MockPeerChaincode("stub2Hash", stub2)
func (stub *Stub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response {
	// Internally we use chaincode name as a composite name
	if channel != "" {
		chaincodeName = chaincodeName + "/" + channel
	}
	// TODO "args" here should possibly be a serialized pb.ChaincodeInput
	otherStub := stub.Invokables[chaincodeName]
	//	function, strings := getFuncArgs(args)
	res := otherStub.MockInvokeTxID(stub.TxID, args)
	return res
}

// TODO: delete when proper identity verification is implemented
const standaloneFakeCertificate = `
-----BEGIN CERTIFICATE-----
MIIEBDCCAuygAwIBAgIDAjppMA0GCSqGSIb3DQEBBQUAMEIxCzAJBgNVBAYTAlVT
MRYwFAYDVQQKEw1HZW9UcnVzdCBJbmMuMRswGQYDVQQDExJHZW9UcnVzdCBHbG9i
YWwgQ0EwHhcNMTMwNDA1MTUxNTU1WhcNMTUwNDA0MTUxNTU1WjBJMQswCQYDVQQG
EwJVUzETMBEGA1UEChMKR29vZ2xlIEluYzElMCMGA1UEAxMcR29vZ2xlIEludGVy
bmV0IEF1dGhvcml0eSBHMjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AJwqBHdc2FCROgajguDYUEi8iT/xGXAaiEZ+4I/F8YnOIe5a/mENtzJEiaB0C1NP
VaTOgmKV7utZX8bhBYASxF6UP7xbSDj0U/ck5vuR6RXEz/RTDfRK/J9U3n2+oGtv
h8DQUB8oMANA2ghzUWx//zo8pzcGjr1LEQTrfSTe5vn8MXH7lNVg8y5Kr0LSy+rE
ahqyzFPdFUuLH8gZYR/Nnag+YyuENWllhMgZxUYi+FOVvuOAShDGKuy6lyARxzmZ
EASg8GF6lSWMTlJ14rbtCMoU/M4iarNOz0YDl5cDfsCx3nuvRTPPuj5xt970JSXC
DTWJnZ37DhF5iR43xa+OcmkCAwEAAaOB+zCB+DAfBgNVHSMEGDAWgBTAephojYn7
qwVkDBF9qn1luMrMTjAdBgNVHQ4EFgQUSt0GFhu89mi1dvWBtrtiGrpagS8wEgYD
VR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAQYwOgYDVR0fBDMwMTAvoC2g
K4YpaHR0cDovL2NybC5nZW90cnVzdC5jb20vY3Jscy9ndGdsb2JhbC5jcmwwPQYI
KwYBBQUHAQEEMTAvMC0GCCsGAQUFBzABhiFodHRwOi8vZ3RnbG9iYWwtb2NzcC5n
ZW90cnVzdC5jb20wFwYDVR0gBBAwDjAMBgorBgEEAdZ5AgUBMA0GCSqGSIb3DQEB
BQUAA4IBAQA21waAESetKhSbOHezI6B1WLuxfoNCunLaHtiONgaX4PCVOzf9G0JY
/iLIa704XtE7JW4S615ndkZAkNoUyHgN7ZVm2o6Gb4ChulYylYbc3GrKBIxbf/a/
zG+FA1jDaFETzf3I93k9mTXwVqO94FntT0QJo544evZG0R0SnU++0ED8Vf4GXjza
HFa9llF7b1cq26KqltyMdMKVvvBulRP/F/A8rLIQjcxz++iPAsbw+zOzlTvjwsto
WHPbqCRiOwY1nQ2pM714A5AuTHhdUDqB1O6gyHA43LL5Z/qHQF1hwFGPa4NrzQU6
yuGnBXj8ytqU0CwIPX4WecigUCAkVDNx
-----END CERTIFICATE-----
`

// GetCreator ...
func (stub *Stub) GetCreator() ([]byte, error) {
	sid := &msp.SerializedIdentity{
		Mspid:   stub.Creator,
		IdBytes: []byte(standaloneFakeCertificate),
	}

	return proto.Marshal(sid)
}

// GetTransient ...
func (stub *Stub) GetTransient() (map[string][]byte, error) {
	return nil, nil
}

// GetBinding ...
func (stub *Stub) GetBinding() ([]byte, error) {
	return nil, nil
}

// GetSignedProposal ...
func (stub *Stub) GetSignedProposal() (*pb.SignedProposal, error) {
	return stub.signedProposal, nil
}

func (stub *Stub) setSignedProposal(sp *pb.SignedProposal) {
	stub.signedProposal = sp
}

// GetArgsSlice ...
func (stub *Stub) GetArgsSlice() ([]byte, error) {
	return nil, nil
}

// GetPrivateDataHash ...
func (stub *Stub) GetPrivateDataHash(collection, key string) ([]byte, error) {
	return nil, nil
}

func (stub *Stub) setTxTimestamp(time *timestamp.Timestamp) {
	// Using a sequential timestamp make the tests' output determinist.
	// We're using a sequence and not a fix value so the seed won't reset at each Invoke.
	stub.TxTimestamp.Seconds = stub.TxTimestamp.Seconds + 1
	stub.TxTimestamp.Nanos = stub.TxTimestamp.Nanos + 1
}

// GetTxTimestamp ...
func (stub *Stub) GetTxTimestamp() (*timestamp.Timestamp, error) {
	if stub.TxTimestamp == nil {
		return nil, errors.New("stub.TxTimestamp not set")
	}
	return stub.TxTimestamp, nil
}

// SetEvent ...
func (stub *Stub) SetEvent(name string, payload []byte) error {
	stub.ChaincodeEventsChannel <- &pb.ChaincodeEvent{EventName: name, Payload: payload}
	return nil
}

// SetStateValidationParameter ...
func (stub *Stub) SetStateValidationParameter(key string, ep []byte) error {
	return stub.SetPrivateDataValidationParameter("", key, ep)
}

// GetStateValidationParameter ...
func (stub *Stub) GetStateValidationParameter(key string) ([]byte, error) {
	return stub.GetPrivateDataValidationParameter("", key)
}

// SetPrivateDataValidationParameter ...
func (stub *Stub) SetPrivateDataValidationParameter(collection, key string, ep []byte) error {
	m, in := stub.EndorsementPolicies[collection]
	if !in {
		stub.EndorsementPolicies[collection] = make(map[string][]byte)
		m, in = stub.EndorsementPolicies[collection]
	}

	m[key] = ep
	return nil
}

// GetPrivateDataValidationParameter ...
func (stub *Stub) GetPrivateDataValidationParameter(collection, key string) ([]byte, error) {
	m, in := stub.EndorsementPolicies[collection]

	if !in {
		return nil, nil
	}

	return m[key], nil
}

// NewStandaloneStub initialises the internal State map
func NewStandaloneStub(name string, cc shim.Chaincode) *Stub {
	s := new(Stub)
	s.Creator = standaloneWorker
	s.Name = name
	s.cc = cc
	s.State = make(map[string][]byte)
	s.PvtState = make(map[string]map[string][]byte)
	s.EndorsementPolicies = make(map[string]map[string][]byte)
	s.Invokables = make(map[string]*Stub)
	s.Keys = list.New()
	s.ChaincodeEventsChannel = make(chan *pb.ChaincodeEvent, OutputPageSize+1) //define large capacity for non-blocking setEvent calls.
	s.Decorations = make(map[string][]byte)
	s.TxTimestamp = &timestamp.Timestamp{}

	return s
}

// NewStandaloneStubWithRegisterNode ...
func NewStandaloneStubWithRegisterNode(name string, cc shim.Chaincode) *Stub {
	s := NewStandaloneStub(name, cc)
	s.MockInvoke([][]byte{[]byte("registerNode")})

	return s
}

// StateRange ...
type StateRange struct {
	Closed      bool
	Stub        *Stub
	StartKey    string
	EndKey      string
	Current     *list.Element
	IsPaginated bool
	PageSize    int32
	Metadata    *pb.QueryResponseMetadata
}

// HasNext returns true if the range query iterator contains additional keys
// and values.
func (iter *StateRange) HasNext() bool {
	if iter.Closed {
		// previously called Close()
		return false
	}

	if iter.Current == nil {
		return false
	}

	current := iter.Current
	for current != nil {
		// if this is an open-ended query for all keys, return true
		if iter.StartKey == "" && iter.EndKey == "" {
			return true
		}
		// iterator has already yielded enough results
		if iter.IsPaginated && iter.Metadata.FetchedRecordsCount == iter.PageSize {
			return false
		}
		comp1 := strings.Compare(current.Value.(string), iter.StartKey)
		comp2 := strings.Compare(current.Value.(string), iter.EndKey)
		if comp1 >= 0 {
			if comp2 < 0 {
				return true
			}
			iter.Metadata.Bookmark = ""
			return false
		}
		current = current.Next()
	}

	// we've reached the end of the underlying values
	return false
}

// Next returns the next key and value in the range query iterator.
func (iter *StateRange) Next() (*queryresult.KV, error) {
	if iter.Closed == true {
		err := errors.New("StandaloneStateRange.Next() called after Close()")
		return nil, err
	}

	if iter.HasNext() == false {
		err := errors.New("StandaloneStateRange.Next() called when it does not HaveNext()")
		return nil, err
	}

	if iter.Current != nil && iter.IsPaginated && iter.Metadata.FetchedRecordsCount >= iter.PageSize {
		return nil, errors.New("Paginated StandaloneStateRange.Next() went past end of range")
	}

	for iter.Current != nil {
		comp1 := strings.Compare(iter.Current.Value.(string), iter.StartKey)
		comp2 := strings.Compare(iter.Current.Value.(string), iter.EndKey)
		// compare to start and end keys. or, if this is an open-ended query for
		// all keys, it should always return the key and value
		if (comp1 >= 0 && comp2 < 0) || (iter.StartKey == "" && iter.EndKey == "") {
			key := iter.Current.Value.(string)
			value, err := iter.Stub.GetState(key)
			iter.Current = iter.Current.Next()
			if iter.Current != nil {
				iter.Metadata.Bookmark = iter.Current.Value.(string)
			} else {
				iter.Metadata.Bookmark = ""
			}
			iter.Metadata.FetchedRecordsCount++
			return &queryresult.KV{Key: key, Value: value}, err
		}
		iter.Current = iter.Current.Next()
	}
	err := errors.New("StandaloneStateRange.Next() went past end of range")
	return nil, err
}

// Close closes the range query iterator. This should be called when done
// reading from the iterator to free up resources.
func (iter *StateRange) Close() error {
	if iter.Closed == true {
		err := errors.New("StandaloneStateRange.Close() called after Close()")
		return err
	}

	iter.Closed = true
	return nil
}

// Print ...
func (iter *StateRange) Print() {
}

// NewStandaloneStateRange ...
func NewStandaloneStateRange(stub *Stub, startKey string, endKey string, isPaginated bool, pageSize int32) *StateRange {
	iter := new(StateRange)
	iter.Closed = false
	iter.Stub = stub
	iter.StartKey = startKey
	iter.EndKey = endKey
	iter.Current = stub.Keys.Front()
	iter.IsPaginated = isPaginated
	iter.PageSize = pageSize
	iter.Metadata = &pb.QueryResponseMetadata{}

	iter.Print()

	return iter
}
