package models

import "time"

type Site struct {
	Dir         string
	PostBaseDir string
}

type HugoContent struct {
	Posts []Post
	Media []Media
}

type Post struct {
	Date    time.Time
	Title   string
	Content string
}

type Media struct {
	Filename string
}
