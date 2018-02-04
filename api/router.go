package api

import (
	"log"

	"github.com/vkuznecovas/mouthful/api/model"

	"github.com/gin-gonic/gin"
	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/global"
)

type router struct {
	db abstraction.Database
}

// New returns a new instance of router
func New(db abstraction.Database) *router {
	r := router{db: db}
	return &r
}

// Status responds with 200 when asked
func (r *router) Status(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

// GetComments returns the comments from thread that is passed as query parameter uri
func (r *router) GetComments(c *gin.Context) {
	path := c.Query("uri")
	if path == "" {
		c.AbortWithStatusJSON(400, global.ErrThreadNotFound)
		return
	}
	comments, err := r.db.GetCommentsByThread(path)

	if err != nil {
		if err == global.ErrThreadNotFound {
			c.AbortWithStatusJSON(404, global.ErrThreadNotFound)
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError)
		return
	}
	c.JSON(200, comments)
}

// CreateComment creates a comment from CreateCommentBody in JSON form
func (r *router) CreateComment(c *gin.Context) {
	var createCommentBody model.CreateCommentBody
	err := c.BindJSON(&createCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest)
		return
	}
	if createCommentBody.Email != nil {
		c.AbortWithStatus(204)
		return
	}
	err = r.db.CreateComment(createCommentBody.Body, createCommentBody.Author, createCommentBody.Path, false)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError)
		return
	}
	c.AbortWithStatus(204)
}

// UpdateComment updates the provided comment in body
func (r *router) UpdateComment(c *gin.Context) {
	var updateCommentBody model.UpdateCommentBody
	err := c.BindJSON(&updateCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest)
		return
	}

	if updateCommentBody.Body == nil && updateCommentBody.Author == nil && updateCommentBody.Confirmed == nil {
		c.AbortWithStatusJSON(400, global.ErrBadRequest)
		return
	}

	comment, err := r.db.GetComment(updateCommentBody.CommentId)
	if err != nil {
		if err == global.ErrCommentNotFound {
			c.AbortWithStatusJSON(404, global.ErrCommentNotFound)
			return
		}
		c.AbortWithStatusJSON(500, global.ErrInternalServerError)
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

	err = r.db.UpdateComment(updateCommentBody.CommentId, body, author, confirmed)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError)
		return
	}
	c.AbortWithStatus(204)
}

// DeleteComment deletes comment by given id
func (r *router) DeleteComment(c *gin.Context) {
	var deleteCommentBody model.DeleteCommentBody
	err := c.BindJSON(&deleteCommentBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, global.ErrBadRequest)
		return
	}
	err = r.db.DeleteComment(deleteCommentBody.CommentId)
	if err != nil {
		if err == global.ErrCommentNotFound {
			c.AbortWithStatusJSON(404, global.ErrCommentNotFound)
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, global.ErrInternalServerError)
		return
	}
	c.AbortWithStatus(204)
}
