package monkebase

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"math/rand"
)

var (
	randable []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

const (
	BCRYPT_ITERS  = 12
	TOKEN_LENGTH  = 24
	SECRET_LENGTH = 128
	RAND_LENGTH   = 62
)

func randStringOfSize(size int) (generated string) {
	var items []rune = make([]rune, size)

	var index int
	for index = range items {
		items[index] = randable[rand.Intn(RAND_LENGTH)]
	}

	generated = string(items)
	return
}

func ReadTokenStat(token string) (owner string, valid bool, err error) {
	return
}

func CreateToken(ID string) (token string, expires int64, err error) {
	return
}

func CreateSecret(ID string) (secret string, err error) {
	return
}

func RevokeToken(token string) (err error) {
	return
}

func RevokeSecret(secret string) (err error) {
	return
}

func RevokeTokenOf(token string) (err error) {
	return
}

func RevokeSecretOf(secret string) (err error) {
	return
}

func CheckPassword(ID, password string) (ok bool, err error) {
	var statement string = "SELECT hash FROM " + AUTH_TABLE + " WHERE id=? LIMIT 1"
	var hash []byte

	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, ID); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	rows.Next()
	if err = rows.Scan(&hash); err != nil {
		return
	}

	ok = bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
	return
}

func SetPassword(ID, password string) (err error) {
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(password), BCRYPT_ITERS); err != nil {
		return
	}

	var statement string = "REPLACE INTO " + AUTH_TABLE + " (id, hash) VALUES (?, ?)"
	_, err = database.Exec(statement, ID, hash)

	return
}
