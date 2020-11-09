package monkebase

const (
	CONTENT_FIELDS = `
id,
file_url,
author,
mime,
like_count,
dislike_count,
repub_count,
view_count,
comment_count,
created,
featured,
featurable,
removed,
nsfw`
	USER_FIELDS = `
id,
email,
nick,
bio,
subscriber_count,
subscription_count,
post_count,
created,
moderator,
admin`
	BAN_FIELDS = `
id,
banner,
banned,
reason,
created,
expires,
forever`
	REPORT_FIELDS = `
id,
reporter,
reported,
type,
reason,
created,
resolved,
resolution`

	READ_CONTENT_ID   = "SELECT " + CONTENT_FIELDS + " FROM " + CONTENT_TABLE + " WHERE id=? LIMIT 1"
	DELETE_CONTENT_ID = "DELETE FROM " + CONTENT_TABLE + " WHERE id=? LIMIT 1"

	// TODO deprecate
	DEPRECATED_READ_MANY_CONTENT           = "SELECT " + CONTENT_FIELDS + " FROM " + CONTENT_TABLE + " ORDER BY created DESC LIMIT ?, ?"
	DEPRECATED_READ_MANY_CONTENT_OF_AUTHOR = "SELECT " + CONTENT_FIELDS + " FROM " + CONTENT_TABLE + " WHERE author=? ORDER BY created DESC LIMIT ?, ?"

	READ_TAGS_OF_ID         = "SELECT tag FROM " + TAG_TABLE + " WHERE id=?"
	READ_TAGS_OF_MANY_ID    = "SELECT id, tag FROM " + TAG_TABLE + " WHERE id IN "
	WRITE_TAGS_OF_MANY_ID   = "REPLACE INTO " + TAG_TABLE + " (id, tag, created) VALUES "
	DELETE_STALE_TAGS_OF_ID = "DELETE FROM " + TAG_TABLE + " WHERE id=? AND tag NOT IN "
	DELETE_TAGS_OF_ID       = "DELETE FROM " + TAG_TABLE + " WHERE id=?"
)
