package monkebase

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"encoding/base64"
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

func randomBytes(size int) (generated []byte, err error) {
	generated = make([]byte, size)
	_, err = rand.Read(generated)
	return
}

func randomString(size int) (generated string, err error) {
	var bytes []byte
	if bytes, err = randomBytes(size); err == nil {
		generated = base64.URLEncoding.EncodeToString(bytes)
	}

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
