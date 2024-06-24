package converter

import (
	"context"
	"github.com/Southclaws/fault"
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/providers/hugodir"
	"github.com/lmika/gopkgs/fp/slices"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"runtime"
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

func (s *Service) WriteToHugo(site models.Site, journalPack models.JournalPack) error {
	posts, err := s.convertToPosts(journalPack.Journal)
	if err != nil {
		return fault.Wrap(err)
	}

	if err := s.writeJournalToHugo(site, posts); err != nil {
		return fault.Wrap(err)
	}

	interestedMoments := posts.InterestedMoments()

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(runtime.NumCPU())

	for _, photo := range journalPack.Photos {
		p := photo
		if _, ok := interestedMoments[filepath.Base(p.Filename)]; !ok {
			continue
		}

		g.Go(func() error {
			if err := s.hugoDir.AddPhoto(site, p); err != nil {
				return fault.Wrap(err)
			}
			return nil
		})
	}

	return g.Wait()
}

func (s *Service) convertToPosts(journal models.Journal) (models.Posts, error) {
	return slices.MapWithError(journal.Entries, func(e models.Entry) (models.Post, error) {
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
