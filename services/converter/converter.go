package converter

import (
	"bytes"
	"fmt"
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/gopkgs/fp/slices"
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

	if err := ast.Walk(n, s.imageURLWalker(entry)); err != nil {
		return models.Post{}, err
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

func (s *Service) imageURLWalker(entry models.Entry) ast.Walker {
	const dayOnePrefix = "dayone-moment://"

	return func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Kind() != ast.KindImage {
			return ast.WalkContinue, nil
		}

		imgNode := n.(*ast.Image)
		dest := string(imgNode.Destination)

		if !strings.HasPrefix(dest, dayOnePrefix) {
			return ast.WalkContinue, nil
		}

		momentID := strings.TrimPrefix(dest, dayOnePrefix)

		photo, found := slices.FindWhere(entry.Photos, func(t models.Moment) bool {
			return t.ID == momentID
		})
		if !found {
			return ast.WalkContinue, nil
		}

		imgNode.Destination = []byte(fmt.Sprintf("/images/%v.%v", photo.MD5, photo.Type))

		return ast.WalkContinue, nil
	}
}
