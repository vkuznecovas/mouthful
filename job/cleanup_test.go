package job_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/sqlite"
	"github.com/vkuznecovas/mouthful/global"
	"github.com/vkuznecovas/mouthful/job"
)

// TODO: this is rather poorly tested and even written(aka untestable). Might need a rethink.

func TestStartCleanupJob(t *testing.T) {
	db := sqlite.CreateTestDatabase()
	job.StartCleanupJob(db, 1, 1, global.Deleted)
}

var cleanupConfig = model.PeriodicCleanUp{
	Enabled:                        true,
	RemoveDeleted:                  true,
	RemoveUnconfirmed:              true,
	DeletedTimeoutSeconds:          1,
	UnconfirmedTimeoutSeconds:      1,
	RemoveDeletedPeriodSeconds:     1,
	RemoveUnconfirmedPeriodSeconds: 1,
}

func TestStartCleanupJobs(t *testing.T) {
	db := sqlite.CreateTestDatabase()
	err := job.StartCleanupJobs(db, &cleanupConfig)
	assert.Nil(t, err)
}

func TestStartCleanupJobsReturnsNilErrorOnNilConfig(t *testing.T) {
	db := sqlite.CreateTestDatabase()
	err := job.StartCleanupJobs(db, nil)
	assert.Nil(t, err)
}

func TestStartCleanupJobsReturnsNilErrorOnDisabled(t *testing.T) {
	cleanupConfig.Enabled = false
	defer func() { cleanupConfig.Enabled = true }()
	db := sqlite.CreateTestDatabase()
	err := job.StartCleanupJobs(db, &cleanupConfig)
	assert.Nil(t, err)
}
func TestStartCleanupJobsReturnsNonNilErrorOnZeroTimeoutForDeleted(t *testing.T) {
	cleanupConfig.DeletedTimeoutSeconds = 0
	defer func() { cleanupConfig.DeletedTimeoutSeconds = 1 }()
	db := sqlite.CreateTestDatabase()
	err := job.StartCleanupJobs(db, &cleanupConfig)
	assert.NotNil(t, err)
	assert.Equal(t, "DeletedTimeoutSeconds not specified but the deletion job is enabled, please specify a value in config", err.Error())
}

func TestStartCleanupJobsReturnsNonNilErrorOnZeroTimeoutForUnconfirmed(t *testing.T) {
	cleanupConfig.UnconfirmedTimeoutSeconds = 0
	defer func() { cleanupConfig.UnconfirmedTimeoutSeconds = 1 }()
	db := sqlite.CreateTestDatabase()
	err := job.StartCleanupJobs(db, &cleanupConfig)
	assert.NotNil(t, err)
	assert.Equal(t, "UnconfirmedTimeoutSeconds not specified but the deletion job is enabled, please specify a value in config", err.Error())
}
