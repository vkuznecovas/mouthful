// Package mysql is responsible for mysql database connections and initialization.
package mysql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	// We absolutely need the mysql driver here, this whole file depends on it
	_ "github.com/go-sql-driver/mysql"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver"
)

// MysqlQueries represents a list of queries for initial table creation in mysql
var MysqlQueries = []string{
	`CREATE TABLE IF NOT EXISTS Thread(
			Id VARCHAR(36) PRIMARY KEY,
			CreatedAt TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) not null,
			Path varchar(255) not null UNIQUE
		)`,
	`CREATE TABLE IF NOT EXISTS Comment(
			Id VARCHAR(36) PRIMARY KEY,
			ThreadId VARCHAR(36) not null,
			Body text not null,
			Author varchar(255) not null,
			Confirmed bool not null default false,
			CreatedAt TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) not null,
			ReplyTo VARCHAR(36) default null,
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
	if config.Port == nil {
		err += "Please specify the database port in config(Database.Port)"
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
	connectionString := fmt.Sprintf("%v:%v@(%v:%v)/%v?parseTime=true", *databaseConfig.Username, *databaseConfig.Password, *databaseConfig.Host, *databaseConfig.Port, *databaseConfig.Database)
	d, err := sqlx.Connect("mysql", connectionString)
	d.MapperFunc(func(s string) string { return strings.Title(s) })
	if err != nil {
		return nil, err
	}
	db = d
	DB := sqlxDriver.Database{
		DB:      db,
		Queries: MysqlQueries,
		Dialect: "mysql",
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
	db, err := sqlx.Open("mysql", "root:@(localhost:3306)/mouthful_test?parseTime=true")
	if err != nil {
		panic(err)
	}
	db.MapperFunc(func(s string) string { return strings.Title(s) })
	db.DB.SetMaxOpenConns(1)
	DB := sqlxDriver.Database{
		DB:      db,
		Queries: MysqlQueries,
		Dialect: "mysql",
		IsTest:  true,
	}
	err = DB.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return &DB
}
