package hugodir

import (
	"fmt"
	"github.com/Southclaws/fault"
	"github.com/lmika/day-one-to-hugo/models"
	exifremove "github.com/neurosnap/go-exif-remove"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type Provider struct {
	StripExifData bool
}

func New(stripExifData bool) *Provider {
	return &Provider{
		stripExifData,
	}
}

func (p *Provider) AddPhoto(site models.Site, media models.Media, moment models.Moment) error {
	subDir := "images"
	if moment.Video {
		subDir = "videos"
	}

	targetFilename := filepath.Join(site.Dir, "static", subDir, filepath.Base(media.Filename))
	if err := p.prepareBaseDir(targetFilename); err != nil {
		return fault.Wrap(err)
	}

	if !p.StripExifData {
		return p.quickCopy(targetFilename, media)
	} else if !moment.CanStripExif() {
		log.Printf("warn: cannot remove exif from %s. Writing out unmodified", targetFilename)
		return p.quickCopy(targetFilename, media)
	}

	imgIn, err := os.ReadFile(media.Filename)
	if err != nil {
		return fault.Wrap(err)
	}

	imgOut, err := exifremove.Remove(imgIn)
	if err != nil {
		return fault.Wrap(err)
	} else if len(imgOut) == 0 {
		log.Printf("warn: cannot remove exif from %s. Writing out unmodified", targetFilename)
		imgOut = imgIn
	}

	if err := os.WriteFile(targetFilename, imgOut, 0644); err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func (p *Provider) AddPost(site models.Site, post models.Post) error {
	postFilename := filepath.Join(site.Dir, "content", site.PostBaseDir, p.postFilename(post))

	if err := p.prepareBaseDir(postFilename); err != nil {
		return fault.Wrap(err)
	}

	f, err := os.Create(postFilename)
	if err != nil {
		return fault.Wrap(err)
	}
	defer f.Close()

	return p.generatePostBody(f, post)
}

func (p *Provider) prepareBaseDir(filename string) error {
	return os.MkdirAll(filepath.Dir(filename), 0755)
}

func (p *Provider) generatePostBody(w io.Writer, post models.Post) error {
	frontMatter := map[string]any{
		"date": post.Date.Format(time.RFC3339),
		"location": map[string]any{
			"placeName": post.Location.PlaceName,
			"locality":  post.Location.Locality,
			"country":   post.Location.Country,
			"lat":       post.Location.Lat,
			"long":      post.Location.Long,
		},
		"weather": map[string]any{
			"code": post.Weather.Code,
			"temp": post.Weather.Temp,
		},
	}
	if post.Title != "" {
		frontMatter["title"] = post.Title
	} else if post.BlankTitle != "" {
		frontMatter["title"] = post.BlankTitle
	}

	fmStr, err := yaml.Marshal(frontMatter)
	if err != nil {
		return fault.Wrap(err)
	}

	fmt.Fprintln(w, "---")
	fmt.Fprint(w, string(fmStr))
	fmt.Fprintln(w, "---")

	fmt.Fprint(w, post.Content)

	return nil
}

func (p *Provider) postFilename(post models.Post) string {
	var wordComponent string
	if post.Title != "" {
		wordComponent = scanNWords(post.Title, 3)
	} else {
		wordComponent = scanNWords(post.Content, 3)
	}

	return fmt.Sprintf("%d/%02d/%d/%v.md", post.Date.Year(), int(post.Date.Month()), post.Date.Day(), wordComponent)
}

func (p *Provider) quickCopy(targetFilename string, media models.Media) error {
	w, err := os.Create(targetFilename)
	if err != nil {
		return fault.Wrap(err)
	}
	defer w.Close()

	r, err := os.Open(media.Filename)
	if err != nil {
		return fault.Wrap(err)
	}
	defer r.Close()

	if _, err := io.Copy(w, r); err != nil {
		return fault.Wrap(err)
	}

	return nil
}

func scanNWords(s string, words int) string {
	var sb strings.Builder

	wordCount := 0
	inWord := false
	for _, r := range s {
		switch {
		case r == '\'':
			// ignore
		case unicode.IsDigit(r) || unicode.IsLetter(r):
			sb.WriteRune(unicode.ToLower(r))
			inWord = true
		case inWord:
			inWord = false
			wordCount += 1
			if wordCount >= words {
				return sb.String()
			}
			sb.WriteRune('-')
		default:
			// ignore
		}
	}

	return sb.String()
}
