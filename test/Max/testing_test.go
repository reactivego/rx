package Max

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	result, err := FromInts(4, 5, 4, 3, 2, 1, 2).Max().ToSingle()
	expect := 5
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
