package Catch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatch(t *testing.T) {
	o123 := FromInts(1, 2, 3)
	o45 := FromInts(4, 5)
	oThrowError := ThrowInt(RxError("error"))

	a, err := o123.Concat(oThrowError).Catch(o45).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}
