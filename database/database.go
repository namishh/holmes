package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseStore struct {
	DB *sql.DB
}

func GetConnection(dbName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	if db != nil {
		return db, nil
	}

	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %s", err)
	}

	log.Println("Connected Successfully to the Database")

	return db, nil
}

func CreateMigrations(DBName string, DB *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) NOT NULL,
    level INT DEFAULT 0,
		password VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS question (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    question TEXT,
    answer TEXT
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	return nil
}

func NewDatabaseStore(path string) (DatabaseStore, error) {
	DB, err := GetConnection(path)
	if err != nil {
		return DatabaseStore{}, err
	}

	if err := CreateMigrations(path, DB); err != nil {
		return DatabaseStore{}, err
	}

	return DatabaseStore{DB: DB}, nil
}
