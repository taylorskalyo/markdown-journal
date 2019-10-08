package ctags

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	tagNamePosition = iota
	tagFilePosition
	tagAddressPosition
	tagFieldsPosition
)

// A TagField has a name, a colon, and a value: “name:value”.
//
// - The name consists only of alphabetical characters. Upper and lower case
// are allowed. Lower case is recommended. Case matters (“kind:” and “Kind: are
// different tagfields).
//
// - The value may be empty. It cannot contain a <Tab>.
type TagField struct {
	Name  string
	Value string
}

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

	// A list of TagFields.
	TagFields []TagField
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

func parseTagField(data string) (tf TagField) {
	fieldPair := strings.SplitN(data, ":", 2)
	if len(fieldPair) > 1 {
		tf.Name = fieldPair[0]
		tf.Value = fieldPair[1]
	}

	return tf
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
			tl.TagAddress = property
		case tagFieldsPosition:
			for _, field := range strings.Fields(property) {
				tf := parseTagField(field)
				if tf != (TagField{}) {
					tl.TagFields = append(tl.TagFields, tf)
				}
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

// Read reads one entry from r. If there is no data left to be read, Read
// returns an empty TagLine and io.EOF.
func (r *Reader) Read() (TagLine, error) {
	r.scan = r.scanner.Scan()
	if !r.scan {
		return TagLine{}, io.EOF
	}
	tl := parseTagLine(r.scanner.Text())

	return tl, nil
}

// ReadAll reads all the remaining entries from r.
func (r *Reader) ReadAll() []TagLine {
	var lines []TagLine
	for tl, err := r.Read(); err == nil; {
		lines = append(lines, tl)
	}

	return lines
}

// String implements Stringer.String() from the strings package.
func (tf TagField) String() string {
	return fmt.Sprintf("%s:%s", tf.Name, tf.Value)
}

// String implements Stringer.String() from the strings package.
func (tl TagLine) String() string {
	fields := []string{}
	for _, field := range tl.TagFields {
		fields = append(fields, field.String())
	}
	properties := []string{
		tl.TagName,
		tl.TagFile,
		tl.TagAddress,
		strings.Join(fields, " "),
	}
	return strings.Join(properties, "\t")
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
