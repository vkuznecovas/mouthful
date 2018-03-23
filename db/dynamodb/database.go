package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

// TODO: tests

// Database is a database instance for sqlite
type Database struct {
	DB          *dynamo.DB
	TablePrefix string
}

// CreateDatabase creates a database instance from the given config
func CreateDatabase(databaseConfig model.Database) (abstraction.Database, error) {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("eu-west-1"), Endpoint: aws.String("http://localhost:8080")})
	prefix := ""
	if databaseConfig.TablePrefix != nil {
		prefix = *databaseConfig.TablePrefix
	}
	return &Database{
		DB:          db,
		TablePrefix: prefix,
	}, nil
	// // put item
	// w := widget{UserID: 613, Time: time.Now(), Msg: "hello"}
	// err := table.Put(w).Run()

	// // get the same item
	// var result widget
	// err = table.Get("UserID", w.UserID).
	// 	Range("Time", dynamo.Equal, w.Time).
	// 	Filter("'Count' = ? AND $ = ?", w.Count, "Message", w.Msg). // placeholders in expressions
	// 	One(&result)

	// // get all items
	// var results []widget
	// err = table.Scan().All(&results)
}
