// Package command contains all the spoon commands
package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db"
	"github.com/vkuznecovas/mouthful/db/tool"
)

// ExportCommandRun exports the database to mouthful format
func ExportCommandRun(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("Couldn't find config file %v", configPath), 1)
	}
	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't read config file %v", configPath), 1)
	}

	// unmarshal config
	config, err := config.ParseConfig(contents)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't parse the config file %v", configPath), 1)
	}

	// set up db according to config
	database, err := db.GetDBInstance(config.Database)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't connect to the database %v", err.Error()), 1)
	}

	err = tool.ExportData("./mouthful.dmp", database.GetAllThreads, database.GetAllComments)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't export data %v", err.Error()), 1)
	}

	log.Println("Dump done!")
	return nil
}
