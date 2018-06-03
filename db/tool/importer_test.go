package tool_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vkuznecovas/mouthful/db/tool"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/global"

	"github.com/vkuznecovas/mouthful/db/model"
)

const path = "./mouthful_tool_test.dmp"

func GetDump() ([]string, error) {
	result := make([]string, 0)
	threads := []model.Thread{
		model.Thread{
			Id:        global.GetUUID(),
			Path:      "/",
			CreatedAt: time.Now().UTC(),
		},
		model.Thread{
			Id:        global.GetUUID(),
			Path:      "/test",
			CreatedAt: time.Now().UTC(),
		},
	}
	ca := time.Now().UTC()
	da := time.Now().UTC()
	comments := []model.Comment{
		model.Comment{
			Id:        global.GetUUID(),
			ThreadId:  threads[0].Id,
			Body:      "something something",
			Author:    "Author1",
			Confirmed: true,
			CreatedAt: ca,
			DeletedAt: &da,
			ReplyTo:   nil,
		},
		model.Comment{
			Id:        global.GetUUID(),
			ThreadId:  threads[1].Id,
			Body:      "something something1",
			Author:    "Author1",
			Confirmed: false,
			CreatedAt: ca,
			DeletedAt: nil,
			ReplyTo:   nil,
		},
	}
	dataDump := model.DataDump{
		ThreadCount:  len(threads),
		CommentCount: len(comments),
	}
	res, err := json.Marshal(dataDump)
	if err != nil {
		return result, err
	}
	result = append(result, string(res))

	for _, v := range threads {
		res, err = json.Marshal(v)
		if err != nil {
			panic(err)
		}
		result = append(result, string(res))
	}

	for _, v := range comments {
		res, err = json.Marshal(v)
		if err != nil {
			panic(err)
		}
		result = append(result, string(res))
	}

	return result, nil
}

func WriteDumpToFile(lines []string) error {
	joined := strings.Join(lines, "\n")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(joined)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDumpFile() error {
	return os.Remove(path)
}

func threadFunc(model.Thread) error {
	return nil
}

func commentFunc(model.Comment) error {
	return nil
}

func threadFuncFail(model.Thread) error {
	return fmt.Errorf("fail")
}

func commentFuncFail(model.Comment) error {
	return fmt.Errorf("fail")
}

func TestImportDataWorksWithCorrectInputs(t *testing.T) {
	dump, err := GetDump()
	assert.Nil(t, err)
	err = WriteDumpToFile(dump)
	assert.Nil(t, err)
	defer func() {
		err := DeleteDumpFile()
		assert.Nil(t, err)
	}()
	err = tool.ImportData(path, threadFunc, commentFunc)
	assert.Nil(t, err)
}

func TestImportDataBadFilePathReturnsError(t *testing.T) {
	err := tool.ImportData("path", threadFunc, commentFunc)
	assert.NotNil(t, err)
	assert.Equal(t, "Could not open data dump at path. \n open path: no such file or directory", err.Error())
}

func TestImportDataBadHeaderReturnsError(t *testing.T) {
	dump, err := GetDump()
	assert.Nil(t, err)
	dump[0] = "asdasd" + dump[0]
	err = WriteDumpToFile(dump)
	assert.Nil(t, err)
	defer func() {
		err := DeleteDumpFile()
		assert.Nil(t, err)
	}()
	err = tool.ImportData(path, threadFunc, commentFunc)
	assert.NotNil(t, err)
	assert.Equal(t, "Corrupted data dump. Could not deserialize the dump header at line 1. \n invalid character 'a' looking for beginning of value", err.Error())
}

func TestImportDataBadCommentJsonReturnsError(t *testing.T) {
	dump, err := GetDump()
	assert.Nil(t, err)
	dump[4] = "asdasd" + dump[4]
	err = WriteDumpToFile(dump)
	assert.Nil(t, err)
	defer func() {
		err := DeleteDumpFile()
		assert.Nil(t, err)
	}()
	err = tool.ImportData(path, threadFunc, commentFunc)
	assert.NotNil(t, err)
	assert.Equal(t, "Corrupted data dump. Could not deserialize comment JSON at line 5. \n invalid character 'a' looking for beginning of value", err.Error())
}

func TestImportDataBadThreadJsonReturnsError(t *testing.T) {
	dump, err := GetDump()
	assert.Nil(t, err)
	dump[1] = "asdasd" + dump[1]
	err = WriteDumpToFile(dump)
	assert.Nil(t, err)
	defer func() {
		err := DeleteDumpFile()
		assert.Nil(t, err)
	}()
	err = tool.ImportData(path, threadFunc, commentFunc)
	assert.NotNil(t, err)
	assert.Equal(t, "Corrupted data dump. Could not deserialize thread JSON at line 2. \n invalid character 'a' looking for beginning of value", err.Error())
}

func TestImportDataFailedThreadInsertionReturnsError(t *testing.T) {
	dump, err := GetDump()
	assert.Nil(t, err)
	err = WriteDumpToFile(dump)
	assert.Nil(t, err)
	defer func() {
		err := DeleteDumpFile()
		assert.Nil(t, err)
	}()
	err = tool.ImportData(path, threadFuncFail, commentFunc)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to insert the thread at line 2. \n fail", err.Error())
}

func TestImportDataFailedCommentInsertionReturnsError(t *testing.T) {
	dump, err := GetDump()
	assert.Nil(t, err)
	err = WriteDumpToFile(dump)
	assert.Nil(t, err)
	defer func() {
		err := DeleteDumpFile()
		assert.Nil(t, err)
	}()
	err = tool.ImportData(path, threadFunc, commentFuncFail)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to insert the comment at line 4. \n fail", err.Error())
}

func TestImportDataEmptyDumpReturnsError(t *testing.T) {
	empty := []string{""}
	err := WriteDumpToFile(empty)
	assert.Nil(t, err)
	defer func() {
		err := DeleteDumpFile()
		assert.Nil(t, err)
	}()
	err = tool.ImportData(path, threadFunc, commentFuncFail)
	assert.NotNil(t, err)
	assert.Equal(t, "Failed to read from data dump at line 1. \n EOF", err.Error())
}
