package Scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	add := func(acc int, value int) int {
		return acc + value
	}
	result, err := FromInts(1, 2, 3, 4, 5).ScanInt(add, 0).ToSlice()
	expect := []int{1, 3, 6, 10, 15}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
