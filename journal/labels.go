package journal

import (
	"fmt"
	"io"
	"strings"
)

// WriteLabels generates a list of entries categorized by label and writes the
// result to a writer.
func (j Journal) WriteLabels(w io.Writer, setters ...WriterOption) error {
	opts := &WriterOptions{
		Level: 1,
	}

	for _, setter := range setters {
		setter(opts)
	}

	baseHeadingDelim := strings.Repeat("#", opts.Level)
	for _, label := range j.Labels {
		fmt.Fprintf(w, "\n%s %s\n", baseHeadingDelim, label.Name)

		for _, occur := range label.Occurrences {
			var name string

			location := occur.TagFile
			if line := occur.Line(); line >= 0 {
				location = fmt.Sprintf("%s:%d", location, line)
			}
			if h, ok := occur.TagFields["heading"]; ok {
				name = h
			} else {
				name = location
			}
			fmt.Fprintf(w, "* [%s](%s)\n", name, location)
		}
	}

	return nil
}
