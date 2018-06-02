package tool

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/vkuznecovas/mouthful/db/model"
)

// WriteLine writes the content as a new line
func WriteLine(writer io.Writer, lineSymbol []byte, content []byte) error {
	_, err := writer.Write(content)
	if err != nil {
		return err
	}
	_, err = writer.Write(lineSymbol)
	if err != nil {
		return err
	}
	return nil
}

// ExportData is responsible for exporting the data from the database
func ExportData(path string, threadGetter func() ([]model.Thread, error), commentGetter func() ([]model.Comment, error)) error {
	comments, err := commentGetter()
	if err != nil {
		return err
	}
	threads, err := threadGetter()
	if err != nil {
		return err
	}

	dump := model.DataDump{
		ThreadCount:  len(threads),
		CommentCount: len(comments),
	}
	marshaledDump, err := json.Marshal(dump)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	newline := []byte("\n")
	WriteLine(w, newline, marshaledDump)
	w.Flush()
	for i, v := range threads {
		marshaledThread, err := json.Marshal(v)
		if err != nil {
			return err
		}
		WriteLine(w, newline, marshaledThread)
		log.Printf("Written %v threads", i)
		if i%100 == 0 {
			w.Flush()
		}
	}
	w.Flush()
	for i, v := range comments {
		marshaledComment, err := json.Marshal(v)
		if err != nil {
			return err
		}
		WriteLine(w, newline, marshaledComment)
		log.Printf("Written %v comments", i)
		if i%100 == 0 {
			w.Flush()
		}
	}
	w.Flush()
	return nil
}
