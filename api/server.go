package api

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

func GetServer(db *abstraction.Database) *gin.Engine {
	r := gin.Default()
	// TODO: config
	store := sessions.NewCookieStore([]byte("mouthful-secret"))
	// TODO: config this
	store.Options(sessions.Options{
		MaxAge: int(360 * time.Minute), //30min
		Path:   "/",
	})
	router := New(*db)
	r.GET("/status", router.Status)
	r.GET("/comments", router.GetComments)
	r.POST("/comments", router.CreateComment)
	r.PATCH("/comments", sessions.Sessions("mouthful-session", store), router.UpdateComment)
	r.DELETE("/comments", sessions.Sessions("mouthful-session", store), router.DeleteComment)
	r.POST("/admin/login", sessions.Sessions("mouthful-session", store), router.Login)
	r.GET("/threads", sessions.Sessions("mouthful-session", store), router.GetAllThreads)
	r.GET("/comments/all", sessions.Sessions("mouthful-session", store), router.GetAllComments)
	return r
}
