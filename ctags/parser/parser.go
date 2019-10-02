package parser

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

type EntryParser struct {
	parser.Parser
}

func NewEntryParser() EntryParser {
	p := goldmark.DefaultParser()
	p.AddOptions(parser.WithInlineParsers(
		util.Prioritized(gextension.NewStrikethroughParser(), 500),
	))
	p.AddOptions(parser.WithInlineParsers(
		util.Prioritized(extension.NewTagParser(), 0),
	))
	p.AddOptions(parser.WithInlineParsers(
		util.Prioritized(gextension.NewTaskCheckBoxParser(), 0),
	))

	return EntryParser{
		Parser: p,
	}
}

func (p EntryParser) Parse(filename string) (lines []ctags.TagLine, err error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return lines, err
	}

	reader := text.NewReader(source)
	tree := p.Parser.Parse(reader)

	// Reset reader, so we can use it to calculate line numbers
	reader.SetPosition(0, text.Segment{})
	reader.ResetPosition()

	line, pos := reader.Position()
	err = gast.Walk(tree, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		s := gast.WalkStatus(gast.WalkContinue)

		if !entering {
			return s, nil
		}

		if t, ok := n.(*ast.Tag); ok {
			segment := t.Value().Segment
			reader.Advance(segment.Start - pos.Start)
			line, pos = reader.Position()

			tl := parseTag(reader, t)
			tl.TagFile = filename
			lines = append(lines, tl)
		}

		return s, nil
	})

	return lines, err
}

func parseTag(reader text.Reader, t *ast.Tag) ctags.TagLine {
	line, _ := reader.Position()
	return ctags.TagLine{
		TagName:    string(t.Value().Text(reader.Source())),
		TagAddress: fmt.Sprintf("%d;\"", line+1),
		TagFields: []ctags.TagField{
			ctags.TagField{
				Name:  "line",
				Value: fmt.Sprintf("%d", line+1),
			},
			ctags.TagField{
				Name:  "kind",
				Value: strings.ToLower(fmt.Sprintf("%s", t.Kind())),
			},
		},
	}
}
