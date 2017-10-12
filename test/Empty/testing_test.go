package Empty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	result, err := EmptyInt().ToSlice()
	expect := []int(nil)
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
