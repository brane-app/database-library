package monkebase

import (
	"github.com/imonke/monketype"

	"database/sql"
)

/**
 * Write some user `user` into USER_TABLE
 * Uses 1 query
 * 		write user: 	REPLACE INTO CONTENT_TABLE (keys...) VALUES (values...)
 * Returns error, if any
 */
func WriteUser(user map[string]interface{}) (err error) {
	var statement string
	var values []interface{}
	statement, values = makeSQLInsertable(USER_TABLE, user)

	_, err = database.Query(statement, values...)
	return
}

/**
 * Delete some user from USER_TABLE
 * Uses 1 query:
 * 		delete user: 	DELETE FROM USER_TABLE WHERE id=ID LIMIT 1
 */
func DeleteUser(ID string) (err error) {
	var statement string = "DELETE FROM " + USER_TABLE + " WHERE id=? LIMIT 1"
	_, err = database.Exec(statement, ID)
	return
}

func readSingleUserKey(key, query string) (user monketype.User, exists bool, err error) {
	var statement string = "SELECT * FROM " + USER_TABLE + " WHERE " + key + "=? LIMIT 1"
	if err = database.QueryRowx(statement, query).StructScan(&user); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	exists = true
	return
}

/**
 * Read some user of id `ID` from USER_TABLE
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE id=ID LIMIT 1
 */
func ReadSingleUser(ID string) (user monketype.User, exists bool, err error) {
	user, exists, err = readSingleUserKey("id", ID)
	return
}

/**
 * Read some user of email `email` from USER_TABLE
 * Works in the same way as ReadSingleUser, but with email
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE email=email LIMIT 1
 */
func ReadSingleUserEmail(email string) (user monketype.User, exists bool, err error) {
	user, exists, err = readSingleUserKey("email", email)
	return
}

/**
 * Read some user of email `email` from USER_TABLE
 * Works in the same way as ReadSingleUser, but with nick
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE nick=nick LIMIT 1
 */
func ReadSingleUserNick(nick string) (user monketype.User, exists bool, err error) {
	user, exists, err = readSingleUserKey("nick", nick)
	return
}

func modifyNamedCount(ID, key string, diff int) (err error) {
	var statement string = "UPDATE " + USER_TABLE + " SET " + key + "=" + key + "+? WHERE id=?"
	_, err = database.Exec(statement, diff, ID)
	return
}

/**
 * Increment the post count of user of id `ID` by one
 * Done in one query
 * 		increment: UPDATE USER_TABLE SET post_count=post_count+1 WHERE id=ID
 */
func IncrementPostCount(ID string) (err error) {
	err = modifyNamedCount(ID, "post_count", 1)
	return
}

func IsModerator(ID string) (moderator bool, err error) {
	var admin bool
	var statement string = "SELECT admin, moderator FROM " + USER_TABLE + " WHERE id=?"
	if err = database.QueryRowx(statement, ID).Scan(&admin, &moderator); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}

		return
	}

	moderator = moderator || admin
	return
}

func IsAdmin(ID string) (admin bool, err error) {
	var statement string = "SELECT admin FROM " + USER_TABLE + " WHERE id=?"
	if err = database.QueryRowx(statement, ID).Scan(&admin); err == sql.ErrNoRows {
		err = nil
	}

	return
}

func SetModerator(ID string, state bool) (err error) {
	var statement string = "UPDATE " + USER_TABLE + " SET moderator=? WHERE id=?"
	_, err = database.Exec(statement, state, ID)
	return
}

func SetAdmin(ID string, state bool) (err error) {
	var statement string = "UPDATE " + USER_TABLE + " SET admin=? WHERE id=?"
	_, err = database.Exec(statement, state, ID)
	return
}
