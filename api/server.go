package api

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cache "github.com/patrickmn/go-cache"
	"github.com/ulule/limiter"
	mgin "github.com/ulule/limiter/drivers/middleware/gin"
	memoryLimiterStore "github.com/ulule/limiter/drivers/store/memory"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/global"
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
	if config.API.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.ForwardedByClientIP = true

	if config.API.Cors.Enabled {
		fmt.Println("coorsing")
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = *config.API.Cors.AllowedOrigins
		corsConfig.AllowMethods = []string{"PUT", "PATCH", "GET", "DELETE", "HEAD", "OPTIONS", "POST"}
		r.Use(cors.New(corsConfig))
	} else {
		r.Use(cors.Default())
	}

	var cacheInstance *cache.Cache
	if config.API.Cache.Enabled {
		expiry := time.Duration(config.API.Cache.ExpiryInSeconds) * time.Second
		interval := time.Duration(config.API.Cache.IntervalInSeconds) * time.Second
		cacheInstance = cache.New(expiry, interval)
	}

	var limitMiddleware *gin.HandlerFunc
	if config.API.RateLimiting.Enabled {
		limit, err := limiter.NewRateFromFormatted(fmt.Sprintf("%v-H", config.API.RateLimiting.PostsHour))
		if err != nil {
			return nil, err
		}
		newInstance := gin.HandlerFunc(mgin.NewMiddleware(limiter.New(memoryLimiterStore.NewStore(), limit)))
		limitMiddleware = &newInstance
	}

	router := New(db, config, cacheInstance)

	if config.Moderation.Enabled {
		fs := static.LocalFile(global.StaticPath, true)
		r.Use(static.Serve("/", fs))
	} else {
		// We only serve client.js then
		customFs := UnmoderatedFs{
			FileSystem: gin.Dir(global.StaticPath, true),
		}
		r.Use(static.Serve("/", customFs))
	}

	r.GET("/status", router.Status)

	v1 := r.Group("/v1")
	v1.GET("/client/config", router.GetClientConfig)
	v1.GET("/comments", router.GetComments)

	if limitMiddleware != nil {
		v1.POST("/comments", *limitMiddleware, router.CreateComment)
	} else {
		v1.POST("/comments", router.CreateComment)
	}

	if config.Moderation.Enabled {
		err := CheckModerationVariables(config)
		if err != nil {
			return nil, err
		}
		store := sessions.NewCookieStore([]byte(config.Moderation.SessionSecret))
		store.Options(sessions.Options{
			// TODO - figure this out
			MaxAge: 0, //int(time.Second * time.Duration(config.Moderation.SessionDurationSeconds)), //30min
			Path:   "/",
		})
		v1.PATCH("/admin/comments", sessions.Sessions(global.DefaultSessionName, store), router.UpdateComment)
		v1.DELETE("/admin/comments", sessions.Sessions(global.DefaultSessionName, store), router.DeleteComment)

		if limitMiddleware != nil {
			v1.POST("/admin/login", *limitMiddleware, sessions.Sessions(global.DefaultSessionName, store), router.Login)
		} else {
			v1.POST("/admin/login", sessions.Sessions(global.DefaultSessionName, store), router.Login)
		}

		v1.POST("/admin/comments/restore", sessions.Sessions(global.DefaultSessionName, store), router.RestoreDeletedComment)
		v1.GET("/admin/threads", sessions.Sessions(global.DefaultSessionName, store), router.GetAllThreads)
		v1.GET("/admin/comments/all", sessions.Sessions(global.DefaultSessionName, store), router.GetAllComments)
	}

	return r, nil
}
