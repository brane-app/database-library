package monkebase

import (
	"github.com/google/uuid"

	"testing"
	"time"
)

var (
	writableUser map[string]interface{} = map[string]interface{}{
		"id":                 uuid.New().String(),
		"email":              "me@imonke.io",
		"nick":               "imonke",
		"bio":                "mmmm, monke",
		"subscriber_count":   0,
		"subscription_count": 0,
		"post_count":         0,
		"created":            time.Now().Unix(),
	}
)

func Test_WriteUser(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{},
		map[string]interface{}{
			"id":  uuid.New().String(),
			"bio": "' or 1=1; DROP TABLE content",
		},
	}

	var err error
	var mod map[string]interface{}
	for _, mod = range mods {
		mod = mapMod(writableUser, mod)
		if err = WriteUser(mod); err != nil {
			test.Fatal(err)
		}
	}
}

func Test_WriteUser_err(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{
			"id":  uuid.New().String(),
			"bio": nil,
		},
		map[string]interface{}{
			"id":     uuid.New().String(),
			"answer": 42,
		},
	}

	var mod map[string]interface{} = mapCopy(writableUser)
	delete(mod, "id")
	mod["bio"] = "foobar"

	var err error
	if err = WriteUser(mod); err == nil {
		test.Errorf("data %+v produced no error!", mod)
	}

	for _, mod = range mods {
		mod = mapMod(writableUser, mod)
		if err = WriteUser(mod); err == nil {
			test.Errorf("data %+v produced no error!", mod)
		}
	}

}
