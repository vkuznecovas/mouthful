package sqlxDriver_test

import (
	"reflect"
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
	suiteReflected := reflect.ValueOf(db.TestSuite{})
	tReflected := reflect.ValueOf(t)
	for _, f := range db.TestFunctions {
		sqliteReflected := reflect.ValueOf(setupSqliteTestDb())
		in := []reflect.Value{suiteReflected, tReflected, sqliteReflected}
		f.Call(in)
	}
}

func TestPostgresDB(t *testing.T) {
	testDb := postgres.CreateTestDatabase()
	driver := testDb.GetUnderlyingStruct()
	driverCasted := driver.(*sqlxDriver.Database)
	// clean out before start
	driverCasted.WipeOutData()
	suiteReflected := reflect.ValueOf(db.TestSuite{})
	dbReflected := reflect.ValueOf(testDb)
	tReflected := reflect.ValueOf(t)
	for _, f := range db.TestFunctions {
		in := []reflect.Value{suiteReflected, tReflected, dbReflected}
		f.Call(in)
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
	suiteReflected := reflect.ValueOf(db.TestSuite{})
	dbReflected := reflect.ValueOf(testDb)
	tReflected := reflect.ValueOf(t)
	for _, f := range db.TestFunctions {
		in := []reflect.Value{suiteReflected, tReflected, dbReflected}
		f.Call(in)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}
