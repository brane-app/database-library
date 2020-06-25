package monkebase

func WriteContent(content Mappable) (err error) {
	var mapped map[string]interface{} = content.Map()
	if err = setTags(mapped["id"].(string), mapped["tags"].([]string)); err != nil {
		return
	}

	delete(mapped, "tags")

	var query string
	var values []interface{}
	query, values = makeSQLInsertable(CONTENT_TABLE, mapped)

	_, err = database.Query(query, values...)
	return
}

func ReadSingleContent(ID string) (content interface{}, exists bool, err error) {
	return
}

func ReadManyContent(index, limit int) (content []interface{}, size int, err error) {
	return
}

func ReadAuthorContent(ID string, index, limit int) (content []interface{}, size, err error) {
	return
}

func setTags(ID string, tags []string) (err error) {
	var p_string = manyParamString(len(tags))
	var statement string = "DELETE FROM " + TAG_TABLE + " WHERE id=? AND tag NOT IN (" + p_string + ")"
	var faces []interface{} = append(
		[]interface{}{ID},
		interfaceStrings(tags...)...,
	)

	if _, err = database.Query(statement, faces...); err != nil {
		return
	}

	statement = "REPLACE INTO " + TAG_TABLE + " (id, tag) VALUES (?, ?)"

	var tag interface{}
	for _, tag = range tags {
		if _, err = database.Query(statement, ID, tag); err != nil {
			break
		}
	}

	return
}
