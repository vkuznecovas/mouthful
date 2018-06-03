package model

// DeleteCommentBody is a struct that represents a delete comment request
type DeleteCommentBody struct {
	CommentId string `json:"commentId"`
	Hard      bool   `json:"hard"`
}
