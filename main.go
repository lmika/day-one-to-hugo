package main

import (
	"flag"
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/providers"
	"github.com/lmika/day-one-to-hugo/providers/hugodir"
	"github.com/lmika/day-one-to-hugo/services/converter"
	"log"
)

func main() {
	flagTargetDir := flag.String("d", "out", "target site dir")
	flagPostsDir := flag.String("posts", "posts", "target post dir relative to site")
	flagTitle := flag.Bool("title", false, "force posts to have a title")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("require journal JSON file")
	}

	de := providers.JournalPackExport(flag.Arg(0))

	site := models.Site{
		Dir:         *flagTargetDir,
		PostBaseDir: *flagPostsDir,
	}

	journalPack, err := de.LoadJournalPack()
	if err != nil {
		log.Fatal(err)
	}

	svc := converter.New(hugodir.New(), *flagTitle)
	if err := svc.WriteToHugo(site, journalPack); err != nil {
		log.Fatal(err)
	}
}
