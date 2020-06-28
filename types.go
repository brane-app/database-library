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

type User struct {
	ID                string `db:"id"`
	Email             string `db:"email"`
	Nick              string `db:"nick"`
	Bio               string `db:"bio"`
	SubscriberCount   int    `db:"subscriber_count"`
	SubscriptionCount int    `db:"subscription_count"`
	PostCount         int    `db:"post_count"`
	Created           int64  `db:"created"`
}
