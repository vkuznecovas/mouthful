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

	"github.com/fatih/color"
	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db"
	"github.com/vkuznecovas/mouthful/job"
)

func main() {

	// Print a warning if user is running as root.
	if os.Geteuid() == 0 {
		color.Set(color.FgYellow)
		log.Println("WARNING: Mouthful is running as root. For security reasons please consider creating a non root user for mouthful. Mouthful does not need root permissions to run.")
		color.Unset()
	}

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

	// check if we're gonna need to override the path in static admin html
	if config.Moderation.Path != nil {
		err := global.RewriteAdminPanelScripts(*config.Moderation.Path)
		if err != nil {
			panic(err)
		}
	}

	// set up db according to config
	database, err := db.GetDBInstance(config.Database)
	if err != nil {
		panic(err)
	}

	// startup cleanup, if enabled
	err = job.StartCleanupJobs(database, config.Moderation.PeriodicCleanUp)
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
	color.Set(color.FgGreen)
	log.Println("Running mouthful server on ", fullAddress)
	color.Unset()
	service.Run(fullAddress)
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
