package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ReplayChan is a channel implementation with replay capabilities. This channel
// is created with a buffer of a fixed size. Multiple endpoints can be created
// to read from the channel. If an endpoint can't keep up with the rate at which
// data is send, the endpoint will fail with an overflow error.
func Example_replayChan() {
	cb := NewReplayChanInt(2, 0)

	cb.Send(1)
	cb.Send(2)

	ep := cb.NewEndpoint()
	for value, ok := ep.Recv(); ok; value, ok = ep.Recv() {
		fmt.Println(value)
	}

	fmt.Println("-")

	cb.Send(3)
	cb.Send(4)

	for value, ok := ep.Recv(); ok; value, ok = ep.Recv() {
		fmt.Println(value)
	}

	fmt.Println("----")

	ep = cb.NewEndpoint()
	for value, ok := ep.Recv(); ok; value, ok = ep.Recv() {
		fmt.Println(value)
	}

	fmt.Println("----")

	cb.Send(5)
	cb.Close(nil)

	ep = cb.NewEndpoint()
	ep.Range(func(value int, err error, closed bool) bool {
		switch {
		case !closed:
			fmt.Println(value)
		case err != nil:
			fmt.Println(err)
		default:
			fmt.Println("closed")
		}
		return true
	})

	// Output:
	// 1
	// 2
	// -
	// 3
	// 4
	// ----
	// 3
	// 4
	// ----
	// 4
	// 5
	// closed
}

func TestReplayChan(t *testing.T) {
	bc := NewReplayChanInt(4, 0)
	bc.Send(1)
	assert.Equal(t, []int64{5, 0, 1}, []int64{bc.size, bc.read, bc.write})
	bc.Send(2)
	assert.Equal(t, []int64{5, 0, 2}, []int64{bc.size, bc.read, bc.write})
	bc.Send(3)
	assert.Equal(t, []int64{5, 0, 3}, []int64{bc.size, bc.read, bc.write})
	bc.Send(4)
	assert.Equal(t, []int64{5, 0, 4}, []int64{bc.size, bc.read, bc.write})
	// bc.Send(5)
	// assert.Equal(t, []int64{5, 1, 5}, []int64{bc.size, bc.read, bc.write})
	// bc.Send(6)
	// assert.Equal(t, []int64{5, 2, 6}, []int64{bc.size, bc.read, bc.write})
	// bc.Send(7)
	// assert.Equal(t, []int64{5, 3, 7}, []int64{bc.size, bc.read, bc.write})

	// ep := bc.NewEndpoint()
	// assert.Equal(t, int64(3), ep.cursor)
	// assert.Nil(t, ep.overflow)

	// value, ok := ep.Recv()
	// assert.Equal(t, 4, value)
	// assert.True(t, ok)

	// bc.Send(8)
	// assert.Equal(t, []int64{5, 4, 8}, []int64{bc.size, bc.read, bc.write})

	// assert.Equal(t, int64(4), ep.cursor)
	// assert.Nil(t, ep.overflow)

	// sending again would block on ep not being able to keepup.

	// bc.RemoveEndpoint(ep)

	// bc.Send(9)
	// assert.Equal(t, []int64{5, 5, 9}, []int64{bc.size, bc.read, bc.write})
}
