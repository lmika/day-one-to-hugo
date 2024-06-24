package converter

// ConvertOptions are the available conversion options.
// TODO: these should be processors
type ConvertOptions struct {
	// UseFirstHeadingAsTitle controls whether to use the first heading of the journal entry as the title.
	UseFirstHeadingAsTitle bool

	// SetMissingTitlesToDate controls whether to use the journal entries date as the title if a title was not found
	// in the entry.
	SetMissingTitlesToDate bool

	// ConvertStarsToFigureCaptions searches for any lines with stars that appear directly after an image as the
	// image caption.
	ConvertStarsToFigureCaptions bool
}

var DefaultConvertOptions = ConvertOptions{
	UseFirstHeadingAsTitle:       true,
	SetMissingTitlesToDate:       true,
	ConvertStarsToFigureCaptions: true,
}
