// Package ast defines AST nodes that represents extension's elements
package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

// A Label struct represents a label.
type Label struct {
	gast.BaseInline
	Value *gast.Text
}

// Dump implements Node.Dump.
func (n *Label) Dump(source []byte, level int) {
	segment := n.Value.Segment
	m := map[string]string{
		"Value": string(segment.Value(source)),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindLabel is a NodeKind of the Label node.
var KindLabel = gast.NewNodeKind("Label")

// Kind implements Node.Kind.
func (n *Label) Kind() gast.NodeKind {
	return KindLabel
}

// NewLabel returns a new Label node.
func NewLabel(value *gast.Text) *Label {
	return &Label{
		Value: value,
	}
}
