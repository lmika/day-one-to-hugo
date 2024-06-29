# day-one-to-hugo

A tool for converting posts from a Day One JSON export to posts for a Hugo site.

## Installation

### MacOS

The recommended way to install this on MacOS is using Homebrew. To install, run the following commands:

```
brew tap lmika/day-one-to-hugo
brew install day-one-to-hugo
```

There is also a regular `.tar.gz` file available.

### Linux

For Linux, both RPM and DEB packages are available from the release artefacts, plus regular `.tar.gz` files.
After installing, the command will be installed at `/usr/local/bin/day-one-to-hugo`.

### Windows

ZIPs with the Windows binaries are available from the release artefacts.

### Go

If you have Go, you can install day-one-to-hugo using this command:

```
go install github.com/lmika/day-one-to-hugo@latest
```

## Basic Usage

1. Make a JSON export from Day One.
2. Unzip it and locate the `.json` files of the journals you'd like to export.
3. Create a [new Hugo site](https://gohugo.io/getting-started/quick-start/).
4. Run `day-one-to-hugo`, passing the `.json` file of the journal you want to export, and setting `-d` to the directory of the Hugo site: 
   ```
   $ day-one-to-hugo -d hugo-site Journal.json
   ```
   This will write out all posts to `<hugo-site>/posts` and all the media — images and videos — to `<hugo-site>/static`.

## Full Usage

```
$ day-one-hugo OPTIONS JSON ...
```
Options are:

- `--site, -d <dir>` — base directory of the Hugo site to export to. Default is `out`
- `--posts <dir>` — directory within the Hugo `content` directory to write the posts to
  (this is usually a property of the layout). Default is `posts`.
- `--from, -f <date>` — export posts that occur on or after this date. Format is `YYYY-MM-DD`
- `--to, -t <date>` — export posts that occur before, but not including, this date. Format is `YYYY-MM-DD`
- `--dry-run, -n` — only print the date, and first line, of the posts that will be exported. Does not export anything.
- `--keep-exif` — keep EXIF metadata on images exported to Hugo. The default is that images are
  exported to Hugo with EXIF tags stripped. Does not affect videos (EXIF tags are not modified).

JSON is one or more Journal `.json` files to export.