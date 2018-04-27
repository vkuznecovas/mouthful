// This is the entry point for mouthful, the self hosted commenting engine.
//
// Upon providing a config, the main program will parse it and start an API.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vkuznecovas/mouthful/global"

	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db"
)

func main() {

	configFlag := flag.String("config", "./data/config.json", "File to read configuration")
	helpFlag := flag.Bool("h", false, "Show help")
	helpFlagLong := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *helpFlag || *helpFlagLong {
		howto()
	}

	// read config.json
	if _, err := os.Stat(*configFlag); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Couldn't find config file:", *configFlag)
		os.Exit(1)
	}

	contents, err := ioutil.ReadFile(*configFlag)
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

func howto() {
	fmt.Println(`
        Welcome to Mouthful

		Mouthful is a lightweight commenting server written in GO and Preact. It's a self hosted alternative to disqus that's ad free.

		Parameters:

		-config			Location of config.json file (Searches in current directory as default)
		-help			Show this screen
        `)
	os.Exit(0)
}
