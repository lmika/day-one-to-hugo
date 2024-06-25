package models

import (
	"fmt"
	"time"
)

type JournalPack struct {
	Journal Journal
	Media   []Media
}

type Journal struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Date     time.Time `json:"creationDate"`
	Text     string    `json:"text"`
	Photos   []Moment  `json:"photos"`
	Videos   []Moment  `json:"videos"`
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

	Video bool `json:"-"`
}

func (m Moment) CanStripExif() bool {
	// Not all the image types, but only those that can have their EXIF data stripped
	return m.Type == "jpeg" || m.Type == "jpg" || m.Type == "png"
}

func (m Moment) BaseName() string {
	return fmt.Sprintf("%v.%v", m.MD5, m.Type)
}
