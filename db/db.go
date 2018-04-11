package db

import (
	"fmt"
	"strings"

	"github.com/vkuznecovas/mouthful/db/dynamodb"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver"

	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

// GetDBInstance looks at the database config object and returns a corresponding database instance
func GetDBInstance(databaseConfig model.Database) (db abstraction.Database, err error) {
	switch strings.ToLower(databaseConfig.Dialect) {
	case "sqlite3":
		db, err = sqlxDriver.CreateDatabase(databaseConfig)
		return db, err
	case "dynamodb":
		db, err = dynamodb.CreateDatabase(databaseConfig)
		return db, err
	default:
		err = fmt.Errorf("unsupported dialect %v", databaseConfig.Dialect)
		return db, err
	}
}
