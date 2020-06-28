package monkebase

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
func ReadSingleUser(ID string) (user User, exists bool, err error) {
	var statement string = "SELECT * FROM " + USER_TABLE + " WHERE id=? LIMIT 1"

	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, ID); err != nil {
		return
	}

	if exists = rows.Next(); exists {
		err = rows.StructScan(&user)
	}

	return
}
