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
	stmt := `CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) NOT NULL,
    	points INT DEFAULT 0,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255) UNIQUE NOT NULL,
		last_answered_question TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	question TEXT,
     	answer TEXT,
      	title TEXT,
       	points INT
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS hints  (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
     	hint TEXT,
      	worth INT,
       parent_question_id INT,
       	FOREIGN KEY (parent_question_id) REFERENCES questions(id)
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	path TEXT,
     	parent_question_id INTEGER,
     	FOREIGN KEY (parent_question_id) REFERENCES questions(id)
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS audios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	path TEXT,
     	parent_question_id INTEGER,
     	FOREIGN KEY(parent_question_id) REFERENCES questions(id)
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	path TEXT,
     	parent_question_id INTEGER,
     	FOREIGN KEY(parent_question_id) REFERENCES questions(id)
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS team_completed_questions (
    team_id INTEGER,
    question_id INTEGER,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (team_id, question_id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
    );`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS team_hint_unlocked (
    team_id INTEGER,
    hint_id INTEGER,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (team_id, hint_id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    FOREIGN KEY (hint_id) REFERENCES hints(id)
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
