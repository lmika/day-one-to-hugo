package converter

import (
	"github.com/lmika/day-one-to-hugo/models"
	"time"
)

type MissingTitleMode int

const (
	MissingTitleModeBlank MissingTitleMode = iota
	MissingTitleModeDate
	MissingTitleModeFirstLine
)

type SelectOption struct {
	FromDate time.Time
	ToDate   time.Time
}

func (o SelectOption) EntryMatch(entry models.Entry) bool {
	if !o.FromDate.IsZero() {
		if entry.Date.Before(o.FromDate) {
			return false
		}
	}

	if !o.ToDate.IsZero() {
		if entry.Date.After(o.ToDate) || entry.Date.Equal(o.ToDate) {
			return false
		}
	}

	return true
}

// ConvertOptions are the available conversion options.
// TODO: these should be processors
type ConvertOptions struct {
	// UseFirstHeadingAsTitle controls whether to use the first heading of the journal entry as the title.
	UseFirstHeadingAsTitle bool

	// MissingTitleMode controls what to use for missing titles.
	MissingTitlesMode MissingTitleMode

	// ConvertStarsToFigureCaptions searches for any lines with stars that appear directly after an image as the
	// image caption.
	ConvertStarsToFigureCaptions bool
}

var DefaultConvertOptions = ConvertOptions{
	UseFirstHeadingAsTitle:       true,
	MissingTitlesMode:            MissingTitleModeDate,
	ConvertStarsToFigureCaptions: true,
}
