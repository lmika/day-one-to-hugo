package hugodir

import (
	"fmt"
	"github.com/Southclaws/fault"
	"github.com/lmika/day-one-to-hugo/models"
	"io"
	"os"
	"path/filepath"
)

type Provider struct {
}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) AddPost(site models.Site, post models.Post) error {
	postFilename := filepath.Join(site.Dir, "content", site.PostBaseDir, "test.md")

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
	if post.Title != "" {
		fmt.Fprintf(w, "title: %v\n", post.Title)
	}
	fmt.Fprintln(w, "---")

	fmt.Fprint(w, post.Content)

	return nil
}
