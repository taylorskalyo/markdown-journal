package journal

import (
	"fmt"
	"io/ioutil"

	"github.com/taylorskalyo/markdown-journal/ctags"
	"github.com/taylorskalyo/markdown-journal/markdown/extension"
	"github.com/taylorskalyo/markdown-journal/markdown/extension/ast"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	gextension "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// FileParser parses entry file markdown contents into ctags tags.
type FileParser struct {
	parser.Parser
}

// NewFileParser returns a new FileParser.
func NewFileParser() FileParser {
	p := goldmark.DefaultParser()
	p.AddOptions(parser.WithInlineParsers(
		util.Prioritized(gextension.NewStrikethroughParser(), 500),
	))
	p.AddOptions(parser.WithInlineParsers(
		util.Prioritized(extension.NewLabelParser(), 0),
	))
	p.AddOptions(parser.WithInlineParsers(
		util.Prioritized(gextension.NewTaskCheckBoxParser(), 0),
	))

	return FileParser{
		Parser: p,
	}
}

// Parse parses the given entry into ctags tags.
func (p FileParser) Parse(filename string) (lines []ctags.TagLine, err error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return lines, err
	}

	return p.parse(filename, source)
}

func (p FileParser) parse(filename string, source []byte) (lines []ctags.TagLine, err error) {
	reader := text.NewReader(source)
	tree := p.Parser.Parse(reader)

	// Reset reader, so we can use it to calculate line numbers
	reader.SetPosition(0, text.Segment{})
	reader.ResetPosition()

	isTitleFound := false

	line, pos := reader.Position()
	err = gast.Walk(tree, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		s := gast.WalkStatus(gast.WalkContinue)

		if !entering {
			return s, nil
		}

		var segment text.Segment
		var tagFields = ctags.TagFields{}

		switch v := n.(type) {
		case *ast.Label:
			segment = v.Value.Segment
			heading := string(v.Heading.Text(reader.Source()))
			tagFields["heading"] = heading
			tagFields["kind"] = "label"
		case *gast.Heading:
			if isTitleFound {
				return s, nil
			}
			isTitleFound = true
			segment = v.Lines().At(0)
			tagFields["kind"] = "title"
		default:
			return s, nil
		}

		reader.Advance(segment.Start - pos.Start)
		line, pos = reader.Position()
		tagFields["line"] = fmt.Sprintf("%d", line+1)
		tl := ctags.TagLine{
			TagName:    string(n.Text(reader.Source())),
			TagFile:    filename,
			TagAddress: fmt.Sprintf("%d", line+1),
			TagFields:  tagFields,
		}
		lines = append(lines, tl)

		return s, nil
	})

	return lines, err
}
