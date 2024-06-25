package converter

import (
	"context"
	"fmt"
	"github.com/Southclaws/fault"
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/providers/hugodir"
	"github.com/lmika/gopkgs/fp/slices"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

type Service struct {
	hugoDir *hugodir.Provider

	convertOptions ConvertOptions
}

func New(hugoDir *hugodir.Provider, convertOptions ConvertOptions) *Service {
	return &Service{
		hugoDir:        hugoDir,
		convertOptions: convertOptions,
	}
}

func (s *Service) WriteToHugo(site models.Site, journalPack models.JournalPack, selectOptions SelectOption) error {
	posts, err := s.convertToPosts(journalPack.Journal, selectOptions)
	if err != nil {
		return fault.Wrap(err)
	}

	log.Printf("writing out %v posts", len(posts))

	if err := s.writeJournalToHugo(site, posts); err != nil {
		return fault.Wrap(err)
	}

	interestedMoments := posts.InterestedMoments()

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(runtime.NumCPU())

	log.Printf("writing out %v media objects", len(interestedMoments))

	for _, photo := range journalPack.Media {
		p := photo
		moment, ok := interestedMoments[filepath.Base(p.Filename)]
		if !ok {
			continue
		}

		g.Go(func() error {
			if err := s.hugoDir.AddPhoto(site, p, moment); err != nil {
				return fault.Wrap(err)
			}
			return nil
		})
	}

	return g.Wait()
}

func (s *Service) PrintSelectedEntries(w io.Writer, pack models.JournalPack, options SelectOption) error {
	posts, err := s.convertToPosts(pack.Journal, options)
	if err != nil {
		return fault.Wrap(err)
	}

	for _, p := range posts {
		var preview string

		if p.Title != "" {
			preview = p.Title
		} else {
			preview = strings.Split(p.Content, "\n")[0]
		}
		if len(preview) > 60 {
			preview = preview[:60]
		}

		if _, err := fmt.Fprintf(w, "%s\t%s\n", p.Date.Local().Format("2006-01-02 15:04:05"), preview); err != nil {
			return fault.Wrap(err)
		}
	}
	return nil
}

func (s *Service) convertToPosts(journal models.Journal, selectOptions SelectOption) (models.Posts, error) {
	selectedEntries := slices.Filter(journal.Entries, selectOptions.EntryMatch)

	return slices.MapWithError(selectedEntries, func(e models.Entry) (models.Post, error) {
		return s.ConvertToPost(e)
	})
}

func (s *Service) writeJournalToHugo(site models.Site, posts models.Posts) error {
	for _, p := range posts {
		if err := s.hugoDir.AddPost(site, p); err != nil {
			return fault.Wrap(err)
		}
	}

	return nil
}
