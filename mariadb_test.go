package monkebase

import (
	"database/sql"
	"os"
	"testing"
)

var (
	CONNECTION string = os.Getenv("MONKEBASE_CONNECTION")
)

type testWritable struct {
	Data map[string]interface{}
}

func (writable testWritable) Map() (data map[string]interface{}) {
	return writable.Data
}

func mapCopy(source map[string]interface{}) (copy map[string]interface{}) {
	copy = map[string]interface{}{}

	var key string
	var value interface{}
	for key, value = range source {
		copy[key] = value
	}

	return
}

func TestMain(main *testing.M) {
	connect(CONNECTION)

	var table string
	for table = range tables {
		database.Query("DROP TABLE IF EXISTS ?", table)
	}

	create()
	if database == nil {
		panic("database nil after being set!")
	}

	os.Exit(main.Run())
}

func _Test_connect(test *testing.T) {
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
