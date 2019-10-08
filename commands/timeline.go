package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylorskalyo/markdown-journal/ctags"
	"github.com/taylorskalyo/markdown-journal/journal"
)

func init() {
	application.AddCommand(timelineCommand)

	tagfileDesc := `read entry info from specified tags file; "-" reads tags from stdin`
	timelineCommand.Flags().StringVarP(&tagfileName, "tagfile", "f", "", tagfileDesc)

	recurseDesc := `recurse into directories`
	timelineCommand.Flags().BoolVarP(&recurse, "recurse", "R", false, recurseDesc)
}

var timelineCommand = &cobra.Command{
	Use:   "timeline [paths]",
	Short: "Display a timeline view",
	Long:  `This command displays a timeline view of journal entries.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var filenames []string
		var tagLines []ctags.TagLine
		var err error

		if tagfileName == "" {

			if len(args) > 0 {
				filenames, err = journal.Files(args, recurse)
			} else {
				filenames, err = journal.Files([]string{"."}, recurse)
			}
			if err != nil {
				log.Fatal(err)
			}

			tagLines, err = generateCtags(filenames)
		} else {

			tagLines, err = readCtags(tagfileName)
		}
		if err != nil {
			log.Fatal(err)
		}

		journal.WriteTimeline(tagLines, os.Stdout)
	},
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
	p := journal.NewEntryParser()
	for _, filename := range filenames {
		lines, err := p.Parse(filename)
		if err != nil {
			return tagLines, err
		}
		tagLines = append(tagLines, lines...)
	}

	return tagLines, err
}
