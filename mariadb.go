package monkebase

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

var (
	database *sql.DB
)

const (
	USER_TABLE    = "user"
	CONTENT_TABLE = "content"
)

func connect(address string) {
	var err error
	if database, err = sql.Open("mysql", address); err != nil {
		panic(err)
	}

	if err = database.Ping(); err != nil {
		panic(err)
	}
}
