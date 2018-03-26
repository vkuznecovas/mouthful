package model

// Thread represents an isso commenting thread
type Thread struct {
	Id    int     `db:"id"`
	Uri   *string `db:"uri"`
	Title *string `db:"title"`
}
