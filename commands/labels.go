package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylorskalyo/markdown-journal/ctags"
	"github.com/taylorskalyo/markdown-journal/journal"
)

func init() {
	application.AddCommand(labelsCommand)

	tagfileDesc := `read entry info from specified tags file; "-" reads tags from stdin`
	labelsCommand.Flags().StringVarP(&tagfileName, "tagfile", "f", "", tagfileDesc)

	recurseDesc := `recurse into directories`
	labelsCommand.Flags().BoolVarP(&recurse, "recurse", "R", false, recurseDesc)
}

var labelsCommand = &cobra.Command{
	Use:   "labels [paths]",
	Short: "Display a entries by label",
	Long:  `This command displays a list of journal entries categorized by label.`,
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
		j.WriteLabels(os.Stdout)
	},
}
