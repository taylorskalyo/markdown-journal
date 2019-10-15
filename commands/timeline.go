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

		j := journal.NewJournal(tagLines)
		j.WriteTimeline(os.Stdout)
	},
}
