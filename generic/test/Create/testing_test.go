package Create

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	s := CreateInt(func(observer IntObserver) {
		observer.Next(0)
		observer.Next(1)
		observer.Next(2)
		observer.Complete()
	})
	a, e := s.ToSlice()
	assert.NoError(t, e)
	b, e := s.ToSlice()
	assert.NoError(t, e)
	assert.Equal(t, []int{0, 1, 2}, a)
	assert.Equal(t, []int{0, 1, 2}, b)
}
