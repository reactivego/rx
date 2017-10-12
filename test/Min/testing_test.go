package Min

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	result, err := FromInts(4, 5, 4, 3, 2, 1, 2).Min().ToSingle()
	expect := 1
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
