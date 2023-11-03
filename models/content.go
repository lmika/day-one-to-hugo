package models

type Site struct {
	Dir         string
	PostBaseDir string
}

type HugoContent struct {
	Posts []Post
	Media []Media
}

type Post struct {
	Title   string
	Content string
}

type Media struct {
	Filename string
}
