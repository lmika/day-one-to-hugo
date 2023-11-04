package models

import "time"

type JournalPack struct {
	Journal Journal
	Photos  []Media
}

type Journal struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Date   time.Time `json:"creationDate"`
	Text   string    `json:"text"`
	Photos []Moment  `json:"photos"`
}

type Moment struct {
	ID     string `json:"identifier"`
	Type   string `json:"type"`
	MD5    string `json:"md5"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
