package Error

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_toSlice(t *testing.T) {
	result, err := ErrorInt(RxError("rx-error")).ToSlice()
	expect := []int(nil)
	assert.EqualError(t, err, "rx-error")
	assert.Equal(t, expect, result)
}

func TestError_println(t *testing.T) {
	err := ErrorInt(RxError("rx-error")).Println()
	assert.EqualError(t, err, "rx-error")
}
