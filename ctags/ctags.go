package ctags

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

const (
	tagNamePosition = iota
	tagFilePosition
	tagAddressPosition
)

var escapeMap = []struct {
	escaped, raw string
}{
	{`\\`, "\\"},
	{`\n`, "\n"},
	{`\r`, "\r"},
	{`\t`, "\t"},
}

// TagFields is a map of name/value pairs.
type TagFields map[string]string

// TagLine represents a single ctags match.
type TagLine struct {
	// Any identifier, not containing white space.
	TagName string

	// The name of the file where {tagname} is defined, relative to the current
	// directory.
	TagFile string

	// Any Ex command. When executed, it behaves like 'magic' was not set. It may
	// be restricted to a line number or a search pattern (Posix).
	TagAddress string

	// A tagfield is a key, value pair.
	TagFields TagFields
}

// A Reader reads ctags entries.
type Reader struct {
	scanner *bufio.Scanner
	scan    bool
}

// A Writer writes ctags entries.
type Writer struct {
	*bufio.Writer
}

// A TagField has a name, a colon, and a value: “name:value”. The value may be
// empty. It cannot contain a <Tab>.
func parseTagField(data string) (string, string) {
	// The name of the "kind:" field can be omitted. A program reading the tags
	// file can recognize the "kind:" field by the missing ':'.
	if !strings.ContainsRune(data, ':') {
		return "kind", data
	}

	var name, value string

	p := strings.SplitN(data, ":", 2)
	if len(p) == 2 {
		value = p[1]
	}
	name = p[0]

	// The name consists only of alphabetical characters. Upper and lower case
	// are allowed. Case matters (“kind:” and “Kind: are different tagfields).
	isInvalid := func(r rune) bool { return r < 'A' || r > 'z' }
	if strings.IndexFunc(name, isInvalid) != -1 {
		return "", ""
	}

	return name, unescape(value)
}

func parseTagLine(data string) (tl TagLine) {
	properties := strings.Split(data, "\t")
	for i, property := range properties {
		switch i {
		case tagNamePosition:
			tl.TagName = property
		case tagFilePosition:
			tl.TagFile = property
		case tagAddressPosition:
			tl.TagAddress = strings.TrimSuffix(property, `;"`)
		default:
			if tl.TagFields == nil {
				tl.TagFields = make(TagFields)
			}

			key, value := parseTagField(property)
			if key != "" {
				tl.TagFields[key] = value
			}
		}
	}

	return tl
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		scanner: bufio.NewScanner(r),
		scan:    true,
	}
}

func (tl TagLine) valid() bool {
	return tl.TagName != "" && tl.TagFile != "" && tl.TagAddress != ""
}

// Read reads one entry from r. If there is no data left to be read, Read
// returns an empty TagLine and io.EOF.
func (r *Reader) Read() (tl TagLine, err error) {
	for !tl.valid() {
		r.scan = r.scanner.Scan()
		if !r.scan {
			return tl, io.EOF
		}
		tl = parseTagLine(r.scanner.Text())
	}

	return tl, nil
}

// ReadAll reads all the remaining entries from r.
func (r *Reader) ReadAll() []TagLine {
	var lines []TagLine
	var tl TagLine
	var err error

	for err != io.EOF {
		tl, err = r.Read()
		if err == nil {
			lines = append(lines, tl)
		}
	}

	return lines
}

// String implements Stringer.String() from the strings package.
func (tl TagLine) String() string {
	properties := []string{
		tl.TagName,
		tl.TagFile,
		fmt.Sprintf(`%s;"`, tl.TagAddress),
	}

	// Sort keys so that output is deterministic.
	var keys []string
	for key := range tl.TagFields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := escape(tl.TagFields[key])
		properties = append(properties, fmt.Sprintf("%s:%s", key, value))
	}

	return strings.Join(properties, "\t")
}

// Line is the line on which this tag was found. If the line can't be
// determined, -1 is returned.
func (tl TagLine) Line() int {
	if v, ok := tl.TagFields["line"]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	if i, err := strconv.Atoi(tl.TagAddress); err == nil {
		return i
	}

	return -1
}

// Kind is the kind of tag this is. If no "kind" tagfield exists, this returns
// an empty string.
func (tl TagLine) Kind() string {
	if v, ok := tl.TagFields["kind"]; ok {
		return v
	}

	return ""
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		bufio.NewWriter(w),
	}
}

// Write writes a single ctags entry to w. Writes are buffered, so Flush must
// eventually be called to ensure that the record is written to the underlying
// io.Writer.
func (w Writer) Write(tl TagLine) error {
	_, err := fmt.Fprintln(w.Writer, tl.String())

	return err
}

// WriteAll writes multiple ctags entries to w using Write and then calls
// Flush, returning any error from the Flush.
func (w Writer) WriteAll(lines []TagLine) (err error) {
	for _, tl := range lines {
		if err = w.Write(tl); err != nil {
			return err
		}
	}

	return w.Flush()
}

func escape(s string) string {
	// Order matters. For example, given `\	` (backslash followed by tab), we
	// need to escape the slash before the tab. Otherwise we end up with extra
	// slashes ("\\\\t" vs the correct result: "\\\t").
	for _, e := range escapeMap {
		s = strings.ReplaceAll(s, e.raw, e.escaped)
	}

	return s
}

func unescape(s string) string {
	// Unescape in reverse order. For example, given `\\\t`, we need to unescape
	// "\t" before we unescape "\\".
	for i := len(escapeMap) - 1; i >= 0; i-- {
		e := escapeMap[i]
		s = strings.ReplaceAll(s, e.escaped, e.raw)
	}

	return s
}
