package monkebase

func WriteContent(content interface{}) (err error) {
	_, err = database.Query("")
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
