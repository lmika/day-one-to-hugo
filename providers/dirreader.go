package providers

import (
	"encoding/json"
	"github.com/Southclaws/fault"
	"github.com/lmika/day-one-to-hugo/models"
	"io/fs"
	"os"
	"path/filepath"
)

type DirExport string

func (de DirExport) LoadJournalPack() (models.JournalPack, error) {
	journal, err := de.LoadJournal()
	if err != nil {
		return models.JournalPack{}, fault.Wrap(err)
	}

	photos, err := de.LoadPhotos()
	if err != nil {
		return models.JournalPack{}, fault.Wrap(err)
	}

	return models.JournalPack{
		Journal: journal,
		Photos:  photos,
	}, nil
}

func (de DirExport) LoadJournal() (j models.Journal, err error) {
	bts, err := os.ReadFile(filepath.Join(string(de), "Journal.json"))
	if err != nil {
		return models.Journal{}, err
	}

	if err := json.Unmarshal(bts, &j); err != nil {
		return models.Journal{}, err
	}

	return j, nil
}

func (de DirExport) LoadPhotos() ([]models.Media, error) {
	photoDir := filepath.Join(string(de), "photos")
	media := make([]models.Media, 0)

	if err := filepath.Walk(photoDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		media = append(media, models.Media{Filename: path})
		return nil
	}); err != nil {
		return nil, err
	}

	return media, nil
}
