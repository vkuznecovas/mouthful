package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db"
	"github.com/vkuznecovas/mouthful/db/model"
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

	comments, err := database.GetAllComments()
	if err != nil {
		panic(err)
	}
	threads, err := database.GetAllThreads()
	if err != nil {
		panic(err)
	}

	dump := model.DataDump{
		ThreadCount:  len(threads),
		CommentCount: len(comments),
	}
	marshaledDump, err := json.Marshal(dump)
	if err != nil {
		panic(err)
	}
	log.Println(string(marshaledDump))
	f, err := os.Create("mouthful.dmp")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	newline := []byte("\n")
	writeLine(w, newline, marshaledDump)
	w.Flush()
	for i, v := range threads {
		marshaledThread, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		writeLine(w, newline, marshaledThread)
		if i%100 == 0 {
			log.Printf("Written %v threads", i)
			w.Flush()
		}
	}
	w.Flush()
	for i, v := range comments {
		marshaledComment, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		writeLine(w, newline, marshaledComment)
		if i%100 == 0 {
			log.Printf("Written %v comments", i)
			w.Flush()
		}
	}
	w.Flush()
	log.Println("Done!")
}
