// Package sqlite is responsible for sqlite database connections and initialization.
package sqlite

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"

	"github.com/jmoiron/sqlx"
	// We absolutely need the sqlite driver here, this whole package depends on it
	_ "github.com/mattn/go-sqlite3"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver"
)

// SqliteQueries represents a list of queries for initial table creation in sqlite
var SqliteQueries = []string{
	`CREATE TABLE IF NOT EXISTS Thread(
			Id BLOB PRIMARY KEY,
			CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP not null,
			Path varchar(1024) not null UNIQUE
		)`,
	`CREATE TABLE IF NOT EXISTS Comment(
			Id BLOB PRIMARY KEY,
			ThreadId BLOB not null,
			Body text not null,
			Author varchar(255) not null,
			Confirmed bool not null default false,
			CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP not null,
			ReplyTo BLOB default null,
			DeletedAt TIMESTAMP DEFAULT null,
			FOREIGN KEY(ThreadId) references Thread(Id)
		)`,
}

// ValidateConfig validates the config for sqlite
func ValidateConfig(config model.Database) error {
	err := ""
	if config.Database == nil {
		err += "Please specify the database file name in Database.Database"
	}
	if err != "" {
		return errors.New(err)
	}
	return nil
}

// CreateDatabase creates a database instance from the given config
func CreateDatabase(databaseConfig model.Database) (abstraction.Database, error) {
	err := ValidateConfig(databaseConfig)
	if err != nil {
		return nil, err
	}
	var db *sqlx.DB
	if *databaseConfig.Database == ":memory:" {
		d, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			return nil, err
		}
		db = d
	} else {
		err := CreateDirectoryIfNotExists(*databaseConfig.Database, afero.NewOsFs())
		if err != nil {
			return nil, err
		}
		d, err := sqlx.Connect("sqlite3", *databaseConfig.Database)
		if err != nil {
			return nil, err
		}
		db = d
	}
	DB := sqlxDriver.Database{
		DB:      db,
		Queries: SqliteQueries,
		Dialect: "sqlite3",
	}
	err = DB.InitializeDatabase()
	if err != nil {
		return &DB, err
	}

	// this is used for the demo page
	if os.Getenv("MOUTHFUL_PERIODIC_DELETE") == "enabled" {
		_, err := DB.CreateComment("Hello world!", "Mouthful", "/", true, nil)
		if err != nil {
			return nil, err
		}
		periodicWipe(DB)

	}

	return &DB, nil
}

// CreateTestDatabase creates a test database in memory
func CreateTestDatabase() abstraction.Database {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	DB := sqlxDriver.Database{
		DB:      db,
		Queries: SqliteQueries,
		Dialect: "sqlite3",
		IsTest:  true,
	}
	err = DB.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return &DB
}

// CreateDirectoryIfNotExists checks for the given directory and if not found - creates it. Used for sqlite database initialization
func CreateDirectoryIfNotExists(path string, fs afero.Fs) error {
	dir := filepath.Dir(path)
	if _, err := fs.Stat(dir); os.IsNotExist(err) {
		err = fs.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func periodicWipe(db sqlxDriver.Database) {
	ticker := time.NewTicker(24 * time.Hour)
	log.Println("MOUTHFUL WILL DELETE ALL DATA ONCE EVERY 24 HOURS")
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("wiping data")
				_, err := db.DB.Exec("DELETE FROM COMMENT")
				if err != nil {
					log.Println(err)
				} else {
					log.Println("data wiped")
					log.Println("Will create comment in a sec...")
					// I'll assume a race condition here with the exec above.
					time.Sleep(time.Second * 1)
					log.Println("Creating comment")
					_, err := db.CreateComment("Hello world!", "Mouthful", "/", true, nil)
					if err != nil {
						log.Println(err)
					} else {
						log.Println("Comment created")
					}
				}
			}
		}
	}()
}
