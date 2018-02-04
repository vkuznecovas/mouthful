package model

// Thread represents a commenting thread
type Thread struct {
	Id   int    `db:"Id"`
	Path string `db:"Path"`
}
