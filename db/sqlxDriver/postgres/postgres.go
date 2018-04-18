// Package postgres is responsible for postgres database connections and initialization.
package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	// We absolutely need the postgres driver here, this whole file depends on it
	_ "github.com/lib/pq"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver"
)

// PostgresQueries represents a list of queries for initial table creation in Postgres
var PostgresQueries = []string{
	`CREATE TABLE IF NOT EXISTS Thread(
			Id uuid PRIMARY KEY,
			CreatedAt TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) not null,
			Path varchar(255) not null UNIQUE
		)`,
	`CREATE TABLE IF NOT EXISTS Comment(
			Id uuid PRIMARY KEY,
			ThreadId uuid not null,
			Body text not null,
			Author varchar(255) not null,
			Confirmed bool not null default false,
			CreatedAt TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) not null,
			ReplyTo uuid default null,
			DeletedAt TIMESTAMP(6) NULL,
			FOREIGN KEY(ThreadId) references Thread(Id)
		)`,
}

// ValidateConfig validates the config for mysql
func ValidateConfig(config model.Database) error {
	err := ""
	if config.Database == nil {
		err += "Please specify the database name in config(Database.Database)"
	}
	if config.Username == nil {
		err += "Please specify the database username in config(Database.Username)"
	}
	if config.Password == nil {
		err += "Please specify the database password in config(Database.Password)"
	}
	if config.Host == nil {
		err += "Please specify the database host in config(Database.Host)"
	}
	if config.Database == nil {
		err += "Please specify the database name in config(Database.Database)"
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
	port := ""
	if databaseConfig.Port != nil {
		port = ":" + *databaseConfig.Port
	}
	connectionString := fmt.Sprintf("postgresql://%v:%v@%v%v/%v?connect_timeout=10", *databaseConfig.Username, *databaseConfig.Password, *databaseConfig.Host, port, *databaseConfig.Database)
	if databaseConfig.SSLEnabled != nil {
		if !*databaseConfig.SSLEnabled {
			connectionString += "&sslmode=disable"
		}
	}
	d, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	db.Mapper = reflectx.NewMapperTagFunc("db",
		nil,
		func(s string) string {
			return strings.ToLower(s)
		},
	)
	db = d
	DB := sqlxDriver.Database{
		DB:      db,
		Queries: PostgresQueries,
		Dialect: "postgres",
		IsTest:  false,
	}
	err = DB.InitializeDatabase()
	if err != nil {
		return &DB, err
	}
	return &DB, nil
}

// CreateTestDatabase creates a test database
func CreateTestDatabase() abstraction.Database {
	db, err := sqlx.Open("postgres", "postgresql://postgres@localhost/mouthful_test?connect_timeout=10&sslmode=disable")
	if err != nil {
		panic(err)
	}
	db.Mapper = reflectx.NewMapperTagFunc("db",
		nil,
		func(s string) string {
			return strings.ToLower(s)
		},
	)
	DB := sqlxDriver.Database{
		DB:      db,
		Queries: PostgresQueries,
		Dialect: "postgres",
		IsTest:  true,
	}
	db.DB.SetMaxOpenConns(1)
	err = DB.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return &DB
}
