package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/vkuznecovas/mouthful/global"

	uuid "github.com/satori/go.uuid"
	"github.com/vkuznecovas/mouthful/api"
	configModel "github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db"
	dbModel "github.com/vkuznecovas/mouthful/db/model"

	"github.com/vkuznecovas/mouthful/db/sqlxDriver"

	"github.com/vkuznecovas/mouthful/cmd/migration/disqus/model"
)

type cpm struct {
	Uid    uuid.UUID
	Parent *string
}

var commentParentMap map[string]cpm
var toDelete []uuid.UUID

func getThread(threads *[]*model.Cthread, id string) *model.Cthread {
	for _, v := range *threads {
		if (v.AttrDsqSpaceid) == id {
			return v
		}
	}
	return nil
}

func insertComment(comment *model.Cpost, comments *[]*model.Cpost, threads *[]*model.Cthread, database sqlxDriver.Database) error {

	// Insert parent if parent exists, and all its parents if needed
	if _, ok := commentParentMap[comment.AttrDsqSpaceid]; ok {
		if commentParentMap[comment.AttrDsqSpaceid].Parent != nil {
			var c *model.Cpost
			for _, v := range *comments {
				if v.AttrDsqSpaceid == *commentParentMap[comment.AttrDsqSpaceid].Parent {
					c = v
					break
				}
			}
			err := insertComment(c, comments, threads, database)
			if err != nil {
				panic(err)
			}

		}
		// insert the comment itself
		t := comment.Cthread
		if len(t) > 0 {
			thread := getThread(threads, t[0].AttrDsqSpaceid)
			u, err := url.Parse(thread.Clink.SValue)
			if err != nil {
				panic(err)
			}
			path := api.NormalizePath(u.Path)
			author := ""
			if comment.Cauthor.Cname.SValue != "" {
				author = comment.Cauthor.Cname.SValue
			} else if comment.Cauthor.Cusername.SValue != "" {
				author = comment.Cauthor.Cusername.SValue
			} else if comment.Cauthor.Cemail.SValue != "" {
				author = comment.Cauthor.Cemail.SValue
			}
			var replyTo *uuid.UUID
			replyTo = nil
			if commentParentMap[comment.AttrDsqSpaceid].Parent != nil {
				cp := commentParentMap[*commentParentMap[comment.AttrDsqSpaceid].Parent]
				replyTo = &cp.Uid
			}
			t, err := database.GetThread(path)
			if err != nil {
				if err == global.ErrThreadNotFound {
					tuid, err := database.CreateThread(path)
					if err != nil {
						panic(err)
					}
					t = dbModel.Thread{
						Id: *tuid,
					}
				}
			}
			uid := uuid.NewV4()
			createdAt, err := time.Parse(time.RFC3339, comment.CcreatedAt.SValue)
			if err != nil {
				panic(err)
			}
			_, err = database.DB.Exec(database.DB.Rebind("INSERT INTO comment(id, threadId, body, author, confirmed, createdAt, replyTo, deletedAt) VALUES(?,?,?,?,?,?,?,?)"), uid, t.Id, comment.Cmessage.SValue, author, true, createdAt, replyTo, nil)
			if err != nil {
				panic(err)
			}
			if comment.CisSpam.SValue == "true" || comment.CisDeleted.SValue == "true" {
				toDelete = append(toDelete, uid)
			}
			editedMapValue := commentParentMap[comment.AttrDsqSpaceid]
			editedMapValue.Uid = uid
			commentParentMap[comment.AttrDsqSpaceid] = editedMapValue
		} else {
			err := fmt.Errorf("Comment %v has no threads", comment.AttrDsqSpaceid)
			return err
		}
	} else {
		err := fmt.Errorf("Comment %v not found", comment.AttrDsqSpaceid)
		return err
	}

	return nil
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic(errors.New("Please provide a source database filename"))
	}
	// read disqus.xml
	contents, err := ioutil.ReadFile(argsWithoutProg[0])
	if err != nil {
		panic(err)
	}

	commentParentMap = make(map[string]cpm, 0)
	dbFile := "./mouthful.db"
	mouthfulDB, err := db.GetDBInstance(configModel.Database{
		Database: &dbFile,
		Dialect:  "sqlite3",
	})
	st := mouthfulDB.GetUnderlyingStruct()
	driverCasted := st.(*sqlxDriver.Database)

	if err != nil {
		panic(err)
	}

	var dis model.Cdisqus
	err = xml.Unmarshal([]byte(contents), &dis)
	if err != nil {
		panic(err)
	}

	// first we form a map for comments, we'll need this to get their parent
	for _, v := range dis.Cpost {
		m := cpm{
			Uid: uuid.NewV4(),
		}
		if v.Cparent != nil {
			m.Parent = &v.Cparent.AttrDsqSpaceid
		} else {
			m.Parent = nil
		}
		commentParentMap[v.AttrDsqSpaceid] = m
	}

	for _, v := range dis.Cpost {
		err = insertComment(v, &dis.Cpost, &dis.Cthread, *driverCasted)
		if err != nil {
			panic(err)
		}
	}

	for _, v := range toDelete {
		err = driverCasted.DeleteComment(v)
		if err != nil {
			panic(err)
		}
	}
}
