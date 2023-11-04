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
	Date     time.Time `json:"creationDate"`
	Text     string    `json:"text"`
	Photos   []Moment  `json:"photos"`
	Location Location  `json:"location"`
	Weather  Weather   `json:"weather"`
}

type Location struct {
	PlaceName string  `json:"placeName"`
	Locality  string  `json:"localityName"`
	Country   string  `json:"country"`
	Lat       float64 `json:"latitude"`
	Long      float64 `json:"longitude"`
}

type Weather struct {
	Code string  `json:"weatherCode"`
	Temp float64 `json:"temperatureCelsius"`
}

type Moment struct {
	ID     string `json:"identifier"`
	Type   string `json:"type"`
	MD5    string `json:"md5"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
