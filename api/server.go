package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

func GetServer(db *abstraction.Database) *gin.Engine {
	r := gin.Default()
	router := New(*db)
	r.GET("/status", router.Status)
	r.GET("/comments", router.GetComments)
	r.POST("/comments", router.CreateComment)
	return r
}
