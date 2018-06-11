package command

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli"
	"github.com/vkuznecovas/mouthful/config"
	"github.com/vkuznecovas/mouthful/db/dynamodb"
	dynamoModel "github.com/vkuznecovas/mouthful/db/dynamodb/model"
	"github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/global"
)

// DynamoCommandRun will migrate all the threads and comments from mouthful sqlite to mouthful dynamodb.
func DynamoCommandRun(sqlitePath, configPath string) error {
	mouthDB, err := sqlx.Connect("sqlite3", sqlitePath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't connect to sqlite instance %v \n Error: %v ", sqlitePath, err.Error()), 1)
	}

	// read config.json
	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't read config %v \n Error: %v ", configPath, err.Error()), 1)
	}

	// unmarshal config
	config, err := config.ParseConfig(contents)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't parse config %v \n Error: %v ", configPath, err.Error()), 1)
	}

	dynamoDb, err := dynamodb.CreateDatabase(config.Database)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't create dynamodb %v ", err.Error()), 1)
	}
	err = dynamoDb.InitializeDatabase()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't initialize dynamodb %v ", err.Error()), 1)
	}

	plainDynamoDriver := dynamoDb.GetUnderlyingStruct().(*dynamodb.Database)
	threads, err := mouthDB.Queryx(mouthDB.Rebind("select * from thread"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't get sqlx threads %v ", err.Error()), 1)
	}
	log.Println("Migration started")
	for threads.Next() {
		var t model.Thread
		err = threads.StructScan(&t)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Couldn't get thread %v \n Error: %v", t, err.Error()), 1)
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
			return cli.NewExitError(fmt.Sprintf("Couldn't insert thread to dynamodb %v \n Error: %v", threadToInsert, err.Error()), 1)
		}
		comments, err := mouthDB.Queryx(mouthDB.Rebind("select * from comment where ThreadId = ? order by createdAt asc"), t.Id)
		for comments.Next() {
			var c model.Comment
			err = comments.StructScan(&c)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Couldn't read sqlite comment %v \n Error: %v", c, err.Error()), 1)
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
				return cli.NewExitError(fmt.Sprintf("Couldn't insert dynamodbComment comment %v \n Error: %v", commentToInsert, err.Error()), 1)
			}
			log.Printf("Comment %v migrated!\n", c.Id)
		}
		log.Printf("Thread %v migrated!\n", t.Path)
	}
	log.Println("Migration done!")
	return nil
}
