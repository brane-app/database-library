package monkebase

import (
	"github.com/google/uuid"

	"testing"
)

func Test_SetPassword(test *testing.T) {
	var password string = randStringOfSize(64)
	var id string = uuid.New().String()

	var err error
	if err = SetPassword(id, password); err != nil {
		test.Fatal(err)
	}

	var ok bool
	if ok, err = CheckPassword(id, password); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Errorf("Set password %s does not match retrieved!", password)
	}
}

func Test_SetPassword_length(test *testing.T) {
	var id string = uuid.New().String()
	var err error

	var index int = 1
	for index != 4*64 {
		index = index * 4
		if err = SetPassword(id, randStringOfSize(index)); err != nil {
			test.Fatal(err)
		}
	}
}

func Test_CheckPassword_wrong(test *testing.T) {
	var sets []string = []string{
		"password",
		"",
	}

	var id string = uuid.New().String()
	var err error
	if err = SetPassword(id, randStringOfSize(64)); err != nil {
		test.Fatal(err)
	}

	var ok bool
	var set string
	for _, set = range sets {
		if ok, err = CheckPassword(id, set); err != nil {
			test.Fatal(err)
		}

		if ok {
			test.Errorf("password %s should not match, but does!", set)
		}
	}
}
