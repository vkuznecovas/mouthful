package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

// Db is a database instance for sqlite
type Database struct {
	DB *sqlx.DB
}

func CreateDatabase() abstraction.Database {
	db, err := sqlx.Connect("sqlite3", "__deleteme.db")
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
