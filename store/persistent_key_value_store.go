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


/** This class is copy of PersistentKeyValueStore.py class in GO with the functionality of initializing a table, setting name-value pairs, and retrieving the value for a given name. */

// type == 0 for string
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
)

type PersistentKeyValueStoreInterface interface {
	InitTable() error
	SetConfigValue(configName string, configValue string) error
	GetConfigValue(configName string) (string, error)
}

type PersistentKeyValueStore struct {
	logger          *log.Logger
	db              *sql.DB
	configTableName string
}

func NewPersistentKeyValueStore(fileName string, tableName string) (*PersistentKeyValueStore, error) {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("PersistentKeyValueStore")

	sqliteDatabase, err := sql.Open("sqlite3", fileName)
	if err != nil {
		logger.Println("Error creating new config store.")
		return nil, err
	}

	var PersistentKeyValueStore *PersistentKeyValueStore = &PersistentKeyValueStore{logger: logger, db: sqliteDatabase, configTableName: tableName}
	return PersistentKeyValueStore, nil
}

func (store *PersistentKeyValueStore) InitTable() error {
	store.logger.Println("Create table...")
	statement, err := store.db.Prepare(fmt.Sprintf(TABLE_CREATE_QUERY, store.configTableName))
	if err != nil {
		store.logger.Println("Error initializing table.")
		return err
	}
	statement.Exec()
	return nil
}

/** Sets or adds a config name-value pair. */
func (store *PersistentKeyValueStore) SetConfigValue(configName string, configValue string) error {
	store.logger.Println(fmt.Sprintf("Setting config name: %s with value: %s", configName, configValue))
	statement, err := store.db.Prepare(fmt.Sprintf(INSERT_OR_REPLACE_QUERY, store.configTableName))
	if err != nil {
		store.logger.Println("Error setting config value.")
		return err
	}
	// It only supports string type (type==0) for now. Can be extented in future.
	statement.Exec(configName, configValue, 0)
	return nil
}

/** Get the value for a given config name. */
func (store *PersistentKeyValueStore) GetConfigValue(configName string) (string, error) {
	statement, err := store.db.Prepare(fmt.Sprintf(READ_WHERE_QUERY, store.configTableName))
	if err != nil {
		store.logger.Println(err.Error())
		return "", err
	}
	rows, err := statement.Query(configName)
	if err != nil {
		store.logger.Println(err.Error())
		return "", err
	}

	defer statement.Close()
	if rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			store.logger.Println(err.Error())
			return "", err
		}
		return value, nil
	}
	return "", errors.New(fmt.Sprintf("Row with Name=%s not found.", configName))
}
