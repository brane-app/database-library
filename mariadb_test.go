package monkebase

import (
	"github.com/jmoiron/sqlx"

	"os"
	"testing"
)

var (
	CONNECTION string = os.Getenv("MONKEBASE_CONNECTION")
)

func mapMod(source map[string]interface{}, mods ...map[string]interface{}) (modified map[string]interface{}) {
	modified = mapCopy(source)

	var key string
	var value interface{}

	var mod map[string]interface{}
	for _, mod = range mods {
		for key, value = range mod {
			modified[key] = value
		}
	}

	return
}

func TestMain(main *testing.M) {
	Connect(CONNECTION)

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

func _Test_Connect(test *testing.T) {
	defer func(test *testing.T) {
		var recovered interface{}

		if recovered = recover(); recovered == nil {
			test.Errorf("recover recovered nil!")
		}
	}(test)

	var existing *sqlx.DB = database
	defer func(existing *sqlx.DB) { database = existing }(existing)

	Connect("foobar")
}
