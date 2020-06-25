package monkebase

import (
	"database/sql"
	"os"
	"testing"
)

var (
	CONNECTION string = os.Getenv("MONKEBASE_CONNECTION")
)

func TestMain(main *testing.T) {
	connect(CONNECTION)

	var table string
	for table = range tables {
		database.Query("DROP TABLE IF EXISTS ?", table)
	}

	create()
}

func Test_connect(test *testing.T) {
	defer func(test *testing.T) {
		var recovered interface{}

		if recovered = recover(); recovered == nil {
			test.Errorf("recover recovered nil!")
		}
	}(test)

	var existing *sql.DB = database
	defer func(existing *sql.DB) { database = existing }(existing)

	connect("foobar")
}
