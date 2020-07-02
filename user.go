package monkebase

import (
	"github.com/imonke/monketype"
	"github.com/jmoiron/sqlx"
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
 * Read some user of id `ID` from USER_TABLE
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE id=ID LIMIT 1
 */
func ReadSingleUser(ID string) (user monketype.User, exists bool, err error) {
	var statement string = "SELECT * FROM " + USER_TABLE + " WHERE id=? LIMIT 1"

	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, ID); err != nil || rows == nil{
		return
	}

	defer rows.Close()

	if exists = rows.Next(); exists {
		err = rows.StructScan(&user)
	}

	return
}

/**
 * Read some user of email `email` from USER_TABLE
 * Works in the same way as ReadSingleUser, but with email
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE email=email LIMIT 1
 */
func ReadSingleUserEmail(email string) (user monketype.User, exists bool, err error) {
	var statement string = "SELECT * FROM "+USER_TABLE+" WHERE email=? LIMIT 1"

	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, email); err != nil || rows == nil {
		return
	}

	defer rows.Close()
	if exists = rows.Next(); exists {
		err = rows.StructScan(&user)
	}

	return
}

/**
 * Read some user of email `email` from USER_TABLE
 * Works in the same way as ReadSingleUser, but with nick
 * Uses 1 query
 * 		read user: 	SELECT * FROM USER_TABLE WHERE nick=nick LIMIT 1
 */
func ReadSingleUserNick(nick string) (user monketype.User, exists bool, err error) {
	var statement string = "SELECT * FROM "+USER_TABLE+" WHERE nick=? LIMIT 1"

	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, nick); err != nil || rows == nil {
		return
	}

	defer rows.Close()
	if exists = rows.Next(); exists {
		err = rows.StructScan(&user)
	}

	return
}
