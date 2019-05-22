package Count

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	count := FromInts(1, 2, 3, 4, 5, 6, 7).Count()
	result, err := count.ToSlice()
	expect := []int{7}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)

	// resubscribe, expect that it will return 7 again
	b, err := count.ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expect, b)
}
