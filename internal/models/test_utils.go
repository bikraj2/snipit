package models

import (
	"database/sql"
	"os"
	"testing"
)
func newTestDB(t *testing.T) *sql.DB {
  db,err:=sql.Open("mysql", "root:bikrajadmin@tcp(localhost:3306)/?parseTime=true")
  if err != nil {
    t.Fatal(err)
  }
script,err := os.ReadFile("./testdata/setup.sql")
  if err != nil {
    t.Fatal(err)
  }

  _, err = db.Exec(string(script))
  if err != nil {
    t.Fatal(err)
  }

  t.Cleanup(func() {
    script,err := os.ReadFile("./testdata/teardown.sql")
    if err!=nil {
      t.Fatal(err)
    }
    _,err = db.Exec(string(script), )
    if err!=nil {
      t.Fatal(err)
    }
    db.Close()
  })
 return db 
}
