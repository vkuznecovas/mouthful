package command

import (
	"fmt"
	"log"
	"math"
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli"
	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/cmd/spoon/command/model"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/sqlite"
	"github.com/vkuznecovas/mouthful/global"
)

type commentParentMapIsso struct {
	Id     int
	Uid    uuid.UUID
	Parent *int
}

// IssoCommandRun takes the issoDbPath connects to it and mmigrates all the comments to a new mouthful sqlite instance
func IssoCommandRun(issoDbPath string) error {
	issoDB, err := sqlx.Connect("sqlite3", issoDbPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't connect to isso db %v \n Error: %v", issoDbPath, err.Error()), 1)
	}
	mouthDB, err := sqlx.Connect("sqlite3", "./mouthful.db")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't connect to mouthful db %v \n Error: %v", "./mouthful.db", err.Error()), 1)
	}
	for _, v := range sqlite.SqliteQueries {
		mouthDB.MustExec(v)
	}
	threads, err := issoDB.Queryx(issoDB.Rebind("select * from threads"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't select threads from isso DB %v \n Error: %v", issoDbPath, err.Error()), 1)
	}
	log.Println("Migration started")
	commentMap := make(map[int]commentParentMapIsso)
	for threads.Next() {
		var t model.Thread
		err = threads.StructScan(&t)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Couldn't select thread from isso DB %v", err.Error()), 1)
		}
		uri := api.NormalizePath(*t.Uri)
		log.Println("Migrating thread " + uri)
		tuid := global.GetUUID()
		_, err = mouthDB.Exec(mouthDB.Rebind("INSERT INTO Thread(Id,Path) VALUES(?, ?)"), tuid, uri)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Couldn't insert thread into mouthful DB %v, Error: %v", t, err.Error()), 1)
		}
		comments, err := issoDB.Queryx(issoDB.Rebind("select * from comments where tid = ? order by created asc"), t.Id)
		for comments.Next() {

			var c model.Comment
			err = comments.StructScan(&c)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Couldn't select comment from isso DB %v", err.Error()), 1)
			}
			log.Printf("Migrating comment %v\n", c.Id)
			commentId := global.GetUUID()
			commentMap[c.Id] = commentParentMapIsso{
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
				return cli.NewExitError(fmt.Sprintf("Couldn't insert comment into mouthful DB %v, Error: %v", c, err.Error()), 1)
			}
		}

	}
	return nil
}
