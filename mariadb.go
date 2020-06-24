package monkebase

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	database *sql.DB
	tables   map[string]string = map[string]string{
		"content": `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			file_url VARCHAR(63) NOT NULL,
			author VARCHAR(63) NOT NULL,
			mime VARCHAR(63) NOT NULL,
			like_count BIGINT UNSIGNED NOT NULL,
			dislike_count BIGINT UNSIGNED NOT NULL,
			repub_count BIGINT UNSIGNED NOT NULL,
			view_count BIGINT UNSIGNED NOT NULL,
			comment_count BIGINT UNSIGNED NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			featured BOOLEAN,
			featurable BOOLEAN,
			removed BOOLEAN,
			nsfw BOOLEAN`,
		"users": `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			nick VARCHAR(63) NOT NULL,
			bio VARCHAR(255),
			subscriber_count BIGINT UNSIGNED NOT NULL,
			subscription_count BIGINT UNSIGNED NOT NULL,
			created BIGINT UNSIGNED NOT NULL`,
		"tags": `
			id CHAR(36) NOT NULL,
			tag CHAR(63) NOT NULL,
			CONSTRAINT no_dupe_tags UNIQUE(id, tag)`,
		"auth": `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			hash BINARY(60)`,
	}
)

const (
	USER_TABLE    = "user"
	CONTENT_TABLE = "content"
)

func connect(address string) {
	var err error
	if database, err = sql.Open("mysql", address); err != nil {
		panic(err)
	}

	if err = database.Ping(); err != nil {
		panic(err)
	}

	create()
}

func create() {
	var err error
	var table, structure string
	for table, structure = range tables {
		if _, err = database.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", table, structure)); err != nil {
			panic(err)
		}
	}
}
