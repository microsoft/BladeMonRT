package store

import (
	"database/sql"
	"os"
	"log"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

// all Configs are persisted in the pair of name and value, plus the type of the value
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
)

type ConfigStore struct {
	db *sql.DB
	configTableName string
}

func NewConfigStore() *ConfigStore {
	_, err := os.OpenFile("sqlite-database.db", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	var configStore *ConfigStore = &ConfigStore{db: sqliteDatabase, configTableName: "ConfigTable"}  
	return configStore
}

func (store *ConfigStore) InitTable() {
	log.Println("Create table...")
	statement, err := store.db.Prepare(fmt.Sprintf(TABLE_CREATE_QUERY, store.configTableName))
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}


/** Sets or adds a config name value pair. */
func (store *ConfigStore) SetConfigValue(configName string, configValue string) {
	log.Println("Setting config name: %s with value: %s", configName, configValue)
	statement, err := store.db.Prepare(fmt.Sprintf(INSERT_OR_REPLACE_QUERY, store.configTableName))
	if err != nil {
		log.Fatal(err.Error())
	}
	// It only supports string type (type==0) for now. Can be extented in future.
	statement.Exec(configName, configValue, 0)
	return
}