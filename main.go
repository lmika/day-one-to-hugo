package main

import (
	"fmt"
	"github.com/lmika/day-one-to-hugo/models"
	"github.com/lmika/day-one-to-hugo/providers"
	"github.com/lmika/day-one-to-hugo/providers/hugodir"
	"github.com/lmika/day-one-to-hugo/services/converter"
	flag "github.com/spf13/pflag"
	"log"
	"os"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	flagTargetDir := flag.StringP("site", "d", "out", "site dir")
	flagPostsDir := flag.String("posts", "posts", "target post dir relative to site")
	flagFrom := flag.StringP("from", "f", "", "journal entries on and after this date")
	flagTo := flag.StringP("to", "t", "", "journal entries up to, but not including, this date")
	flagDryRun := flag.BoolP("dry-run", "n", false, "dry run")
	flagKeepExif := flag.BoolP("keep-exif", "", false, "keep exif data on jpeg and png images")
	flagHelp := flag.BoolP("help", "h", false, "show usage help")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Day One To Hugo, ver %v (%v), date %v\n\n", version, commit, date)
		fmt.Fprintf(os.Stderr, "Usage\n")
		fmt.Fprintf(os.Stderr, "  %s [OPTIONS] JSON ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	if *flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	var selectOptions converter.SelectOption
	if *flagFrom != "" {
		var err error

		selectOptions.FromDate, err = parseDate(*flagFrom)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *flagTo != "" {
		var err error

		selectOptions.ToDate, err = parseDate(*flagTo)
		if err != nil {
			log.Fatal(err)
		}
	}

	convertOptions := converter.DefaultConvertOptions
	convertOptions.UseFirstHeadingAsTitle = true
	convertOptions.MissingTitlesMode = converter.MissingTitleModeDate
	convertOptions.ConvertStarsToFigureCaptions = false

	if *flagDryRun {
		convertOptions.MissingTitlesMode = converter.MissingTitleModeBlank
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "requires at least JSON file")
		fmt.Fprintln(os.Stderr, "see --help for details")
		os.Exit(1)
	}

	site := models.Site{
		Dir:         *flagTargetDir,
		PostBaseDir: *flagPostsDir,
	}

	for _, j := range flag.Args() {
		de := providers.JournalPackExport(j)
		journalPack, err := de.LoadJournalPack()

		if err != nil {
			log.Fatal(err)
		}

		svc := converter.New(hugodir.New(!*flagKeepExif), convertOptions)

		if *flagDryRun {
			if err := svc.PrintSelectedEntries(os.Stdout, journalPack, selectOptions); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := svc.WriteToHugo(site, journalPack, selectOptions); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func parseDate(date string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).UTC(), nil
}
