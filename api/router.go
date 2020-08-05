// Package api contains all the required methods and structs for serving mouthful requests.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth/gothic"
	"github.com/patrickmn/go-cache"

	"github.com/vkuznecovas/mouthful/api/model"
	cfg "github.com/vkuznecovas/mouthful/config"
	configModel "github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	dbModel "github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/global"
	"github.com/vkuznecovas/mouthful/oauth/provider"
)

// Router handles all the different routes as well as stores our  config and db objects
type Router struct {
	db           *abstraction.Database
	config       *configModel.Config
	cache        *cache.Cache
	clientConfig *configModel.ClientConfig
	adminConfig  *configModel.AdminConfig
	providers    map[string]*provider.Provider
}

// SetProviders sets the OAUTH providers for the router
func (r *Router) SetProviders(input map[string]*provider.Provider) {
	r.providers = input
}

// OAuth initializes the OAuth flow by redirecting the user to the providers login page
func (r *Router) OAuth(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// OAuthCallback handles the oauth callback which finishes the auth procedure. It checks for the admin flag for the user, and if found it will set the user as admin for the rest of the session
func (r *Router) OAuthCallback(c *gin.Context) {
	q := c.Request.URL.Query()
	provider := c.Param("provider")
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, v := range r.providers[provider].AdminUserIds {
		if user.UserID == v {
			session := sessions.Default(c)
			session.Set("isAdmin", true)
			session.Save()
		}
	}
	c.Redirect(307, *r.config.Moderation.OAuthCallbackOrigin)
}

// New returns a new instance of router
func New(db *abstraction.Database, config *configModel.Config, cache *cache.Cache) *Router {
	clientConfig := cfg.TransformConfigToClientConfig(config)
	adminConfig := cfg.TransformToAdminConfig(config)
	r := Router{db: db, config: config, cache: cache, clientConfig: clientConfig, adminConfig: adminConfig}
	return &r
}

// Status responds with 200 when asked
func (r *Router) Status(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

// GetClientConfig returns the client config portion
func (r *Router) GetClientConfig(c *gin.Context) {
	c.JSON(200, *r.clientConfig)
}

// GetAdminConfig returns the admin config portion
func (r *Router) GetAdminConfig(c *gin.Context) {
	c.JSON(200, *r.adminConfig)
}

// GetComments returns the comments from thread that is passed as query parameter uri
func (r *Router) GetComments(c *gin.Context) {
	path := c.Query("uri")
	if path == "" {
		c.AbortWithStatusJSON(400, global.ErrThreadNotFound.Error())
		return
	}
	path = NormalizePath(path)
	if r.cache != nil {
		if cacheHit, found := r.cache.Get(path); found {
			jsonString := cacheHit.(*[]byte)
			c.Writer.Header().Set("X-Cache", "HIT")
			c.Data(200, "application/json; charset=utf-8", *jsonString)
			return
		}
	}
	db := *r.db
	comments, err := db.GetCommentsByThread(path)

	if err != nil {
		if err == global.ErrThreadNotFound {
			c.AbortWithStatusJSON(404, global.ErrThreadNotFound.Error())
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}
	if comments != nil {
		js, err := json.Marshal(comments)
		if err != nil {
			c.JSON(500, global.ErrInternalServerError.Error())
			return
		}
		if r.cache != nil {
			r.cache.Set(path, &js, cache.DefaultExpiration)
			c.Writer.Header().Set("X-Cache", "MISS")
		}
		if len(comments) > 0 {
			c.Data(200, "application/json; charset=utf-8", js)
		} else {
			c.JSON(404, global.ErrThreadNotFound.Error())
		}
		return
	}
	c.AbortWithStatusJSON(404, global.ErrThreadNotFound.Error())
}

// GetAllThreads returns an array of threads
func (r *Router) GetAllThreads(c *gin.Context) {
	if !r.isAdmin(c) {
		c.AbortWithStatusJSON(401, global.ErrUnauthorized.Error())
		return
	}
	db := *r.db
	threads, err := db.GetAllThreads()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}
	if threads == nil {
		threads = make([]dbModel.Thread, 0)
	}
	c.JSON(200, threads)
}

// GetAllComments returns an array of comments
func (r *Router) GetAllComments(c *gin.Context) {
	if !r.isAdmin(c) {
		c.AbortWithStatusJSON(401, global.ErrUnauthorized.Error())
		return
	}
	db := *r.db
	comments, err := db.GetAllComments()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}
	if comments == nil {
		comments = make([]dbModel.Comment, 0)
	}
	c.JSON(200, comments)
}

// CreateComment creates a comment from CreateCommentBody in JSON form
func (r *Router) CreateComment(c *gin.Context) {
	var createCommentBody model.CreateCommentBody
	err := c.BindJSON(&createCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	// uuid validation
	var uid *uuid.UUID
	if createCommentBody.ReplyTo != nil {
		uid, err = global.ParseUUIDFromString(*createCommentBody.ReplyTo)
		if err != nil {
			c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
			return
		}
	}

	// length validation
	if r.config.Moderation.MaxCommentLength != nil {
		if len(createCommentBody.Body) > *r.config.Moderation.MaxCommentLength {
			c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
			return
		}
	}

	// author length validation
	if len(createCommentBody.Author) == 0 {
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	maxAuthorLength := global.DefaultAuthorLengthLimit
	if r.config.Moderation.MaxAuthorLength != nil {
		maxAuthorLength = *r.config.Moderation.MaxAuthorLength
	}
	createCommentBody.Author = ShortenAuthor(createCommentBody.Author, maxAuthorLength)

	// body length validation
	createCommentBody.Body = global.ParseAndSaniziteMarkdown(createCommentBody.Body)
	if len(createCommentBody.Body) == 0 {
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}

	createCommentBody.Path = NormalizePath(createCommentBody.Path)
	if r.config.Honeypot && createCommentBody.Email != nil {
		c.AbortWithStatusJSON(200, model.CreateCommentResponse{
			Id:      uuid.Must(uuid.NewV4()).String(),
			Path:    createCommentBody.Path,
			Body:    createCommentBody.Body,
			Author:  createCommentBody.Author,
			Email:   createCommentBody.Email,
			ReplyTo: createCommentBody.ReplyTo,
		})
		return
	}

	db := *r.db
	confirmed := !r.config.Moderation.Enabled
	commentUID, err := db.CreateComment(createCommentBody.Body, createCommentBody.Author, createCommentBody.Path, confirmed, uid)
	if err != nil {
		if err == global.ErrWrongReplyTo {
			c.AbortWithStatusJSON(400, global.ErrWrongReplyTo.Error())
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}

	if r.config.Notification.Webhook.Enabled {
		url := *r.config.Notification.Webhook.URL

		go func() {
			_, err := http.Post(
				url,
				"application/json",
				bytes.NewBufferString(fmt.Sprintf(`{"message":"%s"}`, "Comment received")),
			)

			if err != nil {
				log.Println(err)
			}
		}()
	}

	c.AbortWithStatusJSON(200, model.CreateCommentResponse{
		Id:      commentUID.String(),
		Path:    createCommentBody.Path,
		Body:    createCommentBody.Body,
		Author:  createCommentBody.Author,
		Email:   createCommentBody.Email,
		ReplyTo: createCommentBody.ReplyTo,
	})
}

// UpdateComment updates the provided comment in body
func (r *Router) UpdateComment(c *gin.Context) {
	if !r.isAdmin(c) {
		c.AbortWithStatusJSON(401, global.ErrUnauthorized.Error())
		return
	}
	var updateCommentBody model.UpdateCommentBody
	err := c.BindJSON(&updateCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	commentId, err := global.ParseUUIDFromString(updateCommentBody.CommentId)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	if updateCommentBody.Body == nil && updateCommentBody.Author == nil && updateCommentBody.Confirmed == nil {
		c.AbortWithStatusJSON(400, global.ErrBadRequest)
		return
	}
	db := *r.db
	comment, err := db.GetComment(*commentId)
	if err != nil {
		if err == global.ErrCommentNotFound {
			c.AbortWithStatusJSON(404, global.ErrCommentNotFound.Error())
			return
		}
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}

	body := comment.Body
	author := comment.Author
	confirmed := comment.Confirmed
	if updateCommentBody.Body != nil {
		body = *updateCommentBody.Body
	}
	if updateCommentBody.Author != nil {
		author = *updateCommentBody.Author
	}
	if updateCommentBody.Confirmed != nil {
		confirmed = *updateCommentBody.Confirmed
	}
	err = db.UpdateComment(*commentId, body, author, confirmed)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}
	c.AbortWithStatus(204)
}

// DeleteComment deletes comment by given id
func (r *Router) DeleteComment(c *gin.Context) {
	if !r.isAdmin(c) {
		c.AbortWithStatusJSON(401, global.ErrUnauthorized.Error())
		return
	}
	var deleteCommentBody model.DeleteCommentBody
	err := c.BindJSON(&deleteCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	commentId, err := global.ParseUUIDFromString(deleteCommentBody.CommentId)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	db := *r.db

	if deleteCommentBody.Hard {
		err = db.HardDeleteComment(*commentId)
	} else {
		err = db.DeleteComment(*commentId)
	}

	if err != nil {
		if err == global.ErrCommentNotFound {
			c.AbortWithStatusJSON(404, global.ErrCommentNotFound.Error())
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}
	c.AbortWithStatus(204)
}

// RestoreDeletedComment restores the deleted comment by given id
func (r *Router) RestoreDeletedComment(c *gin.Context) {
	if !r.isAdmin(c) {
		c.AbortWithStatusJSON(401, global.ErrUnauthorized.Error())
		return
	}
	var deleteCommentBody model.DeleteCommentBody
	err := c.BindJSON(&deleteCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	commentId, err := global.ParseUUIDFromString(deleteCommentBody.CommentId)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	db := *r.db
	err = db.RestoreDeletedComment(*commentId)
	if err != nil {
		if err == global.ErrCommentNotFound {
			c.AbortWithStatusJSON(404, global.ErrCommentNotFound.Error())
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}
	c.AbortWithStatus(204)
}

func (r *Router) isAdmin(c *gin.Context) bool {
	session := sessions.Default(c)
	isAdmin := session.Get("isAdmin")
	isAdminParsed, ok := isAdmin.(bool)
	if !ok {
		return false
	}
	return isAdminParsed
}

// Login logs the user in
func (r *Router) Login(c *gin.Context) {
	var loginBody model.LoginBody
	err := c.BindJSON(&loginBody)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}

	if loginBody.Password != r.config.Moderation.AdminPassword {
		c.AbortWithStatusJSON(401, global.ErrBadRequest.Error())
		return
	}

	session := sessions.Default(c)
	session.Set("isAdmin", true)
	session.Save()
	c.AbortWithStatus(204)
}
