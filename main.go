package main

import (
	"log"
	"os"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/sirupsen/logrus"
	"github.com/substrafoundation/chaincode/v2/ledger"
	"github.com/substrafoundation/chaincode/v2/node"
)

var logger = logrus.New()

func main() {

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	chaincode, err := contractapi.NewChaincode(node.NewNodeContract(ledger.GetLedgerFromContext))

	if err != nil {
		log.Fatal("Error create substra chaincode", err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic(err.Error())
	}
}
