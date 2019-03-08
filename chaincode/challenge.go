package main

import (
	"fmt"
	"strings"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gopkg.in/go-playground/validator.v9"
)

// Set is a method of the receiver Challenge. It checks the validity of inputChallenge and uses its fields to set the Challenge.
// Returns the challengeKey and the datasetKey associated to test data
func (challenge *Challenge) Set(stub shim.ChaincodeStubInterface, inp inputChallenge) (challengeKey string, datasetKey string, err error) {
	// checking validity of submitted fields
	validate := validator.New()
	if err = validate.Struct(inp); err != nil {
		err = fmt.Errorf("invalid challenge inputs %s", err.Error())
		return
	}
	datasetKey = strings.Split(inp.TestData, ":")[0]
	dataKeys := strings.Split(strings.Replace(strings.Split(inp.TestData, ":")[1], " ", "", -1), ",")
	testOnly, _, err := checkSameDataset(stub, datasetKey, dataKeys)
	if err != nil {
		err = fmt.Errorf("invalid test data %s", err.Error())
		return
	} else if !testOnly {
		err = fmt.Errorf("test data are not tagged as testOnly data")
		return
	}
	challenge.TestData = &DatasetData{
		DatasetKey: datasetKey,
		DataKeys:   dataKeys,
	}
	challenge.Name = inp.Name
	challenge.DescriptionStorageAddress = inp.DescriptionStorageAddress
	challenge.Metrics = &HashDressName{
		Name:           inp.MetricsName,
		Hash:           inp.MetricsHash,
		StorageAddress: inp.MetricsStorageAddress,
	}
	owner, err := getTxCreator(stub)
	if err != nil {
		return
	}
	challenge.Owner = owner
	challenge.Permissions = inp.Permissions
	challengeKey = inp.DescriptionHash
	return
}

// -------------------------------------------------------------------------------------------
// Smart contract related to challengess
// -------------------------------------------------------------------------------------------

// registerChallenge stores a new challenge in the ledger.
// If the key exists, it will override the value with the new one
func registerChallenge(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	expectedArgs := getFieldNames(&inputChallenge{})
	if nbArgs := len(expectedArgs); nbArgs != len(args) {
		return nil, fmt.Errorf("incorrect arguments, expecting %d args: %s", nbArgs, strings.Join(expectedArgs, ", "))
	}

	// convert input strings args to input struct inputChallenge
	inpc := inputChallenge{}
	stringToInputStruct(args, &inpc)
	// check validity of input args and convert it to Challenge
	challenge := Challenge{}
	challengeKey, datasetKey, err := challenge.Set(stub, inpc)
	if err != nil {
		return nil, err
	}
	// check challenge is not already in ledger
	if elementBytes, _ := stub.GetState(challengeKey); elementBytes != nil {
		return nil, fmt.Errorf("challenge with this description already exists - %s", string(elementBytes))
	}
	// submit to ledger
	challengeBytes, _ := json.Marshal(challenge)
	if err := stub.PutState(challengeKey, challengeBytes); err != nil {
		return nil, fmt.Errorf("failed to submit to ledger the challenge with key %s, error is %s", challengeKey, err.Error())
	}
	// create composite key
	if err := createCompositeKey(stub, "challenge~owner~key", []string{"challenge", challenge.Owner, challengeKey}); err != nil {
		return nil, err
	}
	// add challenge to dataset
	err = addChallengeDataset(stub, datasetKey, challengeKey)
	// return []byte(challengeKey), err
	return []byte(challengeKey), err
}

// -------------------------------------------------------------------------------------------
// Utils for challengess
// -------------------------------------------------------------------------------------------

// addChallengeDataset associates a challenge to a dataset, more precisely, it adds the challenge key to the dataset
func addChallengeDataset(stub shim.ChaincodeStubInterface, datasetKey string, challengeKey string) error {
	dataset := Dataset{}
	if err := getElementStruct(stub, datasetKey, &dataset); err != nil {
		return nil
	}
	if dataset.ChallengeKey != "" {
		return fmt.Errorf("dataset is already associated with a challenge")
	}
	dataset.ChallengeKey = challengeKey
	datasetBytes, _ := json.Marshal(dataset)
	return stub.PutState(datasetKey, datasetBytes)
}
