package job

import (
	"fmt"
	"log"
	"time"

	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/global"
)

// StartCleanupJobs starts the required scheduled tasks to delete stale data
func StartCleanupJobs(db abstraction.Database, config *model.PeriodicCleanUp) error {
	if config == nil || !config.Enabled {
		log.Println("Cleanup jobs not enabled, skipping...")
		return nil
	}
	if config.RemoveDeleted {
		period := global.DefaultCleanupPeriod
		if config.RemoveDeletedPeriodSeconds != 0 {
			period = config.RemoveDeletedPeriodSeconds
		}
		if config.DeletedTimeoutSeconds == 0 {
			return fmt.Errorf("DeletedTimeoutSeconds not specified but the deletion job is enabled, please specify a value in config")
		}
		StartCleanupJob(db, config.DeletedTimeoutSeconds, period, global.Deleted)
	}
	if config.RemoveUnconfirmed {
		period := global.DefaultCleanupPeriod
		if config.RemoveUnconfirmedPeriodSeconds != 0 {
			period = config.RemoveUnconfirmedPeriodSeconds
		}
		if config.UnconfirmedTimeoutSeconds == 0 {
			return fmt.Errorf("UnconfirmedTimeoutSeconds not specified but the deletion job is enabled, please specify a value in config")
		}
		StartCleanupJob(db, config.UnconfirmedTimeoutSeconds, period, global.Unconfirmed)
	}
	return nil
}

// StartCleanupJob starts the cleanup job of a given type
func StartCleanupJob(db abstraction.Database, olderThan int64, every int64, t global.CleanupType) {
	duration := time.Duration(every) * time.Second
	ticker := time.NewTicker(duration)
	log.Printf("Cleanup job %v is up and running.\n", t)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Removing old %v comments...\n", t)
				err := db.CleanUpStaleData(t, olderThan)
				if err != nil {
					log.Printf("An error occurred while deleting old %v comments\n", t)
					log.Println(err.Error())
					log.Println("Mouthful will continue running.")
				}
				log.Printf("Old %v comments removed!\n", t)
			}
		}
	}()
}
