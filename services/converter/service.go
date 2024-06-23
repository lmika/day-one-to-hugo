package converter

import (
	"context"
	"github.com/Southclaws/fault"
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/providers/hugodir"
	"golang.org/x/sync/errgroup"
	"runtime"
)

type Service struct {
	hugoDir  *hugodir.Provider
	incTitle bool
}

func New(hugoDir *hugodir.Provider, incTitle bool) *Service {
	return &Service{
		hugoDir:  hugoDir,
		incTitle: incTitle,
	}
}
func (s *Service) WriteToHugo(site models.Site, journalPack models.JournalPack) error {
	if err := s.writeJournalToHugo(site, journalPack.Journal); err != nil {
		return fault.Wrap(err)
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(runtime.NumCPU())

	for _, photo := range journalPack.Photos {
		p := photo
		g.Go(func() error {
			if err := s.hugoDir.AddPhoto(site, p); err != nil {
				return fault.Wrap(err)
			}
			return nil
		})
	}

	return g.Wait()
}

func (s *Service) writeJournalToHugo(site models.Site, journal models.Journal) error {
	for _, entry := range journal.Entries {
		p, err := s.ConvertToPost(entry)
		if err != nil {
			return fault.Wrap(err)
		}

		if err := s.hugoDir.AddPost(site, p); err != nil {
			return fault.Wrap(err)
		}
	}

	return nil
}
