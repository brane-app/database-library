package monkebase

import (
	"fmt"
	"strings"
)

func getSQLParams(it map[string]interface{}) (keys []string, values []interface{}) {
	var size int = len(it)
	keys, values = make([]string, size), make([]interface{}, size)

	var index int = 0
	var key string
	var value interface{}
	for key, value = range it {
		keys[index] = key
		values[index] = value
		index++
	}

	return
}

func makeSQLInsertable(table string, it map[string]interface{}) (statement string, values []interface{}) {
	var keys []string
	keys, values = getSQLParams(it)

	var q_marks []string = strings.Split(strings.Repeat("?", len(keys)), "")
	statement = fmt.Sprintf(
		"REPLACE INTO %s (%s) VALUES (%s)",
		table, strings.Join(keys, ", "), strings.Join(q_marks, ", "),
	)

	return
}
