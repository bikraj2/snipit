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
	Db *sql.DB
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
  ChangePassword(id int, oldPassword string,newPassword string)(bool, error)
}

func (m *UserModel) Get(id int) (*User, error) {
	var user User
	stmt := `SELECT name,email,created FROM USERS WHERE id  = ?`
	exists, err := m.Exists(id)
	if err != nil {
		return &User{}, err
	} else if !exists {
		return &User{}, ErrNoRecord
	}
	err = m.Db.QueryRow(stmt, id).Scan(&user.Name, &user.Email, &user.Created)
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO USERS (name,email,hashed_password,created) 
VALUES(?,?,?,UTC_TIMESTAMP())
  `
	_, err = m.Db.Exec(stmt, name, email, hashedPassword)

	if err != nil {
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			if mySqlError.Number == 1062 && strings.Contains(mySqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.Db.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials

		} else {

			return 0, err

		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exist bool
	stmt := `SELECT EXISTS(SELECT true from  Users where id =?)`
	err := m.Db.QueryRow(stmt, id).Scan(&exist)
	return exist, err
}
func (m *UserModel) ChangePassword(id int, oldPassword string, newPassword string) (bool, error) {
	var hashedPassword []byte
	stmt := `SELECT hashed_password FROM users WHERE id = ?`

	err := m.Db.QueryRow(stmt, id).Scan(&hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false,ErrNoRecord 
		} else {
			return false, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(oldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, ErrInvalidCredentials
		} else {
			return false, err
		}
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return false, err
	}
	updateStmt := `UPDATE USERS 
SET hashed_password =? 
WHERE ID = ?  `

	_, err = m.Db.Exec(updateStmt, newHashedPassword, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
