package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

// TODO: tests

// Database is a database instance for sqlite
type Database struct {
	DB *sqlx.DB
}

// CreateDatabase creates a database instance from the given config
func CreateDatabase(databaseConfig model.Database) (abstraction.Database, error) {
	var db *sqlx.DB
	if databaseConfig.Database == ":memory:" {
		d, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			return nil, err
		}
		db = d
	} else {
		// TODO: this should come from a file as well
		d, err := sqlx.Connect("sqlite3", databaseConfig.Database)
		if err != nil {
			return nil, err
		}
		db = d
	}
	DB := Database{
		DB: db,
	}
	err := DB.InitializeDatabase()
	if err != nil {
		return &DB, err
	}
	return &DB, nil
}

// CreateTestDatabase creates a test database in memory
func CreateTestDatabase() abstraction.Database {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	DB := Database{
		DB: db,
	}
	err = DB.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return &DB
}
