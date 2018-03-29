package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/global"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vkuznecovas/mouthful/db/dynamodb"
	dynamoModel "github.com/vkuznecovas/mouthful/db/dynamodb/model"

	"github.com/vkuznecovas/mouthful/db/model"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic(errors.New("Please provide a source database filename"))
	}
	mouthDB, err := sqlx.Connect("sqlite3", argsWithoutProg[0])
	if err != nil {
		panic(err)
	}

	// read config.json
	contents, err := ioutil.ReadFile(argsWithoutProg[1])
	if err != nil {
		panic(err)
	}

	// unmarshal config
	config, err := config.ParseConfig(contents)
	if err != nil {
		panic(err)
	}

	dynamoDb, err := dynamodb.CreateDatabase(config.Database)
	if err != nil {
		panic(err)
	}
	err = dynamoDb.InitializeDatabase()
	if err != nil {
		panic(err)
	}

	plainDynamoDriver := dynamoDb.GetUnderlyingStruct().(*dynamodb.Database)
	threads, err := mouthDB.Queryx(mouthDB.Rebind("select * from thread"))
	if err != nil {
		panic(err)
	}
	log.Println("Migration started")
	for threads.Next() {
		var t model.Thread
		err = threads.StructScan(&t)
		if err != nil {
			panic(err)
		}
		log.Println("Migrating thread " + t.Path)
		threadToInsert := dynamoModel.Thread{
			Id:        t.Id,
			Path:      t.Path,
			CreatedAt: t.CreatedAt,
		}
		err = plainDynamoDriver.DB.
			Table(plainDynamoDriver.TablePrefix + global.DefaultDynamoDbThreadTableName).
			Put(threadToInsert).Run()
		if err != nil {
			panic(err)
		}
		comments, err := mouthDB.Queryx(mouthDB.Rebind("select * from comment where ThreadId = ? order by createdAt asc"), t.Id)
		for comments.Next() {
			var c model.Comment
			err = comments.StructScan(&c)
			if err != nil {
				panic(err)
			}
			var deletedAt *int64
			if c.DeletedAt != nil {
				creationTime := c.DeletedAt.UnixNano()
				deletedAt = &creationTime
			}
			var replyTo *string
			if c.ReplyTo != nil {
				uuidString := c.ReplyTo.String()
				replyTo = &uuidString
			}
			commentToInsert := dynamoModel.Comment{
				Id:        c.Id,
				ThreadId:  c.ThreadId,
				Body:      c.Body,
				Author:    c.Author,
				Confirmed: c.Confirmed,
				CreatedAt: c.CreatedAt,
				DeletedAt: deletedAt,
				ReplyTo:   replyTo,
			}
			log.Printf("Migrating comment %v\n", c.Id)
			err = plainDynamoDriver.DB.
				Table(plainDynamoDriver.TablePrefix + global.DefaultDynamoDbCommentTableName).
				Put(commentToInsert).Run()
			if err != nil {
				panic(err)
			}
			log.Printf("Comment %v migrated!\n", c.Id)
		}
		log.Printf("Thread %v migrated!\n", t.Path)
	}
	log.Println("Migration done!")

}
