package journal

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/taylorskalyo/markdown-journal/ctags"
)

const (
	dateFormat  = "2006-01-02"
	yearFormat  = "2006"
	monthFormat = "January"
	dayFormat   = "02 Mon"
)

// Files finds journal entry files. It walks each given path checking for ones
// that look like journal entries. It returns a list of the entries it finds.
// If recurse is true, Files will recurse into subdirectories.
func Files(paths []string, recurse bool) (entries []string, err error) {
	for _, pathArg := range paths {
		err = filepath.Walk(pathArg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Only visit a directory if it was supplied as an argument or recurse
			// option is true.
			if info.IsDir() && path != pathArg && !recurse {
				return filepath.SkipDir
			}

			if !info.IsDir() && isJournalFile(info.Name()) {
				entries = append(entries, path)
			}
			return nil
		})

		if err != nil {
			return entries, err
		}
	}

	return entries, err
}

func isJournalFile(file string) bool {
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}(-.*)?\.md`)
	return re.MatchString(file)
}

func groupByFile(tagLines []ctags.TagLine) map[string][]ctags.TagLine {
	entries := make(map[string][]ctags.TagLine)
	for _, tagLine := range tagLines {
		entryTags, ok := entries[tagLine.TagFile]
		if !ok {
			entryTags = []ctags.TagLine{}
		}
		entries[tagLine.TagFile] = append(entryTags, tagLine)
	}

	for k := range entries {
		sortByLineNumber(entries[k])
	}

	return entries
}

func sortByLineNumber(tagLines []ctags.TagLine) {
	sort.Slice(tagLines, func(i, j int) bool {
		var li, lj int
		var v string
		var ok bool

		if v, ok = tagLines[i].TagFields["line"]; ok {
			li, _ = strconv.Atoi(v)
		}
		if v, ok = tagLines[j].TagFields["line"]; ok {
			lj, _ = strconv.Atoi(v)
		}

		return li > lj
	})
}
