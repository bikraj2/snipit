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
  ID int
  Name string
  Email string
  HashedPassword []byte
  Created time.Time
}


type UserModel struct {
  Db *sql.DB
}
func (m *UserModel) Insert(name,email,password string) error {
  hashedPassword,err:=bcrypt.GenerateFromPassword([]byte(password), 12 )
  if err!=nil {
    return err
  }

  stmt:= `INSERT INTO USERS (name,email,hashed_password,created) 
VALUES(?,?,?,UTC_TIMESTAMP())
  `
  _,err =m.Db.Exec(stmt,name,email,hashedPassword)

if err != nil {
	var mySqlError *mysql.MySQLError
    if errors.As(err,  &mySqlError) {
      if mySqlError.Number == 1062 && strings.Contains(mySqlError.Message,  "users_uc_email") {
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
return id, nil }

func(m *UserModel) Exists(id int) (bool,error) {
  var exist bool
  stmt:= `SELECT EXISTS(SELECT true from  Users where id =?)`
  err:= m.Db.QueryRow(stmt, id).Scan( &exist )
  return exist,err
  }


