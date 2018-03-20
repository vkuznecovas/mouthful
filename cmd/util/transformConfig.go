package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/config/model"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic("Please provide a path to config file")
	}
	data, err := ioutil.ReadFile(argsWithoutProg[0])
	if err != nil {
		panic(err)
	}
	c := model.Config{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}
	transformed := config.TransformConfigToClientConfig(&c)
	res, err := json.Marshal(transformed)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./config.front.json")
	if err != nil {
		panic(err)
	}
	_, err = f.Write(res)
	if err != nil {
		panic(err)
	}
	fmt.Println("generated config.front.json")
}
