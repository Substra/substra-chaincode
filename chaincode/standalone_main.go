package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var scc = new(SubstraChaincode)
var cc = NewStandaloneStub("standalone-chaincode", scc)
var eventIndex = 1
var allEvents = make(map[int]*Event)
var invokeMutex sync.Mutex

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

func startStandaloneServer(port int) {
	logger.Infof("Start  Substra ChaincodeServer on port %v", port)

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/invoke", handleInvoke)
	http.HandleFunc("/events", handleEvents)

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		logger.Errorf("Error starting standalone chaincode server: %s", err)
	}
}

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

	logger.Infof("Function: %v", invokeReq.Function)
	logger.Infof("Arguments: %v", invokeReq.Args)

	var args [][]byte

	if invokeReq.Args == "" {
		args = [][]byte{
			[]byte(invokeReq.Function),
		}
	} else {
		args = [][]byte{
			[]byte(invokeReq.Function),
			[]byte(invokeReq.Args),
		}
	}

	resp, err := doInvoke(invokeReq.Identity, args)

	if err != nil {
		resp, status := _formatErrorResponse(err)
		w.WriteHeader(status)
		fmt.Fprintf(w, "%s", resp)
		return
	}

	fmt.Fprintf(w, "%s", resp)
}

func doInvoke(identity string, args [][]byte) ([]byte, error) {
	invokeMutex.Lock()
	cc.Creator = identity
	cc.args = args
	cc.MockTransactionStart(standaloneMockTxID)

	resp, events, err := scc._Invoke(cc)
	if events != nil {
		allEvents[eventIndex] = events
		eventIndex++
	}

	cc.MockTransactionEnd(standaloneMockTxID)
	invokeMutex.Unlock()

	return resp, err
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
