package converter

import (
	"bytes"
	"github.com/lmika/day-one-to-hugo/models"
	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"strings"
)

func (s *Service) ConvertToPost(entry models.Entry) (post models.Post, err error) {
	entryText := s.replaceBackslashes(entry.Text)

	textSrc := []byte(entryText)
	mdParser := goldmark.DefaultParser()

	n := mdParser.Parse(text.NewReader(textSrc))

	docHead := n.FirstChild()
	if docHead.Kind() == ast.KindHeading {
		post.Title = string(docHead.Text(textSrc))
		n.RemoveChild(n, docHead)
	}

	outBfr := bytes.Buffer{}
	mdRenderer := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	if err := mdRenderer.Render(&outBfr, textSrc, n); err != nil {
		return models.Post{}, err
	}
	post.Content = outBfr.String()
	post.Date = entry.Date

	return post, nil
}

func (s *Service) replaceBackslashes(str string) string {
	var outStr strings.Builder

	inBackSlash := false
	for _, r := range str {
		switch {
		case inBackSlash:
			outStr.WriteRune(r)
			inBackSlash = false
		case r == '\\':
			inBackSlash = true
		default:
			outStr.WriteRune(r)
		}
	}

	return outStr.String()
}
