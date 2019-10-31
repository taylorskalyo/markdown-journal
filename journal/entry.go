package journal

import (
	"errors"
	"path"
	"regexp"
	"strings"
	"time"
)

var errNotEntry = errors.New("not a journal entry")
var reEntryFile = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})(-.*)?\.md`)

// Entry is a file in the journal and its associated tags.
type Entry struct {
	Time time.Time
	File string

	// Linked list of tags.
	FirstTag *tagNode

	name string
}

// NewEntry returns a new Entry.
func NewEntry(file string) (e Entry, err error) {
	matches := reEntryFile.FindStringSubmatch(path.Base(file))
	if len(matches) == 0 {
		err = errNotEntry
		return e, err
	}

	e.Time, err = time.Parse(dateFormat, matches[1])
	if err != nil {
		return e, err
	}

	e.name = matches[2]

	e.File = file

	return e, nil
}

// Title returns a title for the entry. It uses the first heading if one
// exists. If no heading is found, it uses the portion of the entry's filename
// after the date as the title. In the latter case, dashes (`-`) and
// underscores (`_`) are converted to spaces, and the first letter of each word
// is capitalized.
func (e Entry) Title() string {
	for n := e.FirstTag; n != nil; n = n.next {
		if n.Kind() == "title" {
			return n.TagName
		}
	}

	if title := e.name; title != "" {
		title = strings.ReplaceAll(title, "-", " ")
		title = strings.ReplaceAll(title, "_", " ")
		title = strings.Title(title)
		title = strings.TrimSpace(title)
		return title
	}

	return ""
}
