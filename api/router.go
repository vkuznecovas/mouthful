package api

import (
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/vkuznecovas/mouthful/api/model"
	configModel "github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/global"
)

// Router handles all the different routes as well as stores our  config and db objects
type Router struct {
	db     *abstraction.Database
	config *configModel.Config
}

// New returns a new instance of router
func New(db *abstraction.Database, config *configModel.Config) *Router {
	r := Router{db: db, config: config}
	return &r
}

// Status responds with 200 when asked
func (r *Router) Status(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

// GetComments returns the comments from thread that is passed as query parameter uri
func (r *Router) GetComments(c *gin.Context) {
	path := c.Query("uri")
	if path == "" {
		c.AbortWithStatusJSON(400, global.ErrThreadNotFound.Error())
		return
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
	c.JSON(200, comments)
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
	if r.config.Honeypot && createCommentBody.Email != nil {
		c.AbortWithStatus(204)
		return
	}
	db := *r.db
	err = db.CreateComment(createCommentBody.Body, createCommentBody.Author, createCommentBody.Path, false, createCommentBody.ReplyTo)
	if err != nil {
		if err == global.ErrWrongReplyTo {
			c.AbortWithStatusJSON(400, global.ErrWrongReplyTo.Error())
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError.Error())
		return
	}
	c.AbortWithStatus(204)
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

	if updateCommentBody.Body == nil && updateCommentBody.Author == nil && updateCommentBody.Confirmed == nil {
		c.AbortWithStatusJSON(400, global.ErrBadRequest)
		return
	}
	db := *r.db
	comment, err := db.GetComment(updateCommentBody.CommentId)
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
	err = db.UpdateComment(updateCommentBody.CommentId, body, author, confirmed)
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
	fmt.Println(c.Cookie("mouthful-session"))
	var deleteCommentBody model.DeleteCommentBody
	err := c.BindJSON(&deleteCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	db := *r.db
	err = db.DeleteComment(deleteCommentBody.CommentId)
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
	fmt.Println(c.Cookie("mouthful-session"))
	var deleteCommentBody model.DeleteCommentBody
	err := c.BindJSON(&deleteCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest.Error())
		return
	}
	db := *r.db
	err = db.RestoreDeletedComment(deleteCommentBody.CommentId)
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
	// return true // TODO remove once tested
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
