package monkebase

import (
	"github.com/google/uuid"
	"github.com/imonke/monketype"

	"sort"
	"strconv"
	"testing"
	"time"
)

var (
	content monketype.Content = monketype.NewContent(
		"https://gastrodon.io/file/foobar",
		uuid.New().String(),
		"png",
		[]string{"some", "tags"},
		true, false,
	)
	writableContent map[string]interface{} = content.Map()
)

func contentOK(test *testing.T, data map[string]interface{}, have monketype.Content) {
	if data["id"].(string) != have.ID {
		test.Errorf("monketype.Content ID mismatch! have: %s, want: %s", have.ID, data["id"])
	}

	var tags []string = data["tags"].([]string)
	var length = len(have.Tags)
	if length != len(tags) {
		test.Errorf("Tags mismatch! have: %v, want: %v", have.Tags, tags)
	}

	sort.Strings(tags)
	sort.Strings(have.Tags)

	for length != 0 {
		length--
		if tags[length] != have.Tags[length] {
			test.Errorf("Tags mismatch at %d! have: %v, want: %v", length, have.Tags, tags)
		}
	}
}

func populate(many int) (err error) {
	var modified map[string]interface{}

	for many != 0 {
		modified = mapCopy(writableContent)
		modified["id"] = "many_" + strconv.Itoa(many)
		modified["created"] = time.Now().Unix() + int64(100000*many)

		if err = WriteContent(modified); err != nil {
			break
		}

		many--
	}

	return
}

func Test_WriteContent(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{},
		map[string]interface{}{
			"tags": []string{},
		},
		map[string]interface{}{
			"id":   "0",
			"mime": "' or 1=1; DROP TABLE users",
		},
	}

	var err error
	var mod map[string]interface{}
	for _, mod = range mods {
		mod = mapMod(writableContent, mod)
		if err = WriteContent(mod); err != nil {
			test.Fatal(err)
		}
	}
}

func Test_WriteContent_err(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{
			"id": "96910fdf-916b-4664-a38c-be42d0f2c0ce foobar",
		},
		map[string]interface{}{
			"foo": "bar",
		},
	}

	var mod map[string]interface{}
	mod = mapCopy(writableContent)
	delete(mod, "file_url")

	var err error
	if err = WriteContent(mod); err == nil {
		test.Errorf("data %+v produced no error!", mod)
	}

	for _, mod = range mods {
		mod = mapMod(writableContent, mod)
		if err = WriteContent(mod); err == nil {
			test.Errorf("data %+v produced no error!", mod)
		}
	}
}

func Test_DeleteContent(test *testing.T) {
	var id string = uuid.New().String()
	var writable map[string]interface{} = mapMod(
		writableContent,
		map[string]interface{}{"id": id},
	)

	var err error
	if err = WriteContent(writable); err != nil {
		test.Fatal(err)
	}

	if err = DeleteContent(id); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if _, exists, err = ReadSingleContent(id); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("deleted content %s still exists", id)
	}
}

func Test_ReadSingleContent(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableContent)
	modified["id"] = uuid.New().String()

	WriteContent(modified)

	var content monketype.Content = monketype.Content{}
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(modified["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("content of id %s does not exist!", modified["id"])
	}

	contentOK(test, modified, content)
}

func Test_ReadSingleContent_ManyTags(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableContent)
	modified["id"] = uuid.New().String()

	var count int = 255
	var tags []string = make([]string, count)

	var index int = 0
	for index != count {
		tags[index] = "some_" + strconv.Itoa(index)
		index++
	}

	modified["tags"] = tags

	WriteContent(modified)

	var content monketype.Content = monketype.Content{}
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(modified["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("content of id %s does not exist!", modified["id"])
	}

	contentOK(test, modified, content)
}

func Test_ReadSingleContent_NotExists(test *testing.T) {
	var id string = uuid.New().String()

	var content monketype.Content
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(id); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("Query for nonexisting id got %+v", content)
	}
}

func Test_ReadManyContent(test *testing.T) {
	var err error
	var many int = 31
	if err = populate(many); err != nil {
		test.Fatal(err)
	}

	var offset, count int = 4, 15

	var content []monketype.Content
	var size int
	if content, size, err = ReadManyContent(offset, count); err != nil {
		test.Fatal(err)
	}

	if size != count {
		test.Errorf("did not get enough posts! have: %d, want: %d", size, count)
	}

	if len(content) != size {
		test.Errorf("content size %d (%+v) does not match size %d!", len(content), content, size)
	}

	var index, suffix int
	var single monketype.Content
	for index, single = range content {
		suffix = size - index + 12
		if single.ID != "many_"+strconv.Itoa(suffix) {
			test.Errorf("ID %s does not have suffix %d!", single.ID, suffix)
		}
	}
}

func Test_ReadManyContent_Fewer(test *testing.T) {
	var err error
	if _, err = database.Query("DROP TABLE IF EXISTS " + CONTENT_TABLE); err != nil {
		test.Fatal(err)
	}

	create()

	var many int = 12
	if err = populate(many); err != nil {
		test.Fatal(err)
	}

	var offset, count int = 4, 15

	var content []monketype.Content
	var size int
	if content, size, err = ReadManyContent(offset, count); err != nil {
		test.Fatal(err)
	}

	if size != many-offset {
		test.Errorf("Got too many or few posts! have: %d, want: %d", size, many-offset)
	}

	if len(content) != size {
		test.Errorf("content size %d (%+v) does not match size %d!", len(content), content, size)
	}
}

func Test_ReadAuthorContent(test *testing.T) {
	var err error
	if err = populate(20); err != nil {
		test.Fatal(err)
	}

	var author string = uuid.New().String()
	var modified map[string]interface{}

	var index, many int = 0, 20
	for index != many {
		modified = mapCopy(writableContent)
		modified["author"] = author
		modified["created"] = time.Now().Unix()
		modified["id"] = uuid.New().String()

		if err = WriteContent(modified); err != nil {
			test.Fatal(err)
		}

		index++
	}

	var offset int = 4

	var content []monketype.Content
	var size int
	if content, size, err = ReadAuthorContent(author, offset, many); err != nil {
		test.Fatal(err)
	}

	if size != many-offset {
		test.Errorf("Got too many or few posts! have: %d, want: %d", size, many-offset)
	}

	var single monketype.Content
	for _, single = range content {
		if single.Author != author {
			test.Errorf("monketype.Content %s author mismatch! have: %s, want: %s", single.ID, single.Author, author)
		}
	}

}
