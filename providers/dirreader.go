package providers

import (
	"encoding/json"
	"github.com/Southclaws/fault"
	"github.com/lmika/day-one-to-hugo/models"
	"io/fs"
	"os"
	"path/filepath"
)

type JournalPackExport string

func (de JournalPackExport) LoadJournalPack() (models.JournalPack, error) {
	journal, err := de.LoadJournal()
	if err != nil {
		return models.JournalPack{}, fault.Wrap(err)
	}

	photos, err := de.LoadPhotos()
	if err != nil {
		return models.JournalPack{}, fault.Wrap(err)
	}

	videos, err := de.LoadVideos()
	if err != nil {
		return models.JournalPack{}, fault.Wrap(err)
	}

	return models.JournalPack{
		Journal: journal,
		Media:   append(append([]models.Media{}, photos...), videos...),
	}, nil
}

func (de JournalPackExport) LoadJournal() (j models.Journal, err error) {
	bts, err := os.ReadFile(string(de))
	if err != nil {
		return models.Journal{}, err
	}

	if err := json.Unmarshal(bts, &j); err != nil {
		return models.Journal{}, err
	}

	return j, nil
}

func (de JournalPackExport) LoadPhotos() ([]models.Media, error) {
	photoDir := filepath.Join(filepath.Dir(string(de)), "photos")
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

func (de JournalPackExport) LoadVideos() ([]models.Media, error) {
	videoDir := filepath.Join(filepath.Dir(string(de)), "videos")
	media := make([]models.Media, 0)

	if err := filepath.Walk(videoDir, func(path string, info fs.FileInfo, err error) error {
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
