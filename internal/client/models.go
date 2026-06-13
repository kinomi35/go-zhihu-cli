package client

type Me struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URLToken string `json:"url_token"`
}

type FeedResponse struct {
	Data   []FeedItem `json:"data"`
	Paging Paging     `json:"paging"`
}

type Paging struct {
	IsEnd bool   `json:"is_end"`
	Next  string `json:"next"`
}

type FeedItem struct {
	ID     string     `json:"id"`
	Type   string     `json:"type"`
	Target FeedTarget `json:"target"`
}

type FeedTarget struct {
	ID       any       `json:"id"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Excerpt  string    `json:"excerpt"`
	Content  string    `json:"content"`
	Question *Question `json:"question"`
	Author   *Author   `json:"author"`
}

type Question struct {
	ID    any    `json:"id"`
	Title string `json:"title"`
}

type Author struct {
	Name     string `json:"name"`
	URLToken string `json:"url_token"`
}

type Answer struct {
	ID           any       `json:"id"`
	Content      string    `json:"content"`
	Excerpt      string    `json:"excerpt"`
	VoteupCount  int       `json:"voteup_count"`
	CommentCount int       `json:"comment_count"`
	Author       *Author   `json:"author"`
	Question     *Question `json:"question"`
}

type CommentsResponse struct {
	Data []Comment `json:"data"`
}

type Comment struct {
	ID        any     `json:"id"`
	Content   string  `json:"content"`
	VoteCount int     `json:"vote_count"`
	Author    *Author `json:"author"`
}
