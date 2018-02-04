package model

// UpdateCommentBody is a struct that represents an update comment request
type UpdateCommentBody struct {
	CommentId int     `json:"commentId"`
	Body      *string `json:"body,omitempty"`
	Author    *string `json:"author,omitempty"`
	Confirmed *bool   `json:"confirmed,omitempty"`
}
