package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/db"

	"github.com/vkuznecovas/mouthful/config/model"
)

func TestGetDBInstanceSqlite3(t *testing.T) {
	memory := ":memory:"
	database := model.Database{
		Dialect:  "sqlite3",
		Database: &memory,
	}
	sqliteInstance, err := db.GetDBInstance(database)
	assert.Nil(t, err)
	assert.NotNil(t, sqliteInstance)
	assert.Equal(t, "sqlite3", sqliteInstance.GetDatabaseDialect())
}

func TestGetDBInstanceNotFound(t *testing.T) {
	d := "something"
	memory := ":memory:"
	database := model.Database{
		Dialect:  d,
		Database: &memory,
	}
	_, err := db.GetDBInstance(database)
	assert.NotNil(t, err)
	assert.Equal(t, "unsupported dialect "+d, err.Error())
}
