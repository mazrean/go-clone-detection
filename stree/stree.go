package stree

import (
	"errors"
	"fmt"

	"github.com/mazrean/go-clone-detection/domain"
)

type STree struct {
	domainNodes   []*domain.Node
	root          *node
	leafNum       int64
	latestNode    *node
	latestNodeLen int64
}

func NewSTree() *STree {
	tree := &STree{
		domainNodes: []*domain.Node{},
		leafNum:     0,
	}

	rootNode := newRootNode(tree)

	tree.root = rootNode
	tree.latestNode = rootNode
	tree.latestNodeLen = 0

	return tree
}

func (st *STree) AddNode(newDomainNode *domain.Node) error {
	st.domainNodes = append(st.domainNodes, newDomainNode)

	if len(st.domainNodes) == 1 {
		leaf := newLeafNode(st, 0)
		e := newEdge(st, &label{0, finalIndex}, leaf)
		st.latestNode.addEdge(e)

		st.leafNum++
		st.latestNode = st.root
		st.latestNodeLen = 0

		return nil
	}

	for st.leafNum < int64(len(st.domainNodes)) {
		nowNode := st.latestNode
		// 現在位置のノードのrootからのトークン数
		nowNodeLen := st.latestNodeLen
		if st.latestNode.getNodeType() != rootNodeType {
			nowNode = nowNode.getSuffixLink()
			nowNodeLen = st.latestNodeLen - 1
		}

		domainNodes := st.domainNodes[nowNodeLen:]

		nowNodeLen += int64(len(domainNodes))
		nowNode, e, domainNodes, err := st.walk(nowNode, domainNodes)
		if err != nil {
			return fmt.Errorf("error walking: %v", err)
		}
		nowNodeLen -= int64(len(domainNodes))

		// エッジがみつからなかった場合、Rule2適用
		if e == nil {
			leaf := newLeafNode(st, st.leafNum)
			e := newEdge(st, &label{nowNodeLen, finalIndex}, leaf)
			nowNode.addEdge(e)

			st.leafNum++
			st.latestNode = nowNode
			st.latestNodeLen = nowNodeLen

			continue
		}

		domainNode := st.domainNodes[e.getLabel().start+int64(len(domainNodes))]

		if domainNode.GetNodeType() != newDomainNode.GetNodeType() || domainNode.GetToken() != newDomainNode.GetToken() {
			// エッジがみつかり、次の文字が適合しない場合も、Rule2適用

			splitPoint := e.getLabel().start + int64(len(domainNodes))
			suffixLink := nowNode.getSuffixLink()
			var emptyEdge *edge
			suffixLink, emptyEdge, _, err = st.walk(suffixLink, st.domainNodes[e.getLabel().start:splitPoint])
			if err != nil {
				return fmt.Errorf("error walking: %v", err)
			}
			if emptyEdge != nil {
				return errors.New("error walking: no edge")
			}

			nowNode, _ := e.splitEdge(splitPoint, suffixLink)
			nowNodeLen = nowNodeLen + int64(len(domainNodes)) - 1

			leaf := newLeafNode(st, st.leafNum)
			e := newEdge(st, &label{nowNodeLen, finalIndex}, leaf)
			nowNode.addEdge(e)

			st.leafNum++
			st.latestNode = nowNode
			st.latestNodeLen = nowNodeLen

			continue
		} else {
			// エッジがみつかり、次の文字が適合する場合は、Rule3適用
			break
		}
	}

	return nil
}

func (st *STree) walk(nd *node, domainNodes []*domain.Node) (*node, *edge, []*domain.Node, error) {
	e, err := nd.getEdgeByLabel(domainNodes[0])
	if errors.Is(err, ErrNoEdgeFound) {
		return nd, nil, domainNodes, nil
	}
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting edge by label: %v", err)
	}

	for e != nil && e.getLength() < int64(len(domainNodes)) {
		nd = e.getNode()
		domainNodes = domainNodes[e.getLength():]

		e, err = nd.getEdgeByLabel(domainNodes[0])
		if errors.Is(err, ErrNoEdgeFound) {
			return nd, nil, domainNodes, nil
		}
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error getting edge by label: %v", err)
		}
	}

	return nd, e, domainNodes, nil
}
