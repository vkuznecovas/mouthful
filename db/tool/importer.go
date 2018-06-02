// Package tool contains various utilities used with the mouthful database
package tool

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/vkuznecovas/mouthful/db/model"
)

// ImportData is responsible for data imports to mouthful.
func ImportData(pathToDump string, importThread func(model.Thread) error, importComment func(model.Comment) error) error {
	file, err := os.Open(pathToDump)
	if err != nil {
		return fmt.Errorf("Could not open data dump at %v. \n %v", pathToDump, err.Error())
	}

	reader := bufio.NewReader(file)

	currentLine := 1
	// first line of the dump should contain the data dump object in json format
	dataDumpJson, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("Failed to read from data dump at line %v. \n %v", currentLine, err.Error())
	}
	currentLine++

	var dataDumpStruct model.DataDump
	// unmarshal the data dump
	err = json.Unmarshal(dataDumpJson, &dataDumpStruct)
	if err != nil {
		return fmt.Errorf("Corrupted data dump. Could not deserialize the dump header at line 1. \n %v", err.Error())
	}

	log.Println("Importing threads")
	for i := 0; i < dataDumpStruct.ThreadCount; i++ {
		threadJson, _, err := reader.ReadLine()
		var thread model.Thread
		err = json.Unmarshal(threadJson, &thread)
		if err != nil {
			return fmt.Errorf("Corrupted data dump. Could not deserialize thread JSON at line %v. \n %v", currentLine, err.Error())
		}
		err = importThread(thread)
		if err != nil {
			return fmt.Errorf("Failed to insert the thread at line %v. \n %v", currentLine, err.Error())
		}
		currentLine++
		log.Printf("Thread %v done!\n", i)
	}
	log.Println("Threads imported")
	log.Println("Importing comments")
	for i := 0; i < dataDumpStruct.CommentCount; i++ {
		commentJson, _, err := reader.ReadLine()
		var comment model.Comment
		err = json.Unmarshal(commentJson, &comment)
		if err != nil {
			return fmt.Errorf("Corrupted data dump. Could not deserialize comment JSON at line %v. \n %v", currentLine, err.Error())
		}
		err = importComment(comment)
		if err != nil {
			return fmt.Errorf("Failed to insert the comment at line %v. \n %v", currentLine, err.Error())
		}
		currentLine++
		log.Printf("Comment %v done!\n", i)
	}
	log.Println("Comments imported!")
	return nil
}
