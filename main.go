// This is the entry point for mouthful, the self hosted commenting engine.
//
// Upon providing a config, the main program will parse it and start an API.
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
	contents, err := ioutil.ReadFile("./data/config.json")
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
	port := global.DefaultPort
	if config.API.Port != nil {
		port = *config.API.Port
	}

	// add GIN bind address, serving on all by default
	bindAddress := global.DefaultBindAddress
	if config.API.BindAddress != nil {
		bindAddress = *config.API.BindAddress
	}

	// run the server
	fullAddress := fmt.Sprintf("%v:%v", bindAddress, port)
	log.Println("Running server on ", fullAddress)
	service.Run(fullAddress)
}
