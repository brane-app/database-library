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
			email VARCHAR(63) UNIQUE NOT NULL,
			nick VARCHAR(63) UNIQUE NOT NULL,
			bio VARCHAR(255) NOT NULL,
			subscriber_count BIGINT UNSIGNED NOT NULL,
			subscription_count BIGINT UNSIGNED NOT NULL,
			post_count BIGINT UNSIGNED NOT NULL,
			created BIGINT UNSIGNED NOT NULL`,
		AUTH_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			hash BINARY(60) NOT NULL`,
		TOKEN_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			token BINARY(24) UNIQUE,
			created BIGINT UNSIGNED NOT NULL`,
		SECRET_TABLE: `
			id CHAR(36) UNIQUE PRIMARY KEY NOT NULL,
			secret BINARY(128) UNIQUE`,
		TAG_TABLE: `
			id CHAR(36) NOT NULL,
			tag VARCHAR(63) NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			CONSTRAINT no_dupe_tags UNIQUE(id, tag)`,
		SUBSCRIPTION_TABLE: `
			subscriber CHAR(36) NOT NULL,
			subscription CHAR(36) NOT NULL,
			created BIGINT UNSIGNED NOT NULL,
			CONSTRAINT no_dupe_subscriptions UNIQUE(subscriber, subscription)`,
	}
)

const (
	USER_TABLE         = "users"
	CONTENT_TABLE      = "content"
	AUTH_TABLE         = "auth"
	TOKEN_TABLE        = "token"
	SECRET_TABLE       = "secret"
	TAG_TABLE          = "tags"
	SUBSCRIPTION_TABLE = "subs"
)

func Connect(address string) {
	var err error
	if database, err = sqlx.Open("mysql", address); err != nil {
		panic(err)
	}

	if err = database.Ping(); err != nil {
		panic(err)
	}

	database.SetMaxOpenConns(150)
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
