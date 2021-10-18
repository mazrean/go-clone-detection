package domain

import (
	"go/ast"

	"github.com/mazrean/go-clone-detection/domain/values"
)

type Node struct {
	// 本来はドメインにあるべきでないが、idで対応関係を取るのが面倒なのでast.Nodeも入れる
	node       ast.Node
	nodeType   values.NodeType
	position   *values.Position
	childCount values.ChildCount
	token      values.NodeToken
}

func NewNode(
	node ast.Node,
	nodeType values.NodeType,
	position *values.Position,
	childCount values.ChildCount,
	token values.NodeToken,
) *Node {
	return &Node{
		node:       node,
		nodeType:   nodeType,
		position:   position,
		childCount: childCount,
		token:      token,
	}
}

func (n *Node) GetNode() ast.Node {
	return n.node
}

func (n *Node) GetNodeType() values.NodeType {
	return n.nodeType
}

func (n *Node) GetPosition() *values.Position {
	return n.position
}

func (n *Node) GetChildCount() values.ChildCount {
	return n.childCount
}

func (n *Node) IncrementChildCount(childCount values.ChildCount) {
	n.childCount += childCount
}

func (n *Node) GetToken() values.NodeToken {
	return n.token
}
