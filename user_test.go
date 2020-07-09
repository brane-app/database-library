package monkebase

import (
	"github.com/google/uuid"
	"github.com/imonke/monketype"

	"testing"
)

var (
	user         monketype.User         = monketype.NewUser("imonke", "mmm, monke", "me@imonke.io")
	writableUser map[string]interface{} = user.Map()
)

func userOK(test *testing.T, data map[string]interface{}, have monketype.User) {
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

	var user monketype.User
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

	var user monketype.User
	var exists bool
	var err error
	if user, exists, err = ReadSingleUser(id); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("Query for nonexisting id got %+v", user)
	}
}

func Test_DeleteUser(test *testing.T) {
	var mod map[string]interface{} = map[string]interface{}{
		"id":    uuid.New().String(),
		"email": "delete@monke.io",
	}

	var err error
	if err = WriteUser(mapMod(writableUser, mod)); err != nil {
		test.Fatal(err)
	}

	if err = DeleteUser(mod["id"].(string)); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if _, exists, err = ReadSingleUser(mod["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("user %s exists after being deleted!", mod["id"])
	}
}
