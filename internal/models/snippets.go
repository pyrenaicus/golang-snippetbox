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
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	// connect to pool and execute stmt, this returns a sql.Rows result set
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// initialize an empty slice to hold the Snippet structs
	snippets := []*Snippet{}

	// iterate through the rows in the result set with Next(), this prepares each row
	// to be acted on by rows.Scan(). If iteration completes then resultset automatically
	// closes itself and frees-up the underlying db connection
	for rows.Next() {
		// create a pointer to a zeroed struct
		s := &Snippet{}
		// use rows.Scan() to copy the values from each field in the row to the new
		// Snippet object created
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// append it to the slice of snippets
		snippets = append(snippets, s)
	}

	// when rows.Next() finishes, we call rows.Err() to retrieve any error
	// encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// if everything went ok, return Snippets slice
	return snippets, nil
}
