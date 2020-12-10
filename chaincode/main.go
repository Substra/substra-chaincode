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
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/sirupsen/logrus"
)

// SubstraChaincode is a Receiver for Chaincode shim functions
type SubstraChaincode struct {
}

// Create a global logger for the chaincode. Its default level is Info
var logger = logrus.New()

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
// TODO!!!!
func (t *SubstraChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 1 {
		return shim.Error("Incorrect arguments. Expecting nothing...")
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode.
func (t *SubstraChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	resp, _, err := t._Invoke(stub)

	if err != nil {
		return formatErrorResponse(err)
	}
	return shim.Success(resp)
}

func (t *SubstraChaincode) _Invoke(stub shim.ChaincodeStubInterface) ([]byte, *Event, error) {
	start := time.Now()
	// Log all input for potential debug later on.
	logger.Infof("[%s][%s] Args received: '%s'", stub.GetChannelID(), stub.GetTxID()[:10], stub.GetStringArgs())

	// Seed with a timestamp from the channel header so the chaincode's output
	// stay determinist for each transaction. It's necessary because endorsers
	// will compare their own output to the proposal.
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return nil, nil, err
	}
	seedTime := time.Unix(timestamp.GetSeconds(), int64(timestamp.GetNanos()))
	rand.Seed(seedTime.UnixNano())

	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	db := NewLedgerDB(stub)

	var result interface{}
	hasBookmark := false
	var bookmark string

	switch fn {
	case "createComputePlan":
		result, err = createComputePlan(db, args)
	case "createTesttuple":
		result, err = createTesttuple(db, args)
	case "createTraintuple":
		result, err = createTraintuple(db, args)
	case "createCompositeTraintuple":
		result, err = createCompositeTraintuple(db, args)
	case "createAggregatetuple":
		result, err = createAggregatetuple(db, args)
	case "cancelComputePlan":
		result, err = cancelComputePlan(db, args)
	case "logFailTest":
		result, err = logFailTest(db, args)
	case "logFailTrain":
		result, err = logFailTrain(db, args)
	case "logFailCompositeTrain":
		result, err = logFailCompositeTrain(db, args)
	case "logFailAggregate":
		result, err = logFailAggregate(db, args)
	case "logStartTest":
		result, err = logStartTest(db, args)
	case "logStartTrain":
		result, err = logStartTrain(db, args)
	case "logStartCompositeTrain":
		result, err = logStartCompositeTrain(db, args)
	case "logStartAggregate":
		result, err = logStartAggregate(db, args)
	case "logSuccessTest":
		result, err = logSuccessTest(db, args)
	case "logSuccessTrain":
		result, err = logSuccessTrain(db, args)
	case "logSuccessCompositeTrain":
		result, err = logSuccessCompositeTrain(db, args)
	case "logSuccessAggregate":
		result, err = logSuccessAggregate(db, args)
	case "queryAlgo":
		result, err = queryAlgo(db, args)
	case "queryAlgos":
		result, bookmark, err = queryAlgos(db, args)
		hasBookmark = true
	case "queryCompositeAlgo":
		result, err = queryCompositeAlgo(db, args)
	case "queryCompositeAlgos":
		result, bookmark, err = queryCompositeAlgos(db, args)
		hasBookmark = true
	case "queryAggregateAlgo":
		result, err = queryAggregateAlgo(db, args)
	case "queryAggregateAlgos":
		result, bookmark, err = queryAggregateAlgos(db, args)
		hasBookmark = true
	case "queryDataManager":
		result, err = queryDataManager(db, args)
	case "queryDataManagers":
		result, bookmark, err = queryDataManagers(db, args)
		hasBookmark = true
	case "queryDataSamples":
		result, bookmark, err = queryDataSamples(db, args)
		hasBookmark = true
	case "queryDataset":
		result, err = queryDataset(db, args)
	case "queryFilter":
		result, err = queryFilter(db, args)
	case "queryModelDetails":
		result, err = queryModelDetails(db, args)
	case "queryModelPermissions":
		result, err = queryModelPermissions(db, args)
	case "queryModels":
		result, bookmark, err = queryModels(db, args)
		hasBookmark = true
	case "queryObjective":
		result, err = queryObjective(db, args)
	case "queryObjectiveLeaderboard":
		result, err = queryObjectiveLeaderboard(db, args)
	case "queryObjectives":
		result, bookmark, err = queryObjectives(db, args)
		hasBookmark = true
	case "queryTesttuple":
		result, err = queryTesttuple(db, args)
	case "queryTesttuples":
		result, bookmark, err = queryTesttuples(db, args)
		hasBookmark = true
	case "queryTraintuple":
		result, err = queryTraintuple(db, args)
	case "queryCompositeTraintuple":
		result, err = queryCompositeTraintuple(db, args)
	case "queryAggregatetuple":
		result, err = queryAggregatetuple(db, args)
	case "queryTraintuples":
		result, bookmark, err = queryTraintuples(db, args)
		hasBookmark = true
	case "queryCompositeTraintuples":
		result, bookmark, err = queryCompositeTraintuples(db, args)
		hasBookmark = true
	case "queryAggregatetuples":
		result, bookmark, err = queryAggregatetuples(db, args)
		hasBookmark = true
	case "queryComputePlan":
		result, err = queryComputePlan(db, args)
	case "queryComputePlans":
		result, bookmark, err = queryComputePlans(db, args)
		hasBookmark = true
	case "registerAlgo":
		result, err = registerAlgo(db, args)
	case "registerCompositeAlgo":
		result, err = registerCompositeAlgo(db, args)
	case "registerAggregateAlgo":
		result, err = registerAggregateAlgo(db, args)
	case "registerDataManager":
		result, err = registerDataManager(db, args)
	case "registerDataSample":
		result, err = registerDataSample(db, args)
	case "registerObjective":
		result, err = registerObjective(db, args)
	case "updateComputePlan":
		result, err = updateComputePlan(db, args)
	case "updateDataManager":
		result, err = updateDataManager(db, args)
	case "updateDataSample":
		result, err = updateDataSample(db, args)
	case "registerNode":
		result, err = registerNode(db, args)
	case "queryNodes":
		result, err = queryNodes(db, args)
	default:
		err = errors.BadRequest("function \"%s\" not implemented", fn)
	}

	// Invoke duration
	duration := int(time.Since(start).Nanoseconds()) / 1e6
	// Return the result as success payload
	if err != nil {
		logger.Errorf("[%s][%s] Response (%dms): '%#v','%s' - Error: '%s'", stub.GetChannelID(), stub.GetTxID()[:10], duration, result, bookmark, err)
		return nil, nil, err
	}

	// Add bookmark (if any) to response
	if hasBookmark {
		result = map[string]interface{}{
			"results":  result,
			"bookmark": bookmark,
		}
	}

	// Marshal to json the smartcontract result
	resp, err := json.Marshal(result)

	if err != nil {
		logger.Infof("[%s][%s] Response (%dms): '%#v','%s'", stub.GetChannelID(), stub.GetTxID()[:10], duration, result, bookmark)
		return nil, nil, errors.Internal("could not format response: %s", err.Error())
	}

	// Log with no errors
	logger.Infof("[%s][%s] Response (%dms): '%s'", stub.GetChannelID(), stub.GetTxID()[:10], duration, resp)

	// Send event if there is any. It's done in one batch since we can only send
	// one event per call
	err = db.SendEvent()
	if err != nil {
		return nil, nil, errors.Internal("could not send event: %s", err.Error())
	}

	return resp, db.event, nil
}

func formatErrorResponse(err error) peer.Response {
	e := errors.Wrap(err)
	status := e.HTTPStatusCode()

	errStruct := map[string]interface{}{
		"error": e.Error(),
		// Serialize status in the message until fabric-sdk-py allows subtrabac to
		// access the status
		"status": status,
	}
	for k, v := range e.GetContext() {
		errStruct[k] = v
	}

	payload, _ := json.Marshal(errStruct)
	return peer.Response{
		Message: string(payload),
		Payload: payload,
		Status:  int32(status),
	}
}

var scc = new(SubstraChaincode)
var eventIndex = 1
var allEvents = make(map[int]*Event)

type eventsRequest struct {
	Identity string `json:"identity"`
	EventID  int    `json:"event_id"`
}

type invokeRequest struct {
	Identity string `json:"identity"`
	Function string `json:"function"`
	Args     string `json:"args"`
}

func handleError(w http.ResponseWriter, returnCode int, err error) {
	w.WriteHeader(returnCode)
	fmt.Fprintf(w, "%v\n", err)
}

func handleHealth(w http.ResponseWriter, req *http.Request) {
	// logger.Infof("Readiness: %v", req.RequestURI)
	fmt.Fprintf(w, "OK")
}

var cc = NewMockStub("standalone-chaincode", scc)

func handleInvoke(w http.ResponseWriter, req *http.Request) {

	logger.Infof("Request: %v", req.RequestURI)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, 500, err)
		return
	}

	invokeReq := invokeRequest{}
	err = json.Unmarshal(body, &invokeReq)

	if err != nil {
		handleError(w, 500, err)
		return
	}

	args := make([][]byte, 2)

	arr := make([]string, 2)
	json.Unmarshal(body, &arr)

	logger.Infof("Function: %v", invokeReq.Function)
	logger.Infof("Arguments: %v", invokeReq.Args)

	args[0] = []byte(invokeReq.Function)
	args[1] = []byte(invokeReq.Args)

	cc.Creator = invokeReq.Identity
	cc.args = args
	cc.MockTransactionStart(mockTxID)
	resp, events, err := scc._Invoke(cc)

	if events != nil {
		allEvents[eventIndex] = events
		eventIndex++
	}
	cc.MockTransactionEnd(mockTxID)

	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	fmt.Fprintf(w, "%s", resp)
}

func handleEvents(w http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		handleError(w, 500, err)
		return
	}

	eventsReq := eventsRequest{}
	err = json.Unmarshal(body, &eventsReq)

	if err != nil {
		handleError(w, 500, err)
		return
	}

	fmt.Printf("Identity: %v\n", eventsReq.Identity)
	fmt.Printf("requestedEvent: %v\n", eventsReq.EventID)
	cc.Creator = eventsReq.Identity

	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	event, ok := allEvents[eventsReq.EventID]

	if !ok {
		handleError(w, 404, fmt.Errorf("event not found: %d", eventsReq.EventID))
		return
	}

	evt, err := json.Marshal(event)

	if err != nil {
		handleError(w, 500, err)
		return
	}

	fmt.Fprintf(w, "%s", evt)
}

func main() {
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	// logger.Infof("Load TLS certificates")

	// key, err := ioutil.ReadFile(os.Getenv("TLS_KEY_FILE"))
	// if err != nil {
	// 	logger.Errorf("Cannot read key file: %s", err)
	// }

	// cert, err := ioutil.ReadFile(os.Getenv("TLS_CERT_FILE"))
	// if err != nil {
	// 	logger.Errorf("Cannot read cert file: %s", err)
	// }

	// ca, err := ioutil.ReadFile(os.Getenv("TLS_ROOTCERT_FILE"))
	// if err != nil {
	// 	logger.Errorf("Cannot read ca cert file: %s", err)
	// }

	// server := &shim.ChaincodeServer{
	// 	CCID:    os.Getenv("CHAINCODE_CCID"),
	// 	Address: os.Getenv("CHAINCODE_ADDRESS"),
	// 	CC:      new(SubstraChaincode),
	// 	TLSProps: shim.TLSProperties{
	// 		Disabled:      true,
	// 		// Key:           key,
	// 		// Cert:          cert,
	// 		// ClientCACerts: ca,
	// 	},
	// }

	// Start the chaincode external server

	port := 8080

	logger.Infof("Start  Substra ChaincodeServer on port %v", port)

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/invoke", handleInvoke)
	http.HandleFunc("/events", handleEvents)

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)

	if err != nil {
		logger.Errorf("Error starting SubstraChaincode chaincode: %s", err)
	}
}
