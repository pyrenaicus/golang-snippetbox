package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	// create a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// insert into users table
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// if error, check wether the error has type *mysql.MySQLError, if it does,
		// assign it to mySQLError variable. Then check if it relates to our
		// users_uc_email key by checking error code equals 1062 and contents
		// of error message string. If it does, return ErrDuplicateEmail error
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate verifies wether a user exists with provided email & password,
// return user ID if they do
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists checks if a user exists with given ID
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
