package journal

import (
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/taylorskalyo/markdown-journal/ctags"
)

// WriteTimeline generates a timeline view of entries, using ctags, and writes
// the result to a writer.
func WriteTimeline(tagLines []ctags.TagLine, w io.Writer) error {
	sort.Slice(tagLines, func(i, j int) bool {
		return tagLines[i].TagFile < tagLines[j].TagFile
	})

	entries := groupByFile(tagLines)

	// Sort entries by filename so that they are displayed in chronological
	// order. Because the journal filename format is YYYY-MM-DD, we can sort
	// lexicographically to achieve chronological order.
	//
	// Since we can't sort a map, create a sorted slice of the map's keys.
	var entryFiles []string
	for entryFile := range entries {
		entryFiles = append(entryFiles, entryFile)

		// Also sort each entry's tags so we can later find the first heading.
		sortByLineNumber(entries[entryFile])
	}
	sort.Strings(entryFiles)

	var year int
	var month time.Month
	for _, entryFile := range entryFiles {
		var entryTitle string

		for _, entryTag := range entries[entryFile] {
			if kind, ok := entryTag.TagFields["kind"]; ok && kind == "heading" {
				entryTitle = entryTag.TagName
			}
		}
		entryName := strings.TrimSuffix(path.Base(entryFile), path.Ext(entryFile))
		entryDate, err := time.Parse(dateFormat, entryName)
		if err != nil {
			return err
		}

		// Write new year when it changes
		if year != entryDate.Year() {
			year = entryDate.Year()
			fmt.Fprintf(w, "\n# %s\n", entryDate.Format(yearFormat))
		}

		// Write new month when it changes
		if month != entryDate.Month() {
			month = entryDate.Month()
			fmt.Fprintf(w, "\n## %s\n", entryDate.Format(monthFormat))
		}

		// Write day, and link to the entry
		fmt.Fprintf(w, "* [%s](%s)", entryDate.Format(dayFormat), entryFile)
		if entryTitle != "" {
			fmt.Fprintf(w, " - %s\n", entryTitle)
		} else {
			fmt.Fprintf(w, "\n")
		}
	}

	return nil
}
