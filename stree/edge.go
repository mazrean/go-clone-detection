package stree

import "fmt"

type edge struct {
	tree  *STree
	label *label
	node  *node
}

func newEdge(tree *STree, label *label, node *node) *edge {
	return &edge{
		tree:  tree,
		label: label,
		node:  node,
	}
}

func (e *edge) splitEdge(i int64, newSuffixLink *node) (*node, *edge, error) {
	l, err := newLabel(i, e.label.end)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating new label(under): %v", err)
	}

	newEdge := newEdge(e.tree, l, e.node)
	newNode := newInternalNode(e.tree, newSuffixLink, []*edge{newEdge})

	e.label, err = newLabel(e.label.start, i)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating new label(upper): %v", err)
	}

	e.node = newNode

	return newNode, newEdge, nil
}

func (e *edge) getLength() int64 {
	return e.label.end - e.label.start
}

func (e *edge) getNode() *node {
	return e.node
}

func (e *edge) getLabel() *label {
	return e.label
}
