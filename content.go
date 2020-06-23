package monkebase

type Content interface{}

func WriteContent(content Content) (err error) {
	return
}

func ReadSingleContent(ID string) (content Content, exists bool, err error) {
	return
}

func ReadManyContent(index, limit int) (content []Content, size int, err error) {
	return
}

func ReadAuthorContent(ID string, index, limit int) (content []Content, size, err error) {
	return
}
