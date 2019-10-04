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

	recurseDesc := `recurse into directories`
	timelineCommand.Flags().BoolVarP(&recurse, "recurse", "R", false, recurseDesc)
}

var timelineCommand = &cobra.Command{
	Use:   "timeline [files]",
	Short: "Display a timeline view",
	Long:  `This command displays a timeline view of journal entries.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var filenames []string
		var tagLines []ctags.TagLine
		var err error

		if len(args) > 0 {
			filenames, err = journal.Files(args, recurse)
		} else {
			filenames, err = journal.Files([]string{"."}, recurse)
		}
		if err != nil {
			log.Fatal(err)
		}

		p := journal.NewEntryParser()
		for _, filename := range filenames {
			lines, err := p.Parse(filename)
			if err != nil {
				log.Fatal(err)
			}
			tagLines = append(tagLines, lines...)
		}

		journal.WriteTimeline(tagLines, os.Stdout)
	},
}
