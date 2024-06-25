package models

import (
	"time"
)

type Site struct {
	Dir         string
	PostBaseDir string
}

type HugoContent struct {
	Posts []Post
	Media []Media
}

type Post struct {
	Date       time.Time
	Title      string
	BlankTitle string
	Content    string
	Location   Location
	Weather    Weather
	Moments    []Moment
}

type Posts []Post

func (p Posts) InterestedMoments() map[string]Moment {
	moments := make(map[string]Moment)
	for _, post := range p {
		for _, m := range post.Moments {
			moments[m.BaseName()] = m
		}
	}
	return moments
}

type Media struct {
	Filename string
}
