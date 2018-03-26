package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vkuznecovas/mouthful/cmd/migration/isso/model"
)

func main() {
	issoDB, err := sqlx.Connect("sqlite3", "/Users/viktoras/GIT/GO/src/github.com/vkuznecovas/mouthful/isso.db")
	if err != nil {
		panic(err)
	}
	mouthDB, err := sqlx.Connect("sqlite3", "./mouthful")
	if err != nil {
		panic(err)
	}
	for _, v := range mouthDB.SqliteQueries {
		mouthDB.MustExec(v)
	}
	threads := make([]model.Thread, 0)
	err = issoDB.Select(&threads, issoDB.Rebind("select * from threads"))
	if err != nil {
		panic(err)
	}
	for _, v := range threads {
		mouthDB.Exec("insert into thread() values()")
	}
}
