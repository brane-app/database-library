package monkebase

import (
	"github.com/google/uuid"
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
	if database == nil {
		panic("database nil after being set!")
	}

	var err error
	var table string
	for table = range tables {
		if _, err = database.Query("DROP TABLE IF EXISTS " + table); err != nil {
			panic(err)
		}
	}

	create()

	var result int = main.Run()

	for table, _ = range tables {
		if err = EmptyTable(table); err != nil {
			panic(err)
		}
	}

	os.Exit(result)
}

func Test_Connect_malformedAddress(test *testing.T) {
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

func Test_Connect_unreachableAddress(test *testing.T) {
	defer func(test *testing.T) {
		var recovered interface{}

		if recovered = recover(); recovered == nil {
			test.Errorf("recover recovered nil!")
		}
	}(test)

	var existing *sqlx.DB = database
	defer func(existing *sqlx.DB) { database = existing }(existing)

	Connect("foo:bar@tcp(nothing)/table")
}

func Test_EmptyTable(test *testing.T) {
	var writable map[string]interface{} = mapCopy(writableUser)
	writable["id"] = uuid.New().String()

	var err error
	if err = WriteUser(writable); err != nil {
		test.Fatal(err)
	}

	if err = EmptyTable(USER_TABLE); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if _, exists, err = ReadSingleUser(writable["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("table %s was not emptied!", USER_TABLE)
	}
}
