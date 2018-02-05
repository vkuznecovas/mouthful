package main

import (
	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/db/sqlite"
)

func main() {
	db := sqlite.CreateTestDatabase()
	api.GetServer(&db).Run()
}
