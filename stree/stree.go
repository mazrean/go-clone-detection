package stree

import (
	"errors"
	"fmt"
	"sort"

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

func (st *STree) GetClonePairs(threshold int) ([]*domain.CloneSequencePair, error) {
	/*
		考え方:
		- 各クローンをクローンの終了位置で識別する
		- startが大きいedgeのCloneにより長い(=startが小さい)Cloneがある場合上書き
			- sortしてstartが大きいCloneから検出
			- 被った場合上書き
		- startが同じで長さがより長いCloneはsuffix treeの中で長い方のみ残るので、発生しない
	*/
	var rootEdges []*edge
	copy(rootEdges, st.root.getEdges())
	sort.Slice(rootEdges, func(i, j int) bool {
		return rootEdges[i].getLabel().start > rootEdges[j].getLabel().start
	})

	cloneMaps := map[int]map[int]int{}
	for _, e := range rootEdges {
		nd := e.getNode()
		ndType := nd.getNodeType()
		if ndType == leafNodeType {
			continue
		}

		var err error
		_, cloneMaps, err = st.dfs(nd, threshold, int(e.getLength()), cloneMaps)
		if err != nil {
			return nil, fmt.Errorf("error dfs: %v", err)
		}
	}

	var clonePairs []*domain.CloneSequencePair
	for start1, cloneMap := range cloneMaps {
		for start2, length := range cloneMap {
			clonePairs = append(clonePairs, domain.NewCloneSequencePair(
				st.domainNodes[start1:start1+length],
				st.domainNodes[start2:start2+length],
			))
		}
	}

	return clonePairs, nil
}

func (st *STree) dfs(nd *node, threshold int, length int, cloneMap map[int]map[int]int) ([]int, map[int]map[int]int, error) {
	if nd.getNodeType() != internalNodeType {
		return nil, nil, errors.New("error dfs: not internal node")
	}

	//直下にあるleafの値
	directLeafs := []int{}
	//各区分(直下でない場合はedgeごと、直下の場合は直下のグループ)ごとのleafの値
	leafsList := [][]int{}
	for _, e := range nd.getEdges() {
		nd := e.getNode()
		ndType := nd.getNodeType()
		if ndType == leafNodeType {
			ndValue, err := nd.getValue()
			if err != nil {
				return nil, nil, fmt.Errorf("error getting value: %v", err)
			}

			directLeafs = append(directLeafs, int(ndValue))
		} else {
			var newLeafs []int
			var err error
			newLeafs, cloneMap, err = st.dfs(nd, threshold, length+int(e.getLength()), cloneMap)
			if err != nil {
				return nil, nil, fmt.Errorf("error dfs: %v", err)
			}

			leafsList = append(leafsList, newLeafs)
		}
	}

	leafsList = append(leafsList, directLeafs)

	if length > threshold {
		//直下のleaf間のペア検出
		for i, leaf1 := range directLeafs {
			for j := i + 1; j < len(directLeafs); j++ {
				cloneMap[leaf1+length][directLeafs[j]+length] = length
			}
		}

		//各区分のleaf間のペア検出
		for i, leafs := range leafsList {
			for j := i + 1; j < len(leafsList); j++ {
				for _, leaf1 := range leafs {
					for _, leaf2 := range leafsList[j] {
						cloneMap[leaf1+length][leaf2+length] = length
					}
				}
			}
		}
	}

	//子ノードにあるleafの値
	leafs := []int{}
	for _, leafList := range leafsList {
		leafs = append(leafs, leafList...)
	}

	return leafs, cloneMap, nil
}
