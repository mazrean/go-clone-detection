package clone

import (
	"context"
	"go/ast"

	"github.com/mazrean/go-clone-detection/domain"
)

type Serializer interface {
	Serialize(ctx context.Context, root ast.Node, nodeChan chan<- *domain.Node) error
}
