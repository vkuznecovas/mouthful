package tool_test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/vkuznecovas/mouthful/global"

	"github.com/stretchr/testify/assert"
	"github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/db/tool"
)

func TestWriteLine(t *testing.T) {
	f, err := os.Create(path)
	assert.Nil(t, err)
	defer f.Close()
	defer func() {
		err := os.Remove(path)
		assert.Nil(t, err)
	}()
	w := bufio.NewWriter(f)
	newline := []byte("\n")
	test := []byte("test")
	tool.WriteLine(w, newline, test)
	w.Flush()
	f.Close()
	dat, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	assert.Equal(t, "test\n", string(dat))
}

func TestExportDataWorksWellWithExpectedInputs(t *testing.T) {
	t1 := global.GetUUID()
	t2 := global.GetUUID()
	c1 := global.GetUUID()
	c2 := global.GetUUID()
	ca := time.Now().UTC()
	da := time.Now().UTC()

	threadFunc := func() ([]model.Thread, error) {
		threads := []model.Thread{
			model.Thread{
				Id:        t1,
				Path:      "/",
				CreatedAt: time.Now().UTC(),
			},
			model.Thread{
				Id:        t2,
				Path:      "/test",
				CreatedAt: time.Now().UTC(),
			},
		}
		return threads, nil
	}
	commentFunc := func() ([]model.Comment, error) {
		comments := []model.Comment{
			model.Comment{
				Id:        c1,
				ThreadId:  t1,
				Body:      "something something",
				Author:    "Author1",
				Confirmed: true,
				CreatedAt: ca,
				DeletedAt: &da,
				ReplyTo:   nil,
			},
			model.Comment{
				Id:        c2,
				ThreadId:  t2,
				Body:      "something something1",
				Author:    "Author1",
				Confirmed: false,
				CreatedAt: ca,
				DeletedAt: nil,
				ReplyTo:   nil,
			},
		}
		return comments, nil
	}

	defer func() {
		err := os.Remove(path)
		assert.Nil(t, err)
	}()
	err := tool.ExportData(path, threadFunc, commentFunc)
	assert.Nil(t, err)
}

func TestExportDataReturnsErrorOnCommentFuncError(t *testing.T) {
	threadFunc := func() ([]model.Thread, error) {
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
		return threads, nil
	}
	commentFunc := func() ([]model.Comment, error) {
		return nil, fmt.Errorf("test")
	}

	err := tool.ExportData(path, threadFunc, commentFunc)
	assert.NotNil(t, err)
	assert.Equal(t, "test", err.Error())
}

func TestExportDataReturnsErrorOnThreadFuncError(t *testing.T) {
	threadFunc := func() ([]model.Thread, error) {
		return nil, fmt.Errorf("test")

	}
	commentFunc := func() ([]model.Comment, error) {
		comments := make([]model.Comment, 0)
		return comments, nil
	}

	err := tool.ExportData(path, threadFunc, commentFunc)
	assert.NotNil(t, err)
	assert.Equal(t, "test", err.Error())
}
