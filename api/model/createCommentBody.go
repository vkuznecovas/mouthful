package model

// CreateCommentBody is a struct that represents a create comment request
type CreateCommentBody struct {
	Path   string  `json:"path"`
	Body   string  `json:"body"`
	Author string  `json:"author"`
	Email  *string `json:"email,omitempty"`
}
