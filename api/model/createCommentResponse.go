package model

// CreateCommentResponse is a struct that represents a create comment response
type CreateCommentResponse struct {
	Id      string  `json:"id"`
	Path    string  `json:"path"`
	Body    string  `json:"body"`
	Author  string  `json:"author"`
	Email   *string `json:"email,omitempty"`
	ReplyTo *string `json:"replyTo,omitempty"`
}
