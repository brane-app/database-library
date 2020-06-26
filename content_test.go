package monkebase

import (
	"testing"
)

var (
	writableContent map[string]interface{} = map[string]interface{}{
		"id":            "96910fdf-916b-4664-a38c-be42d0f2c0ce",
		"file_url":      "https://gastrodon.io/file/foobar",
		"author":        "8f12f748-58c5-4f08-9f68-ef34ded3a3cf",
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

func Test_WriteContent(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{},
		map[string]interface{}{
			"tags": []string{},
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

		if err = WriteContent(testWritable{Data: copy}); err != nil {
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

		if err = WriteContent(testWritable{Data: copy}); err == nil {
			test.Errorf("data %+v produced no error!", copy)
		}
	}

	copy = mapCopy(writableContent)
	delete(copy, "file_url")

	if err = WriteContent(testWritable{Data: copy}); err == nil {
		test.Errorf("data %+v produced no error!", copy)
	}
}
