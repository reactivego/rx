package Do

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	expect := []int{}
	result, err := FromInts(1, 2, 3, 4, 5).Do(func(v int) {
		expect = append(expect, v)
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

func TestDoOnError(t *testing.T) {
	var oerr error
	_, err := ThrowInt(errors.New("error")).DoOnError(func(err error) { oerr = err }).ToSlice()
	assert.Equal(t, err, oerr)
}

func TestDoOnComplete(t *testing.T) {
	complete := false
	result, err := EmptyInt().DoOnComplete(func() { complete = true }).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int(nil), result)
	assert.True(t, complete)
}

func TestDoOnCompleteSubscribeNext(t *testing.T) {
	wait := make(chan struct{})
	result := []int{}
	_ = FromInts(1, 2, 3, 4, 5).
		DoOnComplete(func() { close(wait) }).
		SubscribeNext(func(v int) { result = append(result, v) }, SubscribeOn(NewGoroutine()))
	<-wait
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestFinally(t *testing.T) {
	finally := false
	_, err := ThrowInt(errors.New("error")).Finally(func() { finally = true }).ToSlice()
	assert.True(t, finally)
	assert.Error(t, err)

	finally = false
	result, err := EmptyInt().Finally(func() { finally = true }).ToSlice()
	assert.True(t, finally)
	assert.Equal(t, []int(nil), result)
	assert.NoError(t, err)
}
