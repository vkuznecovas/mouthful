package model

// Thread represents an isso commenting thread
type Thread struct {
	Id    int     `db:"id"`
	Uri   *string `db:"uri"`
	Title *string `db:"title"`
}

// 'CREATE TABLE IF NOT EXISTS comments (',
// 	'    tid REFERENCES threads(id), id INTEGER PRIMARY KEY, parent INTEGER,',
// 	'    created FLOAT NOT NULL, modified FLOAT, mode INTEGER, remote_addr VARCHAR,',
// 	'    text VARCHAR, author VARCHAR, email VARCHAR, website VARCHAR,',
// 	'    likes INTEGER DEFAULT 0, dislikes INTEGER DEFAULT 0, voters BLOB NOT NULL);'])

// Comment represents an isso commenting comment
type Comment struct {
	Id            int      `db:"id"`
	Tid           int      `db:"tid"`
	Parent        *int     `db:"parent"`
	Created       float64  `db:"created"`
	Modified      *float64 `db:"modified"`
	Mode          *int     `db:"mode"`
	RemoteAddress *string  `db:"remote_addr"`
	Text          *string  `db:"text"`
	Author        *string  `db:"author"`
	Email         *string  `db:"email"`
	Website       *string  `db:"website"`
	Likes         int      `db:"likes"`
	Dislikes      int      `db:"dislikes"`
	Voters        []byte   `db:"voters"`
}
