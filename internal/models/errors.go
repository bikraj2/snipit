package models

import "errors"

var(
  ErrNoRecord = errors.New("models: no matching record Found")
  
  ErrInvalidCredentials = errors.New("modles: invalid credentials")

  ErrDuplicateEmail = errors.New("models: duplicate email")
)
