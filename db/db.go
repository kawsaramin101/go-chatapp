package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitializeDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)

	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := filepath.Join(cwd, "db", "create_user_table.sql")

	err = CreateTables(filePath)
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatal(err)
	}
}

func CreateTables(filePath string) error {

	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		return err
	}
	return nil
}
