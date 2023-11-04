package main

import (
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/providers"
	"github.com/lmika/day-one-to-hugo/providers/hugodir"
	"github.com/lmika/day-one-to-hugo/services/converter"
	"log"
)

func main() {
	de := providers.DirExport("/Users/leonmika/Documents/10-31-2023_9-13-pm")

	site := models.Site{
		Dir:         "/Users/leonmika/tmp/test-journal-export",
		PostBaseDir: "post",
	}

	journalPack, err := de.LoadJournalPack()
	if err != nil {
		log.Fatal(err)
	}

	svc := converter.New(hugodir.New())
	if err := svc.WriteToHugo(site, journalPack); err != nil {
		log.Fatal(err)
	}
}
