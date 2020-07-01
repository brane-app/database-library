package monkebase

import (
	"github.com/google/uuid"

	"testing"
)

func Test_SetPassword(test *testing.T) {
	var password string
	var err error
	if password, err = randomString(64); err != nil {
		test.Fatal(err)
	}

	var id string = uuid.New().String()

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
	var password string
	var index int = 1
	for index != 4*64 {
		index = index * 4
		if password, err = randomString(index); err != nil {
			test.Fatal(err)
		}
		if err = SetPassword(id, password); err != nil {
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
	var password string
	if password, err = randomString(64); err != nil {
		test.Fatal(err)
	}

	if err = SetPassword(id, password); err != nil {
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

func Test_CreateSecret(test *testing.T) {
	var id string = uuid.New().String()

	var err error
	var secret string
	if secret, err = CreateSecret(id); err != nil {
		test.Fatal(err)
	}

	_ = secret
}

func Test_CheckSecret(test *testing.T) {
	var id string = uuid.New().String()

	var secret string
	var err error
	if secret, err = CreateSecret(id); err != nil {
		test.Fatal(err)
	}

	var valid bool
	if valid, err = CheckSecret(id, secret); err != nil {
		test.Fatal(err)
	}

	if !valid {
		test.Errorf("Just set secret %s is invalid for %s!", secret, id)
	}
}

func Test_CheckSecret_invalid(test *testing.T) {
	var id string = uuid.New().String()

	var secret string
	var err error
	if secret, err = CreateSecret(id); err != nil {
		test.Fatal(err)
	}

	var valid bool
	if valid, err = CheckSecret(uuid.New().String(), secret); err == nil {
		test.Errorf("secret for random uuid returned no err!")
	}

	if valid {
		test.Errorf("Just set secret %s is valid for a random uuid!", secret)
	}

	if valid, err = CheckSecret(id, "not_a_secret"); err != nil {
		test.Fatal(err)
	}

	if valid {
		test.Errorf("A bad secret is valid for uuid %s!", id)
	}
}

func Test_RevokeSecretOf(test *testing.T) {
	var id string = uuid.New().String()

	var secret string
	var err error
	if secret, err = CreateSecret(id); err != nil {
		test.Fatal()
	}

	var valid bool
	if valid, err = CheckSecret(id, secret); err != nil {
		test.Fatal(err)
	}

	if !valid {
		test.Errorf("Just set secret %s is invalid for %s!", secret, id)
	}

	if err = RevokeSecretOf(id); err != nil {
		test.Fatal(err)
	}

	if valid, err = CheckSecret(id, secret); err == nil {
		test.Errorf("secret for revoked uuid returned no err!")
	}

	if valid {
		test.Errorf("revoked secret %s for %s is still valid!", secret, id)
	}
}
