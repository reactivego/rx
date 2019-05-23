package Do

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	expect := []int{}
	result, err := FromInts(1, 2, 3, 4, 5).Do(func(v int) {
		expect = append(expect, v)
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
