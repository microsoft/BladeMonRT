package store

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/microsoft/BladeMonRT/logging"
	"log"
)

/** This class is copy of ConfigStore.py class in GO with the functionality of initializing a table, setting name-value pairs, and retrieving the value for a given name. */

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

type ConfigStore struct {
	logger          *log.Logger
	db              *sql.DB
	configTableName string
}

func NewConfigStore(fileName string, tableName string) (*ConfigStore, error) {
	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("ConfigStore")

	sqliteDatabase, err := sql.Open("sqlite3", fileName)
	if err != nil {
		logger.Println("Error creating new config store.")
		return nil, err
	}

	var configStore *ConfigStore = &ConfigStore{logger: logger, db: sqliteDatabase, configTableName: tableName}
	return configStore, nil
}

func (store *ConfigStore) InitTable() error {
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
func (store *ConfigStore) SetConfigValue(configName string, configValue string) error {
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
func (store *ConfigStore) GetConfigValue(configName string) (string, error) {
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
