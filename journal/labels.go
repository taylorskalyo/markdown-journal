package journal

import (
	"fmt"
	"io"
	"regexp"
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
		shouldFilter, err := filterLabel(opts.LabelFilters, label.Name)
		if err != nil {
			return err
		} else if shouldFilter {
			continue
		}

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

func filterLabel(filters []string, label string) (bool, error) {
	var shouldFilter bool

	for _, filter := range filters {
		r, err := regexp.Compile(filter)
		if err != nil {
			return false, err
		}

		if !r.Match([]byte(label)) {
			shouldFilter = true
			break
		}
	}

	return shouldFilter, nil
}
