package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylorskalyo/markdown-journal/ctags"
	"github.com/taylorskalyo/markdown-journal/journal"
)

var (
	tagfileName string
	recurse     bool
	level       int
)

var application = &cobra.Command{
	Use:   "markdown-journal",
	Short: "markdown-journal helps you manage a markdown journal",
	Long:  `A markdown journaling system`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute root command.
func Execute() {
	if err := application.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readCtags(tagfileName string) (tagLines []ctags.TagLine, err error) {
	var tagfile *os.File

	if tagfileName == "-" {
		tagfile = os.Stdin
	} else {
		tagfile, err = os.Open(tagfileName)
		if err != nil {
			return tagLines, err
		}
	}

	r := ctags.NewReader(tagfile)
	tagLines = r.ReadAll()

	return tagLines, err
}

func generateCtags(filenames []string) (tagLines []ctags.TagLine, err error) {
	p := journal.NewFileParser()
	for _, filename := range filenames {
		lines, err := p.Parse(filename)
		if err != nil {
			return tagLines, err
		}
		tagLines = append(tagLines, lines...)
	}

	return tagLines, err
}

func newJournal(filenames []string) (j journal.Journal, err error) {
	var journalFiles []string
	var tagLines []ctags.TagLine

	if len(filenames) > 0 {
		journalFiles, err = journal.Files(filenames, recurse)
	} else {
		journalFiles, err = journal.Files([]string{"."}, recurse)
	}
	if err != nil {
		return j, err
	}

	if tagfileName == "" {
		tagLines, err = generateCtags(journalFiles)
	} else {
		tagLines, err = readCtags(tagfileName)
	}
	if err != nil {
		return j, err
	}

	// Some files may not have any labels or headings and therefore no ctags
	// entries. Ensure every file has at least one ctags entry.
	fileTagLines := make([]ctags.TagLine, len(journalFiles))
	for i, file := range journalFiles {
		fileTagLines[i] = ctags.TagLine{
			TagFile:   file,
			TagFields: ctags.TagFields{"line": "0"},
		}
	}
	tagLines = append(tagLines, fileTagLines...)

	return journal.NewJournal(tagLines), nil
}
