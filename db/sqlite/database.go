package sqlite

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

// TODO: tests

// Database is a database instance for sqlite
type Database struct {
	DB *sqlx.DB
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
		// this is used for the demo page
		if os.Getenv("MOUTHFUL_PERIODIC_DELETE") == "enabled" {
			periodicWipe(d)
		}
	} else {
		d, err := sqlx.Connect("sqlite3", *databaseConfig.Database)
		if err != nil {
			return nil, err
		}
		db = d
	}
	DB := Database{
		DB: db,
	}
	err = DB.InitializeDatabase()
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

func periodicWipe(db *sqlx.DB) {
	ticker := time.NewTicker(24 * time.Hour)
	quit := make(chan struct{})
	log.Println("MOUTHFUL WILL DELETE ALL DATA ONCE EVERY 24 HOURS")
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("wiping data")
				_, err := db.Exec("DELETE FROM COMMENT")
				if err != nil {
					log.Println(err)
				} else {
					log.Println("data wiped")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
