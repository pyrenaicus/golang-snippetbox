package models

import (
	"database/sql"
	"errors"
	"time"
)

// define a snippet type to hold the data for an individual snippet. The fields of the struct
// correspond to the fields in our MySQL snippets table
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

// insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// execute the statement
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// use the LastInsertId() method on the result to get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// the ID returned has the type int64, so we convert it to an int before returning
	return int(id), nil
}

// return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// use QueryRow() to execute SQL statement, this returns a pointer to a sql.Row object
	row := m.DB.QueryRow(stmt, id)
	// initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}
	// use row.Scan() to copy the values from each field in sql.Row to the corresponding
	// field in Snippet struct.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use errors.Is() fn to check and return
		// our ErrNoRecord error instead.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// if everything went ok, return Snippet object
	return s, nil
}

// return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
