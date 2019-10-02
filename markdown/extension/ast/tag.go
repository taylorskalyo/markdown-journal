// Package ast defines AST nodes that represents extension's elements
package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

// A Tag struct represents a tag.
type Tag struct {
	gast.BaseInline
	value *gast.Text
}

// Dump implements Node.Dump.
func (n *Tag) Dump(source []byte, level int) {
	segment := n.value.Segment
	m := map[string]string{
		"Value": string(segment.Value(source)),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindTag is a NodeKind of the Tag node.
var KindTag = gast.NewNodeKind("Tag")

// Kind implements Node.Kind.
func (n *Tag) Kind() gast.NodeKind {
	return KindTag
}

// NewTag returns a new Tag node.
func NewTag(value *gast.Text) *Tag {
	return &Tag{
		value: value,
	}
}

func (n *Tag) Value() *gast.Text {
	return n.value
}
