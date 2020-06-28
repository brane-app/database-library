package monkebase

func WriteUser(user map[string]interface{}) (err error) {
	var statement string
	var values []interface{}
	statement, values = makeSQLInsertable(USER_TABLE, user)

	_, err = database.Query(statement, values...)
	return
}

func ReadSingleUser(ID string) (user interface{}, exists bool, err error) {
	return
}
