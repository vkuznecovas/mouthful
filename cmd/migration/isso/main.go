package main

import (
	"errors"
	"log"
	"math"
	"os"
	"time"

	"github.com/satori/go.uuid"

	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/global"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vkuznecovas/mouthful/cmd/migration/isso/model"
	"github.com/vkuznecovas/mouthful/db/sqlxdriver/sqlite"
)

type CommentMap struct {
	Id     int
	Uid    uuid.UUID
	Parent *int
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic(errors.New("Please provide a source database filename"))
	}
	issoDB, err := sqlx.Connect("sqlite3", argsWithoutProg[0])
	if err != nil {
		panic(err)
	}
	mouthDB, err := sqlx.Connect("sqlite3", "./mouthful.db")
	if err != nil {
		panic(err)
	}
	for _, v := range sqlite.SqliteQueries {
		mouthDB.MustExec(v)
	}
	threads, err := issoDB.Queryx(issoDB.Rebind("select * from threads"))
	if err != nil {
		panic(err)
	}
	log.Println("Migration started")
	commentMap := make(map[int]CommentMap)
	for threads.Next() {
		var t model.Thread
		err = threads.StructScan(&t)
		if err != nil {
			panic(err)
		}
		uri := api.NormalizePath(*t.Uri)
		log.Println("Migrating thread " + uri)
		tuid := global.GetUUID()
		_, err = mouthDB.Exec(mouthDB.Rebind("INSERT INTO Thread(Id,Path) VALUES(?, ?)"), tuid, uri)
		if err != nil {
			panic(err)
		}
		comments, err := issoDB.Queryx(issoDB.Rebind("select * from comments where tid = ? order by created asc"), t.Id)
		for comments.Next() {

			var c model.Comment
			err = comments.StructScan(&c)
			if err != nil {
				panic(err)
			}
			log.Printf("Migrating comment %v\n", c.Id)
			commentId := global.GetUUID()
			commentMap[c.Id] = CommentMap{
				Id:     c.Id,
				Parent: c.Parent,
				Uid:    commentId,
			}

			sec, dec := math.Modf(c.Created)
			createdTime := time.Unix(int64(sec), int64(dec*(1e9)))
			var deletedAt *time.Time
			confirmed := true
			if c.Mode != nil {
				if *c.Mode == 4 {
					currentTime := time.Now()
					deletedAt = &currentTime
				} else if *c.Mode == 2 {
					confirmed = false
				}
			}
			var replyTo *uuid.UUID
			if c.Parent != nil {
				parent := c.Parent
				res := parent
				for parent != nil {
					if commentMap[*parent].Parent == nil {
						res = parent
						parent = nil
						break
					} else {
						parent = commentMap[*parent].Parent
					}
				}
				copied := commentMap[*res].Uid
				replyTo = &copied
			}
			body := global.ParseAndSaniziteMarkdown(*c.Text)
			_, err = mouthDB.Exec(mouthDB.Rebind("INSERT INTO comment(id, threadId, body, author, confirmed, createdAt, replyTo, deletedAt) VALUES(?,?,?,?,?,?,?,?)"), commentId, tuid, body, *c.Author, confirmed, createdTime, replyTo, deletedAt)
			if err != nil {
				panic(err)
			}
		}

	}

}
