package providers

import (
	"encoding/json"
	"github.com/lmika/day-one-to-hugo/models"
	"os"
	"path/filepath"
)

type DirExport string

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
