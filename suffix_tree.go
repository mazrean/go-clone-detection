package clone

import (
	"github.com/mazrean/go-clone-detection/domain"
)

type SuffixTree interface {
	AddNode(node *domain.Node) error
	GetClonePairs(threshold int) ([]*domain.CloneSequencePair, error)
}
