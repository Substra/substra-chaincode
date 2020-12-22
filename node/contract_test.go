package node

import (
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/mock"
	"github.com/substrafoundation/chaincode/v2/database"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) PutState(key string, data []byte) error {
	args := m.Called(key, data)
	return args.Error(0)
}

func (m *MockDatabase) GetState(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Get(1).(error)
}

type MockedContext struct {
	mock.Mock
}

func (m *MockedContext) GetStub() shim.ChaincodeStubInterface {
	args := m.Called()
	return args.Get(0).(shim.ChaincodeStubInterface)
}

func (m *MockedContext) GetClientIdentity() cid.ClientIdentity {
	args := m.Called()
	return args.Get(0).(cid.ClientIdentity)
}

func TestRegistration(t *testing.T) {
	mockDB := new(MockDatabase)
	factory := func(_ interface{}) (database.Database, error) {
		return mockDB, nil
	}

	mockDB.On("PutState", "uuid1", mock.Anything).Return(nil).Once()

	contract := NewNodeContract(factory)

	ctx := new(MockedContext)
	contract.RegisterNode(ctx, "uuid1")
}
