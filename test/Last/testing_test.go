package Last

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLast(t *testing.T) {
	result, err := FromInts(1, 2, 3, 4).Last().ToSingle()
	expect := 4
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
