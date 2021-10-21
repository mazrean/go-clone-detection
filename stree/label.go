package stree

import (
	"errors"
	"math"
)

const finalIndex = math.MaxInt64

type label struct {
	start, end int64
}

func newLabel(start, end int64) (*label, error) {
	if start < 0 || start > end {
		return nil, errors.New("invalid label range")
	}

	return &label{
		start: start,
		end:   end,
	}, nil
}
