package journal

import (
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/taylorskalyo/markdown-journal/ctags"
)

const (
	dateFormat  = "2006-01-02"
	yearFormat  = "2006"
	monthFormat = "January"
	dayFormat   = "02 Mon"
)

// Label is a keyword that appears in a journal entry.
type Label struct {
	Name        string
	Occurrences []LabelTag
}

// LabelTag is an occurrence of a Label within a journal entry.
type LabelTag struct {
	*tagNode
}

// Journal is a collection of entries and labels.
type Journal struct {
	Entries []Entry
	Labels  []Label
}

// TagLines attaches the methods of sort.Interface to []ctags.TagLine, sorting
// in increasing order.
type TagLines []ctags.TagLine

type tagNode struct {
	ctags.TagLine
	next *tagNode
	prev *tagNode
}

// NewJournal returns a new Journal.
func NewJournal(tags TagLines) (j Journal) {
	var err error
	var e Entry
	var l Label
	var occurrences []LabelTag

	sort.Sort(sort.Reverse(tags))
	for _, tag := range tags {
		if tag.TagFile != e.File {
			if e.File != "" {
				j.Entries = append(j.Entries, e)
			}

			e, err = NewEntry(tag.TagFile)
			if err != nil {
				continue
			}
		}

		// Prepend tag. When done, an Entry's tags will appear in increasing order
		// by line number, headings first.
		n := &tagNode{TagLine: tag}
		head := e.FirstTag
		if head != nil {
			head.prev = n
			n.next = head
		}
		e.FirstTag = n

		if tag.Kind() == "label" {
			occurrences = append(occurrences, LabelTag{n})
		}
	}

	sort.Slice(occurrences, func(i, j int) bool {
		return occurrences[i].TagName < occurrences[j].TagName
	})

	for _, o := range occurrences {
		if o.TagName != l.Name {
			l = Label{Name: o.TagName}
			j.Labels = append(j.Labels, l)
		}
		l.Occurrences = append(l.Occurrences, o)
	}

	return j
}

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
	return reEntryFile.MatchString(path.Base(file))
}

func (t TagLines) Len() int      { return len(t) }
func (t TagLines) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Sort by tagfile then line number in increasing order. Headings appear first.
func (t TagLines) Less(i, j int) bool {
	if t[i].TagFile != t[j].TagFile {
		return t[i].TagFile < t[j].TagFile
	}

	if li, lj := t[i].Line(), t[j].Line(); li != lj {
		return li < lj
	}

	return t[i].Kind() == "heading"
}
