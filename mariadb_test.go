package monkebase

import (
	"github.com/jmoiron/sqlx"
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

	var err error

	var table string
	for table = range tables {
		if _, err = database.Query("DROP TABLE IF EXISTS " + table); err != nil {
			panic(err)
		}
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

	var existing *sqlx.DB = database
	defer func(existing *sqlx.DB) { database = existing }(existing)

	connect("foobar")
}
