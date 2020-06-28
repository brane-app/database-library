package monkebase

import (
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

	var rows *sql.Rows
	if rows, err = database.Query(statement, values...); rows != nil {
		rows.Close()
	}
	return
}

/**
 * Read some content of id `ID`
 * Uses 2 queries
 * 		get content: 	SELECT * FROM CONTENT_TABLE WHERE id=ID LIMIT 1
 * 		get tags:		SELECT tag FROM TAG_TABLE WHERE id=ID
 */
func ReadSingleContent(ID string) (content Content, exists bool, err error) {
	var statement string = "SELECT * FROM " + CONTENT_TABLE + " WHERE id=? LIMIT 1"

	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, ID); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	if exists = rows.Next(); !exists {
		return
	}

	if err = rows.StructScan(&content); err != nil {
		return
	}

	content.Tags, err = getTags(ID)

	return
}

func ReadManyContent(offset, count int) (content []Content, size int, err error) {
	var statement string = "SELECT * FROM " + CONTENT_TABLE + " ORDER BY created DESC LIMIT ?, ?"
	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, offset, count); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	var ids []string = make([]string, count)
	content = make([]Content, count)
	size = 0

	for rows.Next() {
		rows.StructScan(&content[size])
		ids[size] = content[size].ID
		size++
	}

	content = content[:size]

	var tags map[string][]string
	if tags, err = getManyTags(ids); err != nil {
		return
	}

	var index int
	for index, _ = range content {
		content[index].Tags = tags[content[index].ID]
	}

	return
}

func ReadAuthorContent(ID string, index, limit int) (content []interface{}, size, err error) {
	return
}

/**
 * Read some tags for post of id `ID`
 * Uses 1 query
 * 		get tags:	SELECT tag FROM TAG_TABLE WHERE id=ID
 */
func getTags(ID string) (tags []string, err error) {
	var statement string = "SELECT tag FROM " + TAG_TABLE + " WHERE id=?"
	var rows *sqlx.Rows
	if rows, err = database.Queryx(statement, ID); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	var tag string
	for rows.Next() {
		if err = rows.Scan(&tag); err != nil {
			break
		}

		tags = append(tags, tag)
	}

	return
}

func getManyTags(IDs []string) (tags map[string][]string, err error) {
	if len(IDs) < 1 {
		return
	}

	var statement string = "SELECT id, tag FROM " + TAG_TABLE + " WHERE id IN (" + manyParamString("?", len(IDs)) + ")"

	var rows *sql.Rows
	if rows, err = database.Query(statement, interfaceStrings(IDs...)...); err != nil || rows == nil {
		return
	}

	defer rows.Close()
	tags = make(map[string][]string, len(IDs))

	var id, tag string
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
 * Uses 2 queries
 * 		delete missing: 	DELETE FROM TAG_TABLE WHERE id=ID AND tag NOT IN (tags...)
 * 		update tags:		REPLACE INTO TAG_TABLE (id, tag, created) VALUES (ID, tags..., now)
 * Or if there are no tags
 * 		queries from dropTags
 */
func setTags(ID string, tags []string) (err error) {
	var length int = len(tags)
	if length < 1 {
		err = dropTags(ID)
		return
	}

	var statement string = "DELETE FROM " + TAG_TABLE + " WHERE id=? AND tag NOT IN (" + manyParamString("?", len(tags)) + ")"
	var faces []interface{} = append(
		[]interface{}{ID},
		interfaceStrings(tags...)...,
	)

	var rows *sql.Rows
	if rows, err = database.Query(statement, faces...); err != nil || rows == nil {
		return
	}

	defer rows.Close()

	var now int64 = time.Now().Unix()
	var insertable []interface{} = make([]interface{}, length*3)

	var index int = 0
	for index < length {
		insertable[index*3] = ID
		insertable[index*3+1] = tags[index]
		insertable[index*3+2] = now
		index++
	}

	statement = "REPLACE INTO " + TAG_TABLE + " (id, tag, created) VALUES " + manyParamString("(?, ?, ?)", length)
	_, err = database.Query(statement, insertable...)
	return
}

/**
 * Deletes all of the tags for some post
 * Uses 1 query:
 * 		DELETE FROM TAG_TABLE WHERE id=ID
 */
func dropTags(ID string) (err error) {
	var rows *sql.Rows
	if rows, err = database.Query("DELETE FROM "+TAG_TABLE+" WHERE id=?", ID); rows != nil {
		rows.Close()
	}
	return
}
