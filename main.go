package main

import (
	"io/ioutil"

	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db"
)

func main() {
	contents, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	config, err := config.ParseConfig(contents)
	if err != nil {
		panic(err)
	}
	database, err := db.GetDBInstance(config.Database)
	if err != nil {
		panic(err)
	}
	api.GetServer(&database).Run()
}
