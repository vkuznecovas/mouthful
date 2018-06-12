package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db"
)

// ImportCommandRun imports the provided dump to the database pointed at by config.json
func ImportCommandRun(configPath, dumpPath string) error {
	// read config.json
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("Couldn't find config file %v", configPath), 1)
	}

	// read dump
	if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("Couldn't find dump file %v", dumpPath), 1)
	}

	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't read config file %v", configPath), 1)
	}

	// unmarshal config
	config, err := config.ParseConfig(contents)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't parse the config file %v", err.Error()), 1)
	}

	// set up db according to config
	database, err := db.GetDBInstance(config.Database)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't connect to the database %v", err.Error()), 1)
	}

	err = database.ImportData(dumpPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't import data to the database %v", err.Error()), 1)
	}

	log.Println("Done!")
	return nil
}
