package journal

import (
	"fmt"
	"io/ioutil"
	"strings"

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

	_, pos := reader.Position()
	err = gast.Walk(tree, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		s := gast.WalkStatus(gast.WalkContinue)

		if !entering {
			return s, nil
		}

		if t, ok := n.(*ast.Label); ok {
			segment := t.Value.Segment
			reader.Advance(segment.Start - pos.Start)
			_, pos = reader.Position()

			tl := parseNode(reader, n)
			tl.TagFile = filename
			lines = append(lines, tl)
		} else if h, ok := n.(*gast.Heading); ok {
			segment := h.Lines().At(0)
			reader.Advance(segment.Start - pos.Start)
			_, pos = reader.Position()

			tl := parseNode(reader, n)
			tl.TagFile = filename
			lines = append(lines, tl)
		}

		return s, nil
	})

	return lines, err
}

func parseNode(reader text.Reader, n gast.Node) ctags.TagLine {
	line, _ := reader.Position()
	return ctags.TagLine{
		TagName:    string(n.Text(reader.Source())),
		TagAddress: fmt.Sprintf("%d", line+1),
		TagFields: ctags.TagFields{
			"line": fmt.Sprintf("%d", line+1),
			"kind": strings.ToLower(n.Kind().String()),
		},
	}
}
