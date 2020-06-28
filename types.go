package monkebase

type Content struct {
	ID           string   `db:"id"`
	FileURL      string   `db:"file_url"`
	Author       string   `db:"author"`
	Mime         string   `db:"mime"`
	Tags         []string `db:"tags"`
	LikeCount    int      `db:"like_count"`
	DislikeCount int      `db:"dislike_count"`
	RepubCount   int      `db:"repub_count"`
	ViewCount    int      `db:"view_count"`
	CommentCount int      `db:"comment_count"`
	Created      int64    `db:"created"`
	Featured     bool     `db:"featured"`
	Featurable   bool     `db:"featurable"`
	Removed      bool     `db:"removed"`
	NSFW         bool     `db:"nsfw"`
}
