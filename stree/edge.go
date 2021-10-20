package stree

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

func (e *edge) splitEdge(i int64, newSuffixLink *node) (*node, *edge) {
	newEdge := newEdge(e.tree, newLabel(i, e.label.end), e.node)
	newNode := newInternalNode(e.tree, newSuffixLink, []*edge{newEdge})

	e.label = newLabel(e.label.start, i)
	e.node = newNode

	return newNode, newEdge
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
