package models

import (
	"database/sql"
	"time"
)

// define a Snippet type to hold the data for an individual snippet. notice how the fields of the struct corresponds to the fields in our mysql snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// this will insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// this will return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// this will return the 10 most recently created snippet
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
