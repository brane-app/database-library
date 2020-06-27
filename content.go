package monkebase

import (
	"time"
)

func WriteContent(content map[string]interface{}) (err error) {
	var copied map[string]interface{} = mapCopy(content)

	if err = setTags(copied["id"].(string), copied["tags"].([]string)); err != nil {
		return
	}

	delete(copied, "tags")

	var query string
	var values []interface{}
	query, values = makeSQLInsertable(CONTENT_TABLE, copied)

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
	if len(tags) < 1 {
		err = dropTags(ID)
		return
	}

	var p_string = manyParamString(len(tags))
	var statement string = "DELETE FROM " + TAG_TABLE + " WHERE id=? AND tag NOT IN (" + p_string + ")"
	var faces []interface{} = append(
		[]interface{}{ID},
		interfaceStrings(tags...)...,
	)

	if _, err = database.Query(statement, faces...); err != nil {
		return
	}

	statement = "REPLACE INTO " + TAG_TABLE + " (id, tag, created) VALUES (?, ?, ?)"
	var now int64 = time.Now().Unix()

	var tag interface{}
	for _, tag = range tags {
		if _, err = database.Query(statement, ID, tag, now); err != nil {
			break
		}
	}

	return
}

func dropTags(ID string) (err error) {
	_, err = database.Query("DELETE FROM "+TAG_TABLE+" WHERE id=?", ID)
	return
}
