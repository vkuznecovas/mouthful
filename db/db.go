package db

import (
	"github.com/vkuznecovas/mouthful/db/sqlite"
)

var DB = sqlite.CreateDatabase()
