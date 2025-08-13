package models

type Article struct {
	Id       uint64 `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorId uint64 `json:"authorId"`
}

type CreateArticle struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorId uint64 `json:"authorId"`
}
