package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vkuznecovas/mouthful/db/tool"

	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic(errors.New("Please provide a config filename"))
	}
	// read config.json
	if _, err := os.Stat(argsWithoutProg[0]); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Couldn't find config file:", argsWithoutProg[0])
		os.Exit(1)
	}

	contents, err := ioutil.ReadFile(argsWithoutProg[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't read the config file")
		panic(err)
	}

	// unmarshal config
	config, err := config.ParseConfig(contents)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't parse the config file")
		panic(err)
	}

	// set up db according to config
	database, err := db.GetDBInstance(config.Database)
	if err != nil {
		panic(err)
	}

	err = tool.ExportData("./mouthful.dmp", database.GetAllThreads, database.GetAllComments)
	if err != nil {
		panic(err)
	}

	log.Println("Done!")
}
