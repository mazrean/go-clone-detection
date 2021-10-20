package domain

type CloneSequencePair struct {
	node1 []*Node
	node2 []*Node
}

func NewCloneSequencePair(node1 []*Node, node2 []*Node) *CloneSequencePair {
	return &CloneSequencePair{
		node1: node1,
		node2: node2,
	}
}

func (cp *CloneSequencePair) GetNodes() ([]*Node, []*Node) {
	return cp.node1, cp.node2
}

func (cp *CloneSequencePair) GetLength() int {
	return len(cp.node1)
}
