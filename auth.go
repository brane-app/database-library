package monkebase

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"
)

const (
	BCRYPT_ITERS  = 12
	TOKEN_LENGTH  = 24
	SECRET_LENGTH = 128
	TOKEN_TTL     = 60 * 60 * 24
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

func CreateSecret(ID string) (secret string, err error) {
	var bytes []byte
	if bytes, err = randomBytes(SECRET_LENGTH); err != nil {
		return
	}

	secret = base64.URLEncoding.EncodeToString(bytes)

	var statement string = "REPLACE INTO " + SECRET_TABLE + " (id, secret) VALUES (?, ?)"
	_, err = database.Exec(statement, ID, bytes)

	return
}

func CheckSecret(ID, secret string) (valid bool, err error) {
	var statement string = "SELECT secret FROM " + SECRET_TABLE + " WHERE id=? LIMIT 1"
	var bytes []byte
	if err = database.QueryRowx(statement, ID).Scan(&bytes); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			valid = false
		}

		return
	}

	valid = secret == base64.URLEncoding.EncodeToString(bytes)
	return
}

func RevokeSecretOf(ID string) (err error) {
	var statement string = "DELETE FROM " + SECRET_TABLE + " WHERE id=? LIMIT 1"
	_, err = database.Exec(statement, ID)
	return
}

func CreateToken(ID string) (token string, expires int64, err error) {
	var bytes []byte
	if bytes, err = randomBytes(TOKEN_LENGTH); err != nil {
		return
	}

	var now int64 = time.Now().Unix()
	expires = now + TOKEN_TTL
	token = base64.URLEncoding.EncodeToString(bytes)

	var statement string = "REPLACE INTO " + TOKEN_TABLE + " (id, token, created) VALUES (?, ?, ?)"
	_, err = database.Exec(statement, ID, bytes, now)
	return
}

func ReadTokenStat(token string) (owner string, valid bool, err error) {
	var bytes []byte
	if bytes, err = base64.URLEncoding.DecodeString(token); err != nil {
		err = nil
		valid = false
		return
	}

	var statement string = "SELECT id, created FROM " + TOKEN_TABLE + " WHERE token=? LIMIT 1"
	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, bytes); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	if !rows.Next() {
		return
	}

	var created int64
	if err = rows.Scan(&owner, &created); err == nil {
		valid = created <= time.Now().Unix() && created+TOKEN_TTL >= time.Now().Unix()
	}

	return
}

func RevokeToken(token string) (err error) {
	var bytes []byte
	if bytes, err = base64.URLEncoding.DecodeString(token); err != nil {
		return
	}

	var statement string = "DELETE FROM " + TOKEN_TABLE + " WHERE token=?"
	_, err = database.Exec(statement, bytes)
	return
}

func RevokeTokenOf(ID string) (err error) {
	var statement string = "DELETE FROM " + TOKEN_TABLE + " WHERE id=?"
	_, err = database.Exec(statement, ID)
	return
}

func CheckPassword(ID, password string) (valid bool, err error) {
	var statement string = "SELECT hash FROM " + AUTH_TABLE + " WHERE id=? LIMIT 1"
	var hash []byte
	if err = database.QueryRowx(statement, ID).Scan(&hash); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			valid = false
		}

		return
	}

	valid = bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
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
