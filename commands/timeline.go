package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylorskalyo/markdown-journal/journal"
)

func init() {
	application.AddCommand(timelineCommand)

	tagfileDesc := `read entry info from specified tags file; "-" reads tags from stdin`
	timelineCommand.Flags().StringVarP(&tagfileName, "tagfile", "f", "", tagfileDesc)

	recurseDesc := `recurse into directories`
	timelineCommand.Flags().BoolVarP(&recurse, "recurse", "R", false, recurseDesc)

	levelDesc := `base heading level`
	timelineCommand.Flags().IntVarP(&level, "level", "H", 1, levelDesc)
}

var timelineCommand = &cobra.Command{
	Use:   "timeline [paths]",
	Short: "Display a timeline view",
	Long:  `This command displays a timeline view of journal entries.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		j, err := newJournal(args)
		if err != nil {
			log.Fatal(err)
		}
		j.WriteTimeline(os.Stdout, journal.HeadingLevel(level))
	},
}
