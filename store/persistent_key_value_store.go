package store

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/microsoft/BladeMonRT/logging"
	"log"
)

// PersistentKeyValueStoreInterface mock generation.
//go:generate mockgen -source=./persistent_key_value_store.go -destination=./mock_persistent_key_value_store.go -package=store

// Type == 0 for string
const (
	TABLE_CREATE_QUERY = `
		CREATE TABLE IF NOT EXISTS %s
		(
			Name                TEXT    PRIMARY KEY   NOT NULL,
			Value               TEXT                  NOT NULL,
			Type                INTEGER               NOT NULL
		);`
	INSERT_OR_REPLACE_QUERY = `INSERT OR REPLACE INTO %s (Name, Value, Type) VALUES(?,?,?);`
	READ_WHERE_QUERY        = `SELECT Value FROM %s WHERE Name = $1;`
	DELETE_ALL_QUERY = `DELETE FROM %s;`
)

/** Interface for the PersistentKeyValueStore that define which methods are implemented by PersistentKeyValueStore. */
type PersistentKeyValueStoreInterface interface {
	InitTable() error
	SetValue(key string, value string) error
	GetValue(key string) (string, error)
	Clear() error
}

/** This class is copy of PersistentKeyValueStore.py class in GO. It stores key value pairs persistently using a database. */ 
type PersistentKeyValueStore struct {
	logger          *log.Logger
	db              *sql.DB
	tableName string
}

func NewPersistentKeyValueStore(fileName string, tableName string) (*PersistentKeyValueStore, error) {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("PersistentKeyValueStore")

	sqliteDatabase, err := sql.Open("sqlite3", fileName)
	if err != nil {
		logger.Println("Error creating new config store:", err)
		return nil, err
	}

	var PersistentKeyValueStore *PersistentKeyValueStore = &PersistentKeyValueStore{logger: logger, db: sqliteDatabase, tableName: tableName}
	return PersistentKeyValueStore, nil
}

func (store *PersistentKeyValueStore) InitTable() error {
	store.logger.Println("Creating table.")
	statement, err := store.db.Prepare(fmt.Sprintf(TABLE_CREATE_QUERY, store.tableName))
	if err != nil {
		store.logger.Println("Error initializing table:", err)
		return err
	}
	statement.Exec()
	return nil
}

func (store *PersistentKeyValueStore) SetValue(key string, value string) error {
	store.logger.Println(fmt.Sprintf("Setting key: %s to value: %s", key, value))
	statement, err := store.db.Prepare(fmt.Sprintf(INSERT_OR_REPLACE_QUERY, store.tableName))
	if err != nil {
		store.logger.Println("Error setting config value:", err)
		return err
	}
	// It only supports string type (type==0) for now. Can be extented in future.
	statement.Exec(key, value, 0)
	return nil
}

func (store *PersistentKeyValueStore) GetValue(key string) (string, error) {
	statement, err := store.db.Prepare(fmt.Sprintf(READ_WHERE_QUERY, store.tableName))
	if err != nil {
		store.logger.Println(err)
		return "", err
	}
	rows, err := statement.Query(key)
	if err != nil {
		store.logger.Println(err)
		return "", err
	}

	defer rows.Close()
	if rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			store.logger.Println(err)
			return "", err
		}
		return value, nil
	}
	return "", errors.New(fmt.Sprintf("Key=%s not found in the store.", key))
}

func (store *PersistentKeyValueStore) Clear() error{
	store.logger.Println("Clear table.")
	statement, err := store.db.Prepare(fmt.Sprintf(DELETE_ALL_QUERY, store.tableName))
	if err != nil {
		store.logger.Println("Error clearing table:", err)
		return err
	}

	statement.Exec()
	return nil
}