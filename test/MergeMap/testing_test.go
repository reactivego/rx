package MergeMap

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeMap(t *testing.T) {
	result, err := Range(1, 2).MergeMapInt(func(n int) ObservableInt { return Range(n, 2) }).ToSlice()
	sort.Ints(result)
	expect := []int{1, 2, 2, 3}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
