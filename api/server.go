package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

// CheckModerationVariables checks to see if the required moderation flags have been set in the config or not
func CheckModerationVariables(config *model.Config) error {
	sessionSecret := config.Moderation.SessionSecret
	if sessionSecret == "" {
		return fmt.Errorf("config.Moderation.SessionSecret is not defined in config")
	}
	sessionDuration := config.Moderation.SessionDurationSeconds
	if sessionDuration == 0 {
		config.Moderation.SessionDurationSeconds = 3600
	}
	if config.Moderation.AdminPassword == "" {
		return fmt.Errorf("config.Moderation.AdminPassword is not defined in config")
	}
	return nil
}

// GetServer returns an instance of the mouthful server
func GetServer(db *abstraction.Database, config *model.Config) (*gin.Engine, error) {
	r := gin.Default()
	// same as
	// config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	// router.Use(cors.New(config))
	r.Use(cors.Default())
	router := New(db, config)
	r.Use(static.Serve("/", static.LocalFile("./admin/build", true)))
	r.GET("/status", router.Status)
	v1 := r.Group("/v1")
	v1.GET("/comments", router.GetComments)
	v1.POST("/comments", router.CreateComment)

	if config.Moderation.Enabled {
		err := CheckModerationVariables(config)
		if err != nil {
			return nil, err
		}
		store := sessions.NewCookieStore([]byte(config.Moderation.SessionSecret))
		store.Options(sessions.Options{
			MaxAge: 0, //int(time.Second * time.Duration(config.Moderation.SessionDurationSeconds)), //30min
			Path:   "/",
		})
		v1.PATCH("/admin/comments", sessions.Sessions("mouthful-session", store), router.UpdateComment)
		v1.DELETE("/admin/comments", sessions.Sessions("mouthful-session", store), router.DeleteComment)
		v1.POST("/admin/login", sessions.Sessions("mouthful-session", store), router.Login)
		v1.POST("/admin/comments/restore", sessions.Sessions("mouthful-session", store), router.RestoreDeletedComment)
		v1.GET("/admin/threads", sessions.Sessions("mouthful-session", store), router.GetAllThreads)
		v1.GET("/admin/comments/all", sessions.Sessions("mouthful-session", store), router.GetAllComments)
	}

	return r, nil
}
