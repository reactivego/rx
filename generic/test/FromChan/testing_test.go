package FromChan

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromChan(t *testing.T) {
	ch := make(chan interface{}, 6)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	ch <- errors.New("error")
	close(ch)
	result, err := FromChan(ch).AsObservableInt().ToSlice()
	expect := []int{0, 1, 2, 3, 4}
	assert.EqualError(t, err, "error")
	assert.Equal(t, expect, result)
}

func TestFromChanInt(t *testing.T) {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	result, err := FromChanInt(ch).ToSlice()
	expect := []int{0, 1, 2, 3, 4}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}
