package journal

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// WriteTimeline generates a timeline view of entries and writes the result to
// a writer.
func (j Journal) WriteTimeline(w io.Writer, setters ...WriterOption) error {
	var year int
	var month time.Month

	opts := &WriterOptions{
		Level: 1,
	}

	for _, setter := range setters {
		setter(opts)
	}

	baseHeadingDelim := strings.Repeat("#", opts.Level)
	for _, entry := range j.Entries {
		// Write new year when it changes
		if year != entry.Time.Year() {
			year = entry.Time.Year()
			fmt.Fprintf(w, "\n%s %s\n", baseHeadingDelim, entry.Time.Format(yearFormat))
		}

		// Write new month when it changes
		if month != entry.Time.Month() {
			month = entry.Time.Month()
			fmt.Fprintf(w, "\n%s# %s\n", baseHeadingDelim, entry.Time.Format(monthFormat))
		}

		// Write day, and link to the entry
		fmt.Fprintf(w, "* [%s](%s)", entry.Time.Format(dayFormat), entry.File)
		if title := entry.Title(); title != "" {
			fmt.Fprintf(w, " - %s\n", title)
		} else {
			fmt.Fprintf(w, "\n")
		}
	}

	return nil
}
