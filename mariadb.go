package monkebase

import (
	"github.com/jmoiron/sqlx"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	database *sqlx.DB
	tables   map[string]string = map[string]string{
		CONTENT_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			file_url VARCHAR(63) NOT NULL,
			author CHAR(36) NOT NULL,
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
		USER_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			email VARCHAR(63) NOT NULL,
			nick VARCHAR(63) NOT NULL,
			bio VARCHAR(255),
			subscriber_count BIGINT UNSIGNED NOT NULL,
			subscription_count BIGINT UNSIGNED NOT NULL,
			post_count BIGINT UNSIGNED NOT NULL,
			created BIGINT UNSIGNED NOT NULL`,
		TAG_TABLE: `
			id CHAR(36) NOT NULL,
			tag VARCHAR(63) NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			CONSTRAINT no_dupe_tags UNIQUE(id, tag)`,
		AUTH_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			hash BINARY(60)`,
	}
)

const (
	USER_TABLE    = "users"
	CONTENT_TABLE = "content"
	TAG_TABLE     = "tags"
	AUTH_TABLE    = "auth"
)

func connect(address string) {
	var err error
	if database, err = sqlx.Open("mysql", address); err != nil {
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
