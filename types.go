package monkebase

type Content struct {
	ID           string   `json:"id"				db:"id"`
	FileURL      string   `json:"file_url"			db:"file_url"`
	Author       string   `json:"author"			db:"author"`
	Mime         string   `json:"mime"				db:"mime"`
	Tags         []string `json:"tags"				db:"tags"`
	LikeCount    int      `json:"like_count"		db:"like_count"`
	DislikeCount int      `json:"dislike_count"		db:"dislike_count"`
	RepubCount   int      `json:"repub_count"		db:"repub_count"`
	ViewCount    int      `json:"view_count"		db:"view_count"`
	CommentCount int      `json:"comment_count"		db:"comment_count"`
	Created      int64    `json:"created"			db:"created"`
	Featured     bool     `json:"featured"			db:"featured"`
	Featurable   bool     `json:"featurable"		db:"featurable"`
	Removed      bool     `json:"removed"			db:"removed"`
	NSFW         bool     `json:"nsfw"				db:"nsfw"`
}
