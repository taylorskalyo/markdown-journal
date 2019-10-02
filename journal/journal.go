package journal

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/taylorskalyo/markdown-journal/ctags"
)

const dateFormat = "2006-01-02"

// Files finds journal entry files. It walks each given path checking for ones
// that look like journal entries. It returns a list of the entries it finds.
// If recurse is true, Files will recurse into subdirectories.
func Files(paths []string, recurse bool) (entries []string, err error) {
	for _, path := range paths {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && !recurse {
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

// WriteTimeline generates a timeline view of entries, using ctags, and writes
// the result to a writer.
func WriteTimeline(tagLines []ctags.TagLine, w io.Writer) error {
	sort.Slice(tagLines, func(i, j int) bool {
		return tagLines[i].TagFile < tagLines[j].TagFile
	})

	entries := make(map[string][]ctags.TagLine)
	for _, tagLine := range tagLines {
		entryTags, ok := entries[tagLine.TagFile]
		if !ok {
			entryTags = []ctags.TagLine{}
		}
		entries[tagLine.TagFile] = append(entryTags, tagLine)
	}

	var entryFiles []string
	for k := range entries {
		entryFiles = append(entryFiles, k)
		sort.Slice(entries[k], func(i, j int) bool {
			var li, lj int
			var err error
			for _, tf := range entries[k][i].TagFields {
				if tf.Name == "line" {
					li, err = strconv.Atoi(tf.Value)
					if err != nil {
						continue
					}
					break
				}
			}
			for _, tf := range entries[k][j].TagFields {
				if tf.Name == "line" {
					lj, err = strconv.Atoi(tf.Value)
					if err != nil {
						continue
					}
					break
				}
			}
			return li > lj
		})
	}
	sort.Strings(entryFiles)

	var year int
	var month time.Month
	for _, entryFile := range entryFiles {
		var entryTitle string

		for _, entryTag := range entries[entryFile] {
			for _, field := range entryTag.TagFields {
				if field.Name == "kind" && field.Value == "heading" {
					entryTitle = entryTag.TagName
				}
			}
		}
		entryName := strings.TrimSuffix(path.Base(entryFile), path.Ext(entryFile))
		entryDate, err := time.Parse(dateFormat, entryName)
		if err != nil {
			return err
		}

		if year != entryDate.Year() {
			year = entryDate.Year()
			fmt.Fprintf(w, "\n# %d\n", year)
		}

		if month != entryDate.Month() {
			month = entryDate.Month()
			fmt.Fprintf(w, "\n## %s\n", month)
		}

		if entryTitle != "" {
			fmt.Fprintf(w, "* [%d %s](%s) - %s\n", entryDate.Day(), entryDate.Weekday(), entryFile, entryTitle)
		} else {
			fmt.Fprintf(w, "* [%d %s](%s)\n", entryDate.Day(), entryDate.Weekday(), entryFile)
		}
	}

	return nil
}
