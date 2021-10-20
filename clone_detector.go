package clone

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/mazrean/go-clone-detection/domain"
	"github.com/mazrean/go-clone-detection/serializer"
	"github.com/mazrean/go-clone-detection/stree"
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

	if config.Serializer == nil {
		config.Serializer = &serializer.Serializer{}
	}

	if config.SuffixTree == nil {
		config.SuffixTree = stree.NewSTree()
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
	cloneSequencePairs, err := cd.suffixTree.GetClonePairs(cd.config.Threshold)
	if err != nil {
		return nil, fmt.Errorf("suffix tree error: %v", err)
	}

	clonePairs := []*domain.ClonePair{}
	for _, cloneSequencePair := range cloneSequencePairs {
		sequence1, sequence2 := cloneSequencePair.GetNodes()
		var i int64
		for i = 0; i < int64(cloneSequencePair.GetLength()); {
			node1 := sequence1[i]

			if int64(node1.GetChildCount()) <= int64(len(sequence1)) {
				clonePairs = append(clonePairs, domain.NewClonePair(
					node1,
					sequence2[i],
				))
				i += int64(node1.GetChildCount())
			} else {
				break
			}
		}
	}

	return clonePairs, nil
}
