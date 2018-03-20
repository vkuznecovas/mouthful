package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/vkuznecovas/mouthful/global"

	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db"
)

func main() {
	// read config.json
	contents, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	// unmarshal config
	config, err := config.ParseConfig(contents)
	if err != nil {
		panic(err)
	}

	// set up db according to config
	database, err := db.GetDBInstance(config.Database)
	if err != nil {
		panic(err)
	}

	// get GIN server
	service, err := api.GetServer(&database, config)
	if err != nil {
		panic(err)
	}

	// set GIN port
	port := fmt.Sprintf(":%v", global.DefaultPort)
	if config.API.Port != nil {
		port = fmt.Sprintf(":%v", *config.API.Port)
	}

	// run the server
	log.Println("Running server on port", port)
	service.Run(port)
}
