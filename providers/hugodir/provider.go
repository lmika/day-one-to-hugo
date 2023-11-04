package hugodir

import (
	"fmt"
	"github.com/Southclaws/fault"
	"github.com/bitfield/script"
	"github.com/lmika/day-one-to-hugo/models"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type Provider struct {
}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) AddPhoto(site models.Site, media models.Media) error {
	//bts, err := os.ReadFile(media.Filename)
	//if err != nil {
	//	return fault.Wrap(err)
	//}

	targetFilename := filepath.Join(site.Dir, "static", "images", filepath.Base(media.Filename))
	if err := p.prepareBaseDir(targetFilename); err != nil {
		return fault.Wrap(err)
	}

	_, err := script.Exec(fmt.Sprintf("magick '%v' -strip '%v'", media.Filename, targetFilename)).Stdout()
	return fault.Wrap(err)
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
	fmt.Fprintln(w, "---")
	fmt.Fprintf(w, "date: %v\n", post.Date.Format(time.RFC3339))
	if post.Title != "" {
		fmt.Fprintf(w, "title: %v\n", post.Title)
	}
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
