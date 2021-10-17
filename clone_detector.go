package clone

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/mazrean/go-clone-detection/domain"
	"golang.org/x/sync/errgroup"
)

type CloneDetector struct {
	config     *Config
	serializer Serializer
	suffixTree SuffixTree
}

func NewCloneDetector(config *Config) *CloneDetector {
	if config == nil {
		config = DefaultConfig
	}

	return &CloneDetector{
		config: config,
	}
}

func (cd *CloneDetector) AddNode(ctx context.Context, root ast.Node) error {
	nodeChan := make(chan *domain.Node)

	eg := errgroup.Group{}
	eg.Go(func() error {
		defer close(nodeChan)

		err := cd.serializer.Serialize(ctx, root, nodeChan)
		if err != nil {
			return fmt.Errorf("serialization error: %v", err)
		}

		return nil
	})

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case node, ok := <-nodeChan:
				if !ok {
					return nil
				}

				err := cd.suffixTree.AddNode(node)
				if err != nil {
					return fmt.Errorf("suffix tree error: %v", err)
				}
			}
		}
	})

	err := eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (cd *CloneDetector) GetClones() ([]*domain.ClonePair, error) {
	clonePairs, err := cd.suffixTree.GetClonePairs(cd.config.Threshold)
	if err != nil {
		return nil, fmt.Errorf("suffix tree error: %v", err)
	}

	return clonePairs, nil
}
