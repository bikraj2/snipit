package main
import "database/sql"

type ExampleModel struct {
	DB *sql.DB
}

func (app *ExampleModel) ExampleTransc() error {
	tx, err := app.DB.Begin()

	se
} 
