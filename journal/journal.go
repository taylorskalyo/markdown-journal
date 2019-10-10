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

type journal struct {
	// entries maps tagfiles to tags. These tags represent each tag in the file.
	entries map[string]*tagNode

	// labels maps label names to a slice of tags. For each tag, kind == "label".
	labels map[string][]*tagNode
}

// tagLines attaches the methods of sort.Interface to []ctags.TagLine, sorting
// in increasing order.
type tagLines []ctags.TagLine

type tagNode struct {
	ctags.TagLine
	next *tagNode
	prev *tagNode
}

func newJournal(tags tagLines) journal {
	entries := make(map[string]*tagNode)
	labels := make(map[string][]*tagNode)

	sort.Sort(sort.Reverse(tags))
	for _, tag := range tags {

		// Prepend tag. When done, tags will appear in increasing order.
		n := &tagNode{TagLine: tag}
		head, ok := entries[tag.TagFile]
		if ok {
			head.prev = n
			n.next = head
		}
		entries[tag.TagFile] = n

		if k, ok := tag.TagFields["kind"]; ok && k == "label" {
			l, ok := labels[tag.TagName]
			if !ok {
				labels[tag.TagName] = []*tagNode{}
			}
			labels[tag.TagName] = append(l, n)
		}
	}

	return journal{
		entries: entries,
		labels:  labels,
	}
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
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}(-.*)?\.md`)
	return re.MatchString(file)
}

func (t tagLines) Len() int      { return len(t) }
func (t tagLines) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Sort by tagfile and line number in increasing order. Headings appear first.
func (t tagLines) Less(i, j int) bool {
	var li, lj int
	var v string
	var ok bool

	if t[i].TagFile != t[j].TagFile {
		return t[i].TagFile < t[j].TagFile
	}

	if v, ok = t[i].TagFields["line"]; ok {
		li, _ = strconv.Atoi(v)
	}
	if v, ok = t[j].TagFields["line"]; ok {
		lj, _ = strconv.Atoi(v)
	}

	if li != lj {
		return li < lj
	}

	v, ok = t[i].TagFields["kind"]
	return ok && v == "heading"
}
