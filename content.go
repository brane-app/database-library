package monkebase

import (
	"github.com/imonke/monketype"
	"github.com/jmoiron/sqlx"

	"database/sql"
	"time"
)

/**
 * Write some content `content` to the table CONTENT_TABLE
 * Uses 3 query
 * 		write content: 	REPLACE INTO CONTENT_TABLE (keys...) VALUES (values...)
 * 		queries from: setTags
 * Returns error, if any
 */
func WriteContent(content map[string]interface{}) (err error) {
	var copied map[string]interface{} = mapCopy(content)

	if err = setTags(copied["id"].(string), copied["tags"].([]string)); err != nil {
		return
	}

	delete(copied, "tags")

	var statement string
	var values []interface{}
	statement, values = makeSQLInsertable(CONTENT_TABLE, copied)
	_, err = database.Exec(statement, values...)
	return
}

/**
 * Delete some content of id `ID`
 * Uses 2 queries
 * 		delete content:		DELETE FROM CONTENT_TABLE WHERE id=ID LIMIT 1
 * 		queries from: 		dropTags
 */
func DeleteContent(ID string) (err error) {
	if _, err = database.Exec(DELETE_CONTENT_ID, ID); err == nil {
		err = dropTags(ID)
	}
	return
}

/**
 * Read some content of id `ID`
 * Uses 2 queries
 * 		get content: 	SELECT * FROM CONTENT_TABLE WHERE id=ID LIMIT 1
 * 		get tags:		SELECT tag FROM TAG_TABLE WHERE id=ID
 */
func ReadSingleContent(ID string) (content monketype.Content, exists bool, err error) {
	if err = database.QueryRowx(READ_CONTENT_ID, ID).StructScan(&content); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}

	exists = true
	content.Tags, err = getTags(ID)
	return
}

/**
 * Read `count` number of contents, starting from `offset`
 * Newest posts are returned first
 * Uses 2 queries
 * 		get content: 	SELECT * FROM CONTENT_TABLE ORDER BY created DESC LIMIT offset, count
 * 		queries from: 	getManyTags
 */
func DEPRECATED_ReadManyContent(offset, count int) (content []monketype.Content, size int, err error) {
	var rows *sqlx.Rows
	if rows, err = database.Queryx(DEPRECATED_READ_MANY_CONTENT, offset, count); err != nil || rows == nil {
		return
	}

	defer rows.Close()
	content, size, err = scanManyContent(rows, count)
	return
}

/**
 * Same as DEPRECATED_ReadManyContent but for some author of id `ID`
 * Uses 2 queries
 * 		get content: 	SELECT * FROM CONTENT_TABLE ORDER BY created DESC LIMIT offset, count
 * 		queries from: 	getManyTags
 */
func DEPRECATED_ReadAuthorContent(ID string, offset, count int) (content []monketype.Content, size int, err error) {
	var rows *sqlx.Rows
	if rows, err = database.Queryx(DEPRECATED_READ_MANY_CONTENT_OF_AUTHOR, ID, offset, count); err != nil || rows == nil {
		return
	}

	defer rows.Close()
	content, size, err = scanManyContent(rows, count)
	return
}

/**
 * Read some tags for post of id `ID`
 * Uses 1 query
 * 		get tags:	SELECT tag FROM TAG_TABLE WHERE id=ID
 */
func getTags(ID string) (tags []string, err error) {
	var rows *sqlx.Rows
	if rows, err = database.Queryx(READ_TAGS_OF_ID, ID); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	tags = make([]string, 0)
	var tag string
	for rows.Next() {
		if err = rows.Scan(&tag); err != nil {
			break
		}

		tags = append(tags, tag)
	}

	return
}

/**
 * Get the tags for every post of id in `IDs`
 * Returns a map where
 * 		id -> []tags
 * Uses 1 query:
 * 		get tags: SELECT id, tag FROM TAG_TABLE WHERE id IN (IDs...)
 */
func getManyTags(IDs []string) (tags map[string][]string, err error) {
	var size int = len(IDs)
	if size < 1 {
		return
	}

	tags = make(map[string][]string, size)
	var id string
	for _, id = range IDs {
		tags[id] = make([]string, 0)
	}

	var paramString string = "(" + manyParamString("?", len(IDs)) + ")"
	var rows *sql.Rows
	if rows, err = database.Query(READ_TAGS_OF_MANY_ID+paramString, interfaceStrings(IDs...)...); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	var tag string
	for rows.Next() {
		if err = rows.Scan(&id, &tag); err != nil {
			break
		}

		tags[id] = append(tags[id], tag)
	}

	return
}

/**
 * Updates the tags of a post
 * Done in two queries if there are tags
 * Or one if there are no tags
 */
func setTags(ID string, tags []string) (err error) {
	var length int = len(tags)
	if length < 1 {
		err = dropTags(ID)
		return
	}

	var paramString string = "(" + manyParamString("?", len(tags)) + ")"
	var faces []interface{} = append(
		[]interface{}{ID},
		interfaceStrings(tags...)...,
	)
	if _, err = database.Exec(DELETE_STALE_TAGS_OF_ID+paramString, faces...); err != nil {
		return
	}

	var now int64 = time.Now().Unix()
	var insertable []interface{} = make([]interface{}, length*3)

	var index int = 0
	for index < length {
		insertable[index*3] = ID
		insertable[index*3+1] = tags[index]
		insertable[index*3+2] = now
		index++
	}

	paramString = manyParamString("(?, ?, ?)", length)
	_, err = database.Exec(WRITE_TAGS_OF_MANY_ID+paramString, insertable...)
	return
}

/**
 * Deletes all of the tags for some post
 * Done in one query
 */
func dropTags(ID string) (err error) {
	_, err = database.Exec(DELETE_TAGS_OF_ID, ID)
	return
}
