package extension

import (
	"unicode"

	"github.com/taylorskalyo/markdown-journal/markdown/extension/ast"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

/* Valid labels:
 *   **:label:** is strong or emphaized
 *   a label, :label:, neighbors punctuation (except ':' or '/')
 *
 * Not labels:
 *   https://notalabel.com:3000
 *   Module::notalabel::CONSTANT
 */

type labelParser struct {
}

var defaultLabelParser = &labelParser{}

// NewLabelParser return a new InlineParser that parses
// label expressions.
func NewLabelParser() parser.InlineParser {
	return defaultLabelParser
}

func (s *labelParser) Trigger() []byte {
	return []byte{':'}
}

func isLabelRune(r rune) bool {
	return r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isBoundaryRune(r rune) bool {
	return r != ':' && r != '/' && !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func (s *labelParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	before := block.PrecendingCharacter()
	if !isBoundaryRune(before) {
		return nil
	}

	line, segment := block.PeekLine()
	stop := 1
	for ; stop < len(line) && line[stop] != ':'; stop++ {
		r := util.ToRune(line, stop)
		if !isLabelRune(r) {
			return nil
		}
	}

	if stop >= len(line) {
		return nil
	}

	if stop+1 < len(line) {
		after := util.ToRune(line, stop+1)
		if !isBoundaryRune(after) {
			return nil
		}
	}

	labelSegment := text.NewSegment(segment.Start+1, segment.Start+stop)
	value := gast.NewTextSegment(labelSegment)
	node := ast.NewLabel(value)
	gast.MergeOrAppendTextSegment(node, labelSegment)
	block.Advance(stop + 1)
	return node
}

func (s *labelParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// LabelHTMLRenderer is a renderer.NodeRenderer implementation that
// renders Label nodes.
type LabelHTMLRenderer struct {
	html.Config
}

// NewLabelHTMLRenderer returns a new LabelHTMLRenderer.
func NewLabelHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &LabelHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *LabelHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// nothing to do
}

type label struct {
}

// Label is an extension that allow you to use label expression like ':text:' .
var Label = &label{}

func (e *label) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewLabelParser(), 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewLabelHTMLRenderer(), 0),
	))
}
