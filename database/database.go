// Package database holds everything related to data persistence
package database

// Factory is capable of deriving a database from the context
type Factory = func(ctx interface{}) (Database, error) // ctx should be a ctx contractapi.TransactionContextInterface, but that would mean coupling

// Database is the main interface to act on the persistence layer
// This covers all CRUD operations
type Database interface {
	DBWriter
	DBReader
}

// DBWriter handles persisting and updating data
type DBWriter interface {
	PutState(key string, data []byte) error
}

type DBReader interface {
	GetState(key string) ([]byte, error)
}
