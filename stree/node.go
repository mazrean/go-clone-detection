package stree

import (
	"database/sql"
	"errors"
	"sort"

	"github.com/mazrean/go-clone-detection/domain"
)

type nodeType int

const (
	rootNodeType nodeType = iota
	internalNodeType
	leafNodeType
)

type node struct {
	tree       *STree
	nodeType   nodeType
	value      sql.NullInt64
	edges      []*edge
	suffixLink *node
}

func newRootNode(tree *STree) *node {
	return &node{
		tree:     tree,
		nodeType: rootNodeType,
	}
}

func newInternalNode(tree *STree, suffixLink *node, edges []*edge) *node {
	return &node{
		tree:       tree,
		nodeType:   internalNodeType,
		edges:      edges,
		suffixLink: suffixLink,
	}
}

func newLeafNode(tree *STree, value int64) *node {
	return &node{
		tree:     tree,
		nodeType: leafNodeType,
		value: sql.NullInt64{
			Int64: value,
			Valid: true,
		},
	}
}

func (n *node) getNodeType() nodeType {
	return n.nodeType
}

func (n *node) getEdges() []*edge {
	return n.edges
}

var (
	ErrNoEdgeFound = errors.New("no edge found")
)

func (n *node) getEdgeByLabel(domainNode *domain.Node) (*edge, error) {
	if n.nodeType != internalNodeType {
		return nil, errors.New("node is not an internal node")
	}

	id := sort.Search(len(n.edges), func(i int) bool {
		if n.tree.domainNodes[n.edges[i].label.start].GetNodeType() == domainNode.GetNodeType() {
			return n.tree.domainNodes[n.edges[i].label.start].GetToken() >= domainNode.GetToken()
		}

		return n.tree.domainNodes[n.edges[i].label.start].GetNodeType() >= domainNode.GetNodeType()
	})

	if n.tree.domainNodes[n.edges[id].label.start].GetNodeType() != domainNode.GetNodeType() || n.tree.domainNodes[n.edges[id].label.start].GetToken() != domainNode.GetToken() {
		return nil, ErrNoEdgeFound
	}

	return n.edges[id], nil
}

func (n *node) addEdge(e *edge) {
	n.edges = append(n.edges, e)

	sort.Slice(n.edges, func(i, j int) bool {
		if n.tree.domainNodes[n.edges[i].label.start].GetNodeType() == n.tree.domainNodes[n.edges[j].label.start].GetNodeType() {
			return n.tree.domainNodes[n.edges[i].label.start].GetToken() < n.tree.domainNodes[n.edges[j].label.start].GetToken()
		}

		return n.tree.domainNodes[n.edges[i].label.start].GetNodeType() < n.tree.domainNodes[n.edges[j].label.start].GetNodeType()
	})
}

func (n *node) getSuffixLink() *node {
	return n.suffixLink
}

func (n *node) getValue() (int64, error) {
	if n.nodeType != leafNodeType {
		return 0, errors.New("node is not a leaf node")
	}

	return n.value.Int64, nil
}
