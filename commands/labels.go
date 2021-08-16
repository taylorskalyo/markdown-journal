package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylorskalyo/markdown-journal/journal"
)

func init() {
	application.AddCommand(labelsCommand)

	tagfileDesc := `read entry info from specified tags file; "-" reads tags from stdin`
	labelsCommand.Flags().StringVarP(&tagfileName, "tagfile", "f", "", tagfileDesc)

	recurseDesc := `recurse into directories`
	labelsCommand.Flags().BoolVarP(&recurse, "recurse", "R", false, recurseDesc)

	levelDesc := `base heading level`
	labelsCommand.Flags().IntVarP(&level, "level", "H", 1, levelDesc)

	filterDesc := `filter`
	labelsCommand.Flags().StringArrayVarP(&filters, "filter", "Q", []string{}, filterDesc)
}

var labelsCommand = &cobra.Command{
	Use:   "labels [paths]",
	Short: "Display a entries by label",
	Long:  `This command displays a list of journal entries categorized by label.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		j, err := newJournal(args)
		if err != nil {
			log.Fatal(err)
		}
		j.WriteLabels(
			os.Stdout,
			journal.HeadingLevel(level),
			journal.LabelFilters(filters),
		)
	},
}
