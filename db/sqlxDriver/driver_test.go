package sqlxDriver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/db"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/mysql"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/postgres"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/sqlite"
)

func setupSqliteTestDb() abstraction.Database {
	database := sqlite.CreateTestDatabase()
	return database
}

func TestSqliteDb(t *testing.T) {
	for _, f := range db.TestFunctions {
		f.(func(*testing.T, abstraction.Database))(t, setupSqliteTestDb())
	}
}

func TestPostgresDB(t *testing.T) {
	testDb := postgres.CreateTestDatabase()
	driver := testDb.GetUnderlyingStruct()
	driverCasted := driver.(*sqlxDriver.Database)
	// clean out before start
	driverCasted.WipeOutData()
	for _, f := range db.TestFunctions {
		f.(func(*testing.T, abstraction.Database))(t, testDb)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}

func TestMysqlDB(t *testing.T) {
	testDb := mysql.CreateTestDatabase()
	driver := testDb.GetUnderlyingStruct()
	driverCasted := driver.(*sqlxDriver.Database)
	// clean out before start
	driverCasted.WipeOutData()
	for _, f := range db.TestFunctions {
		f.(func(*testing.T, abstraction.Database))(t, testDb)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}
