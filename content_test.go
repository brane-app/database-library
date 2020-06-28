package monkebase

import (
	"github.com/google/uuid"

	"sort"
	"testing"
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
		"created":       1593108723,
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

	var key string
	var value interface{}
	var err error

	var mod, copy map[string]interface{}
	for _, mod = range mods {
		copy = mapCopy(writableContent)

		for key, value = range mod {
			copy[key] = value
		}

		if err = WriteContent(copy); err != nil {
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

	var key string
	var value interface{}
	var err error

	var mod, copy map[string]interface{}
	for _, mod = range mods {
		copy = mapCopy(writableContent)

		for key, value = range mod {
			copy[key] = value
		}

		if err = WriteContent(copy); err == nil {
			test.Errorf("data %+v produced no error!", copy)
		}
	}

	copy = mapCopy(writableContent)
	delete(copy, "file_url")

	if err = WriteContent(copy); err == nil {
		test.Errorf("data %+v produced no error!", copy)
	}
}

func Test_ReadSingleContent(test *testing.T) {
	var copy map[string]interface{} = mapCopy(writableContent)
	copy["id"] = uuid.New().String()

	WriteContent(copy)

	var content Content = Content{}
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(copy["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("content of id %s does not exist!", copy["id"])
	}

	contentOK(test, copy, content)
}

func Test_ReadSingleContent_ManyTags(test *testing.T) {
	var copy map[string]interface{} = mapCopy(writableContent)
	copy["id"] = uuid.New().String()

	var count int = 63
	var tags []string = make([]string, count)

	var index int = 0
	for index != count {
		tags[index] = "some_" + string(index)
		index++
	}

	copy["tags"] = tags

	WriteContent(copy)

	var content Content = Content{}
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(copy["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("content of id %s does not exist!", copy["id"])
	}

	contentOK(test, copy, content)
}
