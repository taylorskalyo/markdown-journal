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

/* Valid tags:
 *   **:tag:** is strong or emphaized
 *   a tag, :tag:, neighbors punctuation (except ':' or '/')
 *
 * Not tags:
 *   https://notatag.com:3000
 *   Module::notatag::CONSTANT
 */

type tagParser struct {
}

var defaultTagParser = &tagParser{}

// NewTagParser return a new InlineParser that parses
// tag expressions.
func NewTagParser() parser.InlineParser {
	return defaultTagParser
}

func (s *tagParser) Trigger() []byte {
	return []byte{':'}
}

func isTagRune(r rune) bool {
	return r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isBoundaryRune(r rune) bool {
	return r != ':' && r != '/' && !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func (s *tagParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	before := block.PrecendingCharacter()
	if !isBoundaryRune(before) {
		return nil
	}

	line, segment := block.PeekLine()
	stop := 1
	for ; stop < len(line) && line[stop] != ':'; stop++ {
		r := util.ToRune(line, stop)
		if !isTagRune(r) {
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

	tagSegment := text.NewSegment(segment.Start+1, segment.Start+stop)
	value := gast.NewTextSegment(tagSegment)
	node := ast.NewTag(value)
	gast.MergeOrAppendTextSegment(node, tagSegment)
	block.Advance(stop + 1)
	return node
}

func (s *tagParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// TagHTMLRenderer is a renderer.NodeRenderer implementation that
// renders Tag nodes.
type TagHTMLRenderer struct {
	html.Config
}

// NewTagHTMLRenderer returns a new TagHTMLRenderer.
func NewTagHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &TagHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *TagHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// nothing to do
}

type tag struct {
}

// Tag is an extension that allow you to use tag expression like '~~text~~' .
var Tag = &tag{}

func (e *tag) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTagParser(), 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewTagHTMLRenderer(), 0),
	))
}
