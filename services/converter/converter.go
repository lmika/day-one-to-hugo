package converter

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/gopkgs/fp/slices"
	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	maxMissingTitleModeLength = 45
)

func (s *Service) ConvertToPost(entry models.Entry) (post models.Post, err error) {
	entryText := s.replaceBackslashes(entry.Text)

	textSrc := []byte(entryText)
	mdParser := goldmark.DefaultParser()

	n := mdParser.Parse(text.NewReader(textSrc))

	// TODO: make this a processor
	if s.convertOptions.UseFirstHeadingAsTitle {
		docHead := n.FirstChild()
		if docHead.Kind() == ast.KindHeading {
			post.Title = string(docHead.Text(textSrc))
			n.RemoveChild(n, docHead)
		}
	}

	// TODO: make this a processor
	if post.Title == "" {
		switch s.convertOptions.MissingTitlesMode {
		case MissingTitleModeDate:
			post.BlankTitle = entry.Date.Format("Jan 02, 2006")
		case MissingTitleModeFirstLine:
			post.BlankTitle = strings.Split(strings.TrimSpace(entry.Text), "\n")[0]
			if len(post.Title) > maxMissingTitleModeLength {
				post.BlankTitle = post.BlankTitle[:maxMissingTitleModeLength]
			}
		default:
			// leave blank
		}
	}

	foundMoments := make([]models.Moment, 0)
	if err := ast.Walk(n, s.imageURLWalker(entry, &foundMoments)); err != nil {
		return models.Post{}, err
	}

	outBfr := bytes.Buffer{}
	mdRenderer := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	if err := mdRenderer.Render(&outBfr, textSrc, n); err != nil {
		return models.Post{}, err
	}
	post.Content = outBfr.String()
	post.Content = s.convertToVideo(post.Content)

	// TODO: make this a processor
	if s.convertOptions.ConvertStarsToFigureCaptions {
		post.Content = s.figureMaker(post.Content)
	}

	post.Date = entry.Date
	post.Location = entry.Location
	post.Weather = entry.Weather
	post.Moments = foundMoments

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

func (s *Service) figureMaker(md string) string {
	var (
		imgLine     = regexp.MustCompile(`^!\[]\((.*)\)`)
		captionLine = regexp.MustCompile(`^[*]([^*]+)[*]`)
	)

	type pendingSeen struct {
		imgURL string
		mode   int
	}

	var bts bytes.Buffer
	seen := pendingSeen{}
	scnr := bufio.NewScanner(strings.NewReader(md))

	for scnr.Scan() {
		text := scnr.Text()
		switch {
		case seen.mode == 2 && captionLine.MatchString(text):
			caption := captionLine.FindStringSubmatch(text)[1]
			bts.WriteString(fmt.Sprintf(`<figure><img src="%v"><figcaption>%v</figcaption></figure>`,
				seen.imgURL, caption))
			bts.WriteString("\n\n")
			seen = pendingSeen{}
		case seen.mode == 1 && strings.TrimSpace(text) == "":
			seen.mode = 2
		case imgLine.MatchString(text):
			if seen.imgURL != "" {
				bts.WriteString(fmt.Sprintf(`<img src="%v">`, seen.imgURL))
				bts.WriteString("\n\n")
			}

			seen.imgURL = imgLine.FindStringSubmatch(text)[1]
			seen.mode = 1
		default:
			if seen.imgURL != "" {
				bts.WriteString(fmt.Sprintf(`<img src="%v">`, seen.imgURL))
				bts.WriteString("\n\n")
			}
			seen = pendingSeen{}
			bts.WriteString(text)
			bts.WriteString("\n")
		}
	}

	if seen.imgURL != "" {
		bts.WriteString(fmt.Sprintf(`<img src="%v">`, seen.imgURL))
		bts.WriteString("\n\n")
	}

	return bts.String()
}

func (s *Service) convertToVideo(md string) string {
	var (
		imgLine = regexp.MustCompile(`^!\[]\((.*)\)`)
	)

	type pendingSeen struct {
		imgURL string
		mode   int
	}

	var bts bytes.Buffer
	scnr := bufio.NewScanner(strings.NewReader(md))

	for scnr.Scan() {
		text := scnr.Text()
		switch {
		case imgLine.MatchString(text):
			imgURL := imgLine.FindStringSubmatch(text)[1]
			if strings.HasPrefix(imgURL, "/videos/") {
				urlPart, queryPart, found := strings.Cut(imgURL, "?")
				if found {
					query, _ := url.ParseQuery(queryPart)
					bts.WriteString(fmt.Sprintf(`<video src="%v" controls width="%v" height="%v"></video>`, urlPart,
						query.Get("w"), query.Get("h")))
					bts.WriteString("\n\n")
				} else {
					bts.WriteString(fmt.Sprintf(`<video src="%v" controls></video>`, imgURL))
					bts.WriteString("\n\n")
				}
			} else {
				bts.WriteString(text)
				bts.WriteString("\n")
			}
		default:
			bts.WriteString(text)
			bts.WriteString("\n")
		}
	}

	return bts.String()
}

func (s *Service) imageURLWalker(entry models.Entry, foundMoments *[]models.Moment) ast.Walker {
	const dayOnePrefix = "dayone-moment:/"

	return func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Kind() != ast.KindImage {
			return ast.WalkContinue, nil
		}

		imgNode := n.(*ast.Image)
		dest := string(imgNode.Destination)

		if !strings.HasPrefix(dest, dayOnePrefix) {
			return ast.WalkContinue, nil
		}

		momentID := filepath.Base(strings.TrimPrefix(dest, dayOnePrefix))

		// Search for image
		photo, found := slices.FindWhere(entry.Photos, func(t models.Moment) bool { return t.ID == momentID })
		if found {
			imgNode.Destination = []byte(fmt.Sprintf("/images/%v.%v", photo.MD5, photo.Type))
			*foundMoments = append(*foundMoments, photo)

			return ast.WalkContinue, nil
		}

		// Search for video
		video, found := slices.FindWhere(entry.Videos, func(t models.Moment) bool { return t.ID == momentID })
		if found {
			// This is a bit of a hack, but I'd like to pass the width and height through to the video processor.
			imgNode.Destination = []byte(fmt.Sprintf("/videos/%v.%v?w=%v&h=%v", video.MD5, video.Type, video.Width, video.Height))
			video.Video = true
			*foundMoments = append(*foundMoments, video)

			return ast.WalkContinue, nil
		}

		return ast.WalkContinue, nil
	}
}
