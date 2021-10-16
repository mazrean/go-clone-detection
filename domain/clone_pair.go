package domain

type ClonePair struct {
	node1 *Node
	node2 *Node
}

func NewClonePair(node1 *Node, node2 *Node) *ClonePair {
	return &ClonePair{
		node1: node1,
		node2: node2,
	}
}

func (cp *ClonePair) GetNodes() (*Node, *Node) {
	return cp.node1, cp.node2
}
