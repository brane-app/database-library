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

func userOK(test *testing.T, data map[string]interface{}, have User) {
	if data["id"].(string) != have.ID {
		test.Errorf("User ID mismatch! have: %s, want: %s", have.ID, data["id"])
	}

	if data["bio"].(string) != have.Bio {
		test.Errorf("User bio mismatch! have: %s, want: %s", have.Bio, data["bio"])
	}
}

func Test_WriteUser(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{},
		map[string]interface{}{
			"id":  uuid.New().String(),
			"bio": "' or 1=1; DROP TABLE user",
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

func Test_ReadSingleUser(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableUser)
	modified["id"] = uuid.New().String()

	WriteUser(modified)

	var user User
	var exists bool
	var err error
	if user, exists, err = ReadSingleUser(modified["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("user of id %s does not exist!", modified["id"])
	}

	userOK(test, modified, user)
}

func Test_ReadSingleUser_NotExists(test *testing.T) {
	var id string = uuid.New().String()

	var user User
	var exists bool
	var err error
	if user, exists, err = ReadSingleUser(id); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("Query for nonexisting id got %+v", user)
	}
}
