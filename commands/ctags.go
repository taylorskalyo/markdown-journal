package commands

import (
	"log"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/taylorskalyo/markdown-journal/ctags"
	"github.com/taylorskalyo/markdown-journal/journal"
)

var tagfileName string
var recurse bool
var nosort bool

func init() {
	application.AddCommand(ctagsCommand)

	tagfileDesc := `write tags to specified file; "-" writes tags to stdout`
	ctagsCommand.Flags().StringVarP(&tagfileName, "tagfile", "f", "tags", tagfileDesc)

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
		var filenames []string
		var tagLines []ctags.TagLine
		var err error
		var tagfile *os.File

		if len(args) > 0 {
			filenames, err = journal.Files(args, recurse)
		} else {
			filenames, err = journal.Files([]string{"."}, recurse)
		}
		if err != nil {
			log.Fatal(err)
		}

		for _, filename := range filenames {
			p := journal.NewFileParser()
			lines, err := p.Parse(filename)
			if err != nil {
				log.Fatal(err)
			}
			tagLines = append(tagLines, lines...)
		}

		if tagfileName == "-" {
			tagfile = os.Stdout
		} else {
			tagfile, err = os.Create(tagfileName)
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
