package clone

import (
	"context"
	"errors"
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
		config:     config,
		serializer: config.Serializer,
		suffixTree: config.SuffixTree,
	}
}

func (cd *CloneDetector) AddNode(ctx context.Context, root ast.Node) error {
	nodeChan := make(chan *domain.Node)

	if root == nil {
		return errors.New("root node is nil")
	}

	eg, ctx := errgroup.WithContext(ctx)
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

type ClonePair struct {
	Node1 ast.Node
	Node2 ast.Node
}

func (cd *CloneDetector) GetClones() ([]*ClonePair, error) {
	cloneSequencePairs, err := cd.suffixTree.GetClonePairs(cd.config.Threshold)
	if err != nil {
		return nil, fmt.Errorf("suffix tree error: %v", err)
	}

	clonePairs := []*ClonePair{}
	for _, cloneSequencePair := range cloneSequencePairs {
		sequence1, sequence2 := cloneSequencePair.GetNodes()
		var i int64
		for i = 0; i < int64(cloneSequencePair.GetLength()); {
			node1 := sequence1[i]

			if int64(node1.GetChildCount()) <= int64(len(sequence1)) {
				clonePairs = append(clonePairs, &ClonePair{
					node1.GetNode(),
					sequence2[i].GetNode(),
				})
				i += int64(node1.GetChildCount()) + 1
			} else {
				break
			}
		}
	}

	return clonePairs, nil
}
