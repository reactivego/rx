package Concat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5}
	c := []int{6, 7}
	oa := FromSliceInt(a)
	ob := FromSliceInt(b)
	oc := FromSliceInt(c)
	s, _ := oa.Concat(ob).Concat(oc).ToSlice()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
	s, _ = oa.Concat(ob, oc).ToSlice()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
}
