package commands

import (
	"log"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/taylorskalyo/markdown-journal/ctags"
	"github.com/taylorskalyo/markdown-journal/journal"
)

var (
	nosort           bool
	ctagsTagfileName string
)

func init() {
	application.AddCommand(ctagsCommand)

	tagfileDesc := `write tags to specified file; "-" writes tags to stdout`
	ctagsCommand.Flags().StringVarP(&ctagsTagfileName, "tagfile", "f", "tags", tagfileDesc)

	nosortDesc := `do not sort tags by tagname`
	ctagsCommand.Flags().BoolVar(&nosort, "no-sort", false, nosortDesc)

	recurseDesc := `recurse into subdirectories`
	ctagsCommand.Flags().BoolVarP(&recurse, "recurse", "R", false, recurseDesc)
}

var ctagsCommand = &cobra.Command{
	Use:   "ctags [paths]",
	Short: "Generate ctags",
	Long:  `This command generates a ctags compatible tags file.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var journalFiles []string
		var tagLines []ctags.TagLine
		var err error
		var tagfile *os.File

		if len(args) > 0 {
			journalFiles, err = journal.Files(args, recurse)
		} else {
			journalFiles, err = journal.Files([]string{"."}, recurse)
		}
		if err != nil {
			log.Fatal(err)
		}

		tagLines, err = generateCtags(journalFiles)
		if err != nil {
			log.Fatal(err)
		}

		if ctagsTagfileName == "-" {
			tagfile = os.Stdout
		} else {
			tagfile, err = os.Create(ctagsTagfileName)
			if err != nil {
				log.Fatal(err)
			}
		}

		if !nosort {
			sort.Slice(tagLines, func(i, j int) bool {
				return tagLines[i].TagName < tagLines[j].TagName
			})
		}

		w := ctags.NewWriter(tagfile)
		w.WriteAll(tagLines)
	},
}
