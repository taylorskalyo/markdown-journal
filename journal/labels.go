package journal

import (
	"fmt"
	"io"
)

// WriteLabels generates a list of entries categorized by label and writes the
// result to a writer.
func (j Journal) WriteLabels(w io.Writer) error {
	for _, label := range j.Labels {
		fmt.Fprintf(w, "\n# %s\n", label.Name)

		for _, occur := range label.Occurrences {
			var name string

			location := occur.TagFile
			if line := occur.Line(); line >= 0 {
				location = fmt.Sprintf("%s:%d", location, line)
			}
			section := occur.section()
			if section != nil {
				name = section.TagName
			} else {
				name = location
			}
			fmt.Fprintf(w, "* [%s](%s)\n", name, location)
		}
	}

	return nil
}
