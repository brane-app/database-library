package monkebase

import (
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
	statement = "REPLACE INTO " + table + " (" + strings.Join(keys, ", ") + ") VALUES " + "(" + manyParamString(len(keys)) + ")"

	return
}

func manyParamString(size int) (param_string string) {
	param_string = strings.Join(strings.Split(strings.Repeat("?", size), ""), ", ")
	return
}

func interfaceStrings(them ...string) (faces []interface{}) {
	faces = make([]interface{}, len(them))

	var index int
	for index, _ = range them {
		faces[index] = them[index]
	}

	return
}
