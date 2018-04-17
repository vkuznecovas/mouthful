// Package db is responsible for database access.
package db

import (
	"fmt"
	"strings"

	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/dynamodb"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/mysql"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/postgres"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/sqlite"
)

// GetDBInstance looks at the database config object and returns a corresponding database instance
func GetDBInstance(databaseConfig model.Database) (db abstraction.Database, err error) {
	switch strings.ToLower(databaseConfig.Dialect) {
	case "sqlite3":
		db, err = sqlite.CreateDatabase(databaseConfig)
		return db, err
	case "postgres":
		db, err = postgres.CreateDatabase(databaseConfig)
		return db, err
	case "mysql":
		db, err = mysql.CreateDatabase(databaseConfig)
		return db, err
	case "dynamodb":
		db, err = dynamodb.CreateDatabase(databaseConfig)
		return db, err
	default:
		err = fmt.Errorf("unsupported dialect %v", databaseConfig.Dialect)
		return db, err
	}
}
