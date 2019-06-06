package Error

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	result, err := ErrorInt(RxError("rx-error")).ToSlice()
	expect := []int(nil)
	assert.EqualError(t, err, "rx-error")
	assert.Equal(t, expect, result)
}
