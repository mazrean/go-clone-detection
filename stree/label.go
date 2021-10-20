package stree

import "math"

const finalIndex = math.MaxInt64

type label struct {
	start, end int64
}

func newLabel(start, end int64) *label {
	return &label{
		start: start,
		end:   end,
	}
}
