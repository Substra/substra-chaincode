package node

import (
	"encoding/json"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/substrafoundation/chaincode/v2/database"
)

// SmartContract manages nodes
type SmartContract struct {
	contractapi.Contract
	dbFactory database.Factory
}

func NewNodeContract(factory database.Factory) *SmartContract {
	return &SmartContract{
		dbFactory: factory,
	}
}

// Node stores informations about node registered into the network,
// would be used to list authorized nodes for permissions
type Node struct {
	ID string `json:"id"`
}

// RegisterNode creates a new node in world state
func (s *SmartContract) RegisterNode(ctx contractapi.TransactionContextInterface, id string) error {
	node := Node{ID: id}

	nodeAsBytes, _ := json.Marshal(node)

	log.Println("register", node)

	db, err := s.dbFactory(ctx)

	if err != nil {
		log.Println("Failed to get DB")
		return err
	}

	return db.PutState(id, nodeAsBytes)
}

// QueryNodes will return all known nodes
func (s *SmartContract) QueryNodes(ctx contractapi.TransactionContextInterface) ([]Node, error) {
	startKey := ""
	endKey := ""

	// database, err := s.dbFactory(ctx)
	// then iter := database.GetAll() or whatever
	// the example here will get ALL keys so it's not a real use case, we need to filter by asset type
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []Node{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		node := new(Node)
		_ = json.Unmarshal(queryResponse.Value, node)

		results = append(results, *node)
	}

	return results, nil
}
