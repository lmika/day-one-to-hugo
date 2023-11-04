package main

import (
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/providers"
	"github.com/lmika/day-one-to-hugo/providers/hugodir"
	"github.com/lmika/day-one-to-hugo/services/converter"
	"log"
)

func main() {
	de := providers.DirExport("/Users/leonmika/Documents/11-04-2023_9-44-pm")

	site := models.Site{
		Dir:         "/Users/leonmika/Developer/Websites/untraveller-web-v2",
		PostBaseDir: "posts",
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
