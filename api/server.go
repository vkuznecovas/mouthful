package api

import (
	"fmt"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
)

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
	router := New(db, config)
	r.GET("/status", router.Status)
	r.GET("/comments", router.GetComments)
	r.POST("/comments", router.CreateComment)

	if config.Moderation.Enabled {
		err := CheckModerationVariables(config)
		if err != nil {
			return nil, err
		}
		store := sessions.NewCookieStore([]byte(config.Moderation.SessionSecret))
		store.Options(sessions.Options{
			MaxAge: int(time.Second * time.Duration(config.Moderation.SessionDurationSeconds)), //30min
			Path:   "/",
		})
		r.PATCH("/comments", sessions.Sessions("mouthful-session", store), router.UpdateComment)
		r.DELETE("/comments", sessions.Sessions("mouthful-session", store), router.DeleteComment)
		r.POST("/admin/login", sessions.Sessions("mouthful-session", store), router.Login)
		r.GET("/threads", sessions.Sessions("mouthful-session", store), router.GetAllThreads)
		r.GET("/comments/all", sessions.Sessions("mouthful-session", store), router.GetAllComments)
	}

	return r, nil
}
