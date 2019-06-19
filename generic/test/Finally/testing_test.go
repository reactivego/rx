package Finally

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinally(t *testing.T) {
	finally := false
	_, err := ThrowInt(RxError("error")).Finally(func() { finally = true }).ToSlice()
	assert.True(t, finally)
	assert.Error(t, err)

	finally = false
	result, err := EmptyInt().Finally(func() { finally = true }).ToSlice()
	assert.True(t, finally)
	assert.Equal(t, []int(nil), result)
	assert.NoError(t, err)
}
