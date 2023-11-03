package converter

import (
	"github.com/Southclaws/fault"
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/providers/hugodir"
)

type Service struct {
	hugoDir *hugodir.Provider
}

func New(hugoDir *hugodir.Provider) *Service {
	return &Service{
		hugoDir: hugoDir,
	}
}

func (s *Service) WriteToHugo(site models.Site, journal models.Journal) error {
	for _, entry := range journal.Entries[:1] {
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
