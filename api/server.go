package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	cache "github.com/patrickmn/go-cache"
	"github.com/ulule/limiter"
	mgin "github.com/ulule/limiter/drivers/middleware/gin"
	memoryLimiterStore "github.com/ulule/limiter/drivers/store/memory"
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/global"
	"github.com/vkuznecovas/mouthful/oauth"
	"github.com/vkuznecovas/mouthful/oauth/provider"
)

// CheckModerationVariables checks to see if the required moderation flags have been set in the config or not
func CheckModerationVariables(config *model.Config) error {
	if config.Moderation.AdminPassword == "" && !config.Moderation.DisablePasswordLogin {
		return fmt.Errorf("config.Moderation.AdminPassword is not defined in config")
	}

	// force default password change if password login is not disabled
	if config.Moderation.AdminPassword == "somepassword" && !config.Moderation.DisablePasswordLogin {
		return fmt.Errorf("Please change the config.Moderation.AdminPassword value in config. Do not leave the default there")
	}

	// determine if oauth is in use
	hasEnabledAuthProviders := false
	if config.Moderation.OAauthProviders != nil && len(*config.Moderation.OAauthProviders) != 0 {
		for _, v := range *config.Moderation.OAauthProviders {
			if v.Enabled == true {
				hasEnabledAuthProviders = true
				break
			}
		}
	}

	// check if password is disabled but no OAUTH providers are enabled
	if config.Moderation.DisablePasswordLogin == true {
		err := fmt.Errorf("You have moderation enabled with no enabled OAUTH providers or admin password functionality. You will not be able to login on the admin panel. Please check your configuration")
		if !hasEnabledAuthProviders {
			return err
		}
	}

	// if we have providers, we do need the origin specified as well
	if hasEnabledAuthProviders {
		if config.Moderation.OAuthCallbackOrigin == nil || *config.Moderation.OAuthCallbackOrigin == "" {
			return fmt.Errorf("Please provide a OAuthCallbackOrigin in the config moderation section. For more info, refer to the documentation on github")
		}
		if !strings.HasSuffix(*config.Moderation.OAuthCallbackOrigin, "/") {
			properOrigin := *config.Moderation.OAuthCallbackOrigin + "/"
			config.Moderation.OAuthCallbackOrigin = &properOrigin
		}
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

	r := gin.New()
	r.Use(gin.Recovery())
	if config.API.Logging {
		r.Use(gin.Logger())
	}
	r.ForwardedByClientIP = true
	if config.API.Cors.Enabled {
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
		store := cookie.NewStore([]byte(config.Moderation.AdminPassword))
		store.Options(sessions.Options{
			MaxAge: int(time.Second * time.Duration(config.Moderation.SessionDurationSeconds)), //30min
			Path:   "/",
		})
		r.Use(sessions.Sessions("mouthful", store))
		v1.GET("/admin/config", router.GetAdminConfig)
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

		if config.Moderation.OAauthProviders != nil {
			gothic.Store = store
			callbackUrl := *config.Moderation.OAuthCallbackOrigin + "v1/oauth/callbacks/"
			providers, err := oauth.GetProviders(config.Moderation.OAauthProviders, callbackUrl)
			if err != nil {
				return nil, err
			}
			providerMap := make(map[string]*provider.Provider)
			for i := range providers {
				providerMap[providers[i].Name] = &providers[i]
				goth.UseProviders(*providers[i].Implementation)
			}

			router.SetProviders(providerMap)
			v1.GET("/oauth/callbacks/:provider", sessions.Sessions(global.DefaultSessionName, store), router.OAuthCallback)
			v1.GET("/oauth/auth/:provider", sessions.Sessions(global.DefaultSessionName, store), router.OAuth)
		}
	}

	return r, nil
}
