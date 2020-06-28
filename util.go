package monkebase

import (
	"database/sql"
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
	statement = "REPLACE INTO " + table + " (" + strings.Join(keys, ", ") + ") VALUES " + "(" + manyParamString("?", len(keys)) + ")"

	return
}

func manyParamString(param string, size int) (param_string string) {
	var param_slice []string = make([]string, size)
	for size != 0 {
		size--
		param_slice[size] = param
	}

	param_string = strings.Join(param_slice, ", ")
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

func mapCopy(source map[string]interface{}) (copy map[string]interface{}) {
	copy = map[string]interface{}{}

	var key string
	var value interface{}
	for key, value = range source {
		copy[key] = value
	}

	return
}
