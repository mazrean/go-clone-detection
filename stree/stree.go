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
	nextNode      *node
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

	for st.leafNum < int64(len(st.domainNodes)) {
		nowNode := st.latestNode
		// 現在位置のノードのrootからのトークン数
		nowNodeLen := st.latestNodeLen
		if nowNode.getNodeType() != rootNodeType {
			nowNode = nowNode.getSuffixLink()
			nowNodeLen = st.latestNodeLen - 1
		}
		oldNextNode := st.nextNode

		restDomainNodes := st.domainNodes[st.leafNum+nowNodeLen:]

		nowNodeLen += int64(len(restDomainNodes))
		nowNode, e, restDomainNodes, err := st.walk(nowNode, restDomainNodes)
		if err != nil {
			return fmt.Errorf("error walking(node): %w", err)
		}
		nowNodeLen -= int64(len(restDomainNodes))

		// エッジがみつからなかった場合、Rule2適用
		if e == nil && len(restDomainNodes) > 0 {
			if st.nextNode != nil {
				return errors.New("nextNode is not nil")
			}

			l, err := newLabel(int64(len(st.domainNodes))-1, finalIndex)
			if err != nil {
				return fmt.Errorf("error creating label(no edge): %w", err)
			}

			leaf := newLeafNode(st, st.leafNum)
			e := newEdge(st, l, leaf)
			err = nowNode.addEdge(e)
			if err != nil {
				return fmt.Errorf("error adding edge(no edge): %w", err)
			}

			st.leafNum++
			st.latestNode = nowNode
			st.latestNodeLen = nowNodeLen
			st.nextNode = nil

			continue
		}

		if e == nil {
			// ちょうどノード上まで行った時、Rule3適用
			break
		}

		domainNode := st.domainNodes[e.getLabel().start+int64(len(restDomainNodes))-1]

		if domainNode.GetNodeType() != newDomainNode.GetNodeType() ||
			domainNode.GetToken() != newDomainNode.GetToken() ||
			domainNode.GetChildCount() != newDomainNode.GetChildCount() {
			// エッジがみつかり、次の文字が適合しない場合も、Rule2適用

			splitPoint := e.getLabel().start + int64(len(restDomainNodes)) - 1
			var suffixLink *node
			var linkDomainNodes []*domain.Node
			if nowNode.getNodeType() != rootNodeType {
				suffixLink = nowNode.getSuffixLink()
				suffixLink, _, linkDomainNodes, err = st.walk(suffixLink, st.domainNodes[e.getLabel().start:splitPoint])
				if err != nil {
					return fmt.Errorf("error walking(suffix tree): %w", err)
				}
			} else {
				if e.getLabel().start+1 == splitPoint {
					suffixLink = st.root
				} else {
					suffixLink, _, linkDomainNodes, err = st.walk(st.root, st.domainNodes[e.getLabel().start+1:splitPoint])
					if err != nil {
						return fmt.Errorf("error walking(suffix tree): %w", err)
					}
				}
			}

			if len(linkDomainNodes) > 0 {
				/*
					次のノード追加までにsuffix linkのノードは作られるので,
					メモリの確保のみしておく
				*/
				st.nextNode = &node{}
				suffixLink = st.nextNode
			} else {
				st.nextNode = nil
			}

			newNode, _, err := e.splitEdge(splitPoint, suffixLink)
			if err != nil {
				return fmt.Errorf("error splitting edge: %w", err)
			}
			newNodeLen := nowNodeLen + e.getLength()

			if oldNextNode != nil {
				*oldNextNode = *newNode
				newNode = oldNextNode
				e.node = newNode
			}

			if len(linkDomainNodes) == 0 {
				nowNode = newNode
				nowNodeLen = newNodeLen
			}

			l, err := newLabel(int64(len(st.domainNodes))-1, finalIndex)
			if err != nil {
				return fmt.Errorf("error creating label(char): %w", err)
			}

			leaf := newLeafNode(st, st.leafNum)
			e := newEdge(st, l, leaf)
			err = newNode.addEdge(e)
			if err != nil {
				return fmt.Errorf("error adding edge(char): %w", err)
			}

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
		return nil, nil, nil, fmt.Errorf("error getting edge by label: %w", err)
	}

	for e.getLength() < int64(len(domainNodes)) {
		nd = e.getNode()
		domainNodes = domainNodes[e.getLength():]

		e, err = nd.getEdgeByLabel(domainNodes[0])
		if errors.Is(err, ErrNoEdgeFound) {
			return nd, nil, domainNodes, nil
		}
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error getting edge by label: %w", err)
		}
	}

	edgeLastNode := st.domainNodes[e.getLabel().start+int64(len(domainNodes))-1]
	restLastNode := domainNodes[len(domainNodes)-1]
	if e.getLength() == int64(len(domainNodes)) &&
		edgeLastNode.GetNodeType() == restLastNode.GetNodeType() &&
		edgeLastNode.GetToken() == restLastNode.GetToken() &&
		edgeLastNode.GetChildCount() == restLastNode.GetChildCount() {
		return e.getNode(), nil, domainNodes[e.getLength():], nil
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

	rootEdges := make([]*edge, len(st.root.getEdges()))
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
			return nil, fmt.Errorf("error dfs: %w", err)
		}
	}

	var clonePairs []*domain.CloneSequencePair
	for end1, cloneMap := range cloneMaps {
		for end2, length := range cloneMap {
			clonePairs = append(clonePairs, domain.NewCloneSequencePair(
				st.domainNodes[end1-length:end1],
				st.domainNodes[end2-length:end2],
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
				return nil, nil, fmt.Errorf("error getting value: %w", err)
			}

			directLeafs = append(directLeafs, int(ndValue))
		} else {
			var newLeafs []int
			var err error
			newLeafs, cloneMap, err = st.dfs(nd, threshold, length+int(e.getLength()), cloneMap)
			if err != nil {
				return nil, nil, fmt.Errorf("error dfs: %w", err)
			}

			leafsList = append(leafsList, newLeafs)
		}
	}

	leafsList = append(leafsList, directLeafs)

	if length > threshold {
		//直下のleaf間のペア検出
		for i, leaf1 := range directLeafs {
			for j := i + 1; j < len(directLeafs); j++ {
				_, ok := cloneMap[leaf1+length]
				if !ok {
					cloneMap[leaf1+length] = map[int]int{}
				}
				cloneMap[leaf1+length][directLeafs[j]+length] = length
			}
		}

		//各区分のleaf間のペア検出
		for i, leafs := range leafsList {
			for j := i + 1; j < len(leafsList); j++ {
				for _, leaf1 := range leafs {
					for _, leaf2 := range leafsList[j] {
						_, ok := cloneMap[leaf1+length]
						if !ok {
							cloneMap[leaf1+length] = map[int]int{}
						}
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
