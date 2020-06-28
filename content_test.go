package monkebase

import (
	"github.com/google/uuid"

	"sort"
	"strconv"
	"testing"
	"time"
)

var (
	writableContent map[string]interface{} = map[string]interface{}{
		"id":            uuid.New().String(),
		"file_url":      "https://gastrodon.io/file/foobar",
		"author":        uuid.New().String(),
		"mime":          "png",
		"tags":          []string{"some", "tags"},
		"like_count":    10,
		"dislike_count": 1,
		"repub_count":   3,
		"view_count":    400,
		"comment_count": 4,
		"created":       time.Now().Unix(),
		"featured":      false,
		"featurable":    true,
		"removed":       false,
		"nsfw":          false,
	}
)

func contentOK(test *testing.T, data map[string]interface{}, have Content) {
	if data["id"].(string) != have.ID {
		test.Errorf("Content ID mismatch! have: %s, want: %s", have.ID, data["id"])
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

func mapMod(source map[string]interface{}, mods ...map[string]interface{}) (modified map[string]interface{}) {
	modified = mapCopy(source)

	var key string
	var value interface{}

	var mod map[string]interface{}
	for _, mod = range mods {
		for key, value = range mod {
			modified[key] = value
		}
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

func Test_ReadSingleContent(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableContent)
	modified["id"] = uuid.New().String()

	WriteContent(modified)

	var content Content = Content{}
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

	var content Content = Content{}
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
