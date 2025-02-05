package models

import (
	"database/sql"
	"errors"
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

	// writing the sql statement we want to execute. the reason why ? are used is that they indicate placeholder parameters for the data we want to insert, because the data will be provided by the untrusted user input from a form, its a good practice to use placeholder parameters instead of interpolating data in sql query
	stmt := `INSERT INTO snippets (title, content,created, expires) VALUES (?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// DB.Exec() is used for statements which dont return rows(like INSERT and DELETE)
	// use the Exec() method on the embedded connection pool to execute the statement. the first parameter is the sql sttement, followed by the title, content and expiry value for the placeholder parameter. this methods returns a sql.Result type, which contains some basic information about what happened whent the statement was executed
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// use the LastInsertId() method on the result to get the ID of our newly inserted record in the snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// the ID returned has the type int64, sw we convert it to an int type before returning
	return int(id), nil
}

// this will return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*Snippet, error) {

	// write the sql statement we want to execute
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires>UTC_TIMESTAMP() AND id = ?`

	// user the QueryRow() method on the connection pool to execute our sql statement, passing in the untrusted id variable as the value for the placeholder parameter. this returns a pointer to a sql.Row object which holds the result from the database
	row := m.DB.QueryRow(stmt, id)

	// initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// use row.Scan() to copy the values from each field in sql.Row to the corresponding field in the Snippet struct. notice that the arguments to row.Scan are *pointers* to the place you want to copy the data into, and the number of arguments must be exactly the same as the number of columns returned by your statement
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// if the query return no rows, then row.Scan() will return a sql.ErrNoRows error. we use the errors.Is() function check for that error specifically and return our own ErrNoRecord error instead
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// if everything went well, return the snippet object
	return s, nil
}

// this will return the 10 most recently created snippet
func (m *SnippetModel) Latest() ([]*Snippet, error) {

	//write the sql statement we want to execute
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// use the Query() method on the connection pool to execute the query. this returns a sql.Rows resultset containing the result of our query
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// we defer rows.Close() to ensure the sql.Rows resultset is always properlt closed before Latest() method returns. this defer statement should come *after* you check for an error from the Query() method. otherwise, if Query() returns an error, you'll get a panic trying to close a nil result set
	defer rows.Close()

	// initialize an empty slice of pointers to Snippet structs
	snippets := []*Snippet{}

	// use rows.Next to iterate through the rows in the resultset. this prepares the first (and then each subsequent) row to be acted on by the rows.Scan() method. if iteration over all the rows completes then the resultset automatically closes itself and frees up the underlying database connection.
	for rows.Next() {
		// create a pointer to a new zeroed snippet struct
		s := &Snippet{}

		// use rows.Scan() to copy the values from each field in the current row into the corresponding field in the Snippet struct. notice that the arguments to rows.Scan are *pointers* to the place you want to copy the data into, and the number of arguments must be exactly the same as the number of columns returned by the SELECT statement in the sql statement. if there's an error during this scan, we return the error immediately, so we don't continue scanning the
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// append it to the slice of snippets
		snippets = append(snippets, s)
	}
	// when the rows.Next() loop has finished, we call rows.Err() to retrieve any error that was encountered during the iteration, its important to call this - dont assume that a successful iteration was completed over the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// if everything went well, return the slice of pointers to Snippet structs
	return snippets, nil
}
