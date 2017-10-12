package Replay

import (
	"testing"
	// "time"

	"github.com/stretchr/testify/assert"
)

// func TestReplay(t *testing.T) {
// 	ch := make(chan int, 5)
// 	for i := 0; i < 5; i++ {
// 		ch <- i
// 	}
// 	close(ch)
// 	source := FromChanInt(ch).Replay(0, 0)

// 	println("source ready")

// 	source.Connect()

// 	println("connected")
// 	expected := []int{0, 1, 2, 3, 4}
// 	result, err := source.ToSlice()
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, result)
// 	result, err = source.ToSlice()
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected, result)
// }

func TestReplayWithSize(t *testing.T) {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	source := FromChanInt(ch).Replay(2, 0).AutoConnect(1)
	println("before ToSlice")
	result, err := source.ToSlice()
	println("after ToSlice")
	expect := []int{0, 1, 2, 3, 4}
	assert.NoError(t, err)
	assert.Equal(t, expect, result)
	println("done")

	// result, err = source.ToSlice()
	// expect = []int{3, 4}
	// assert.NoError(t, err)
	// assert.Equal(t, expect, result)

	// result, err = source.ToSlice()
	// expect = []int{3, 4}
	// assert.NoError(t, err)
	// assert.Equal(t, expect, result)
}

// func TestReplayWithExpiry(t *testing.T) {
// 	ch := make(chan int)
// 	go func() {
// 		for i := 0; i < 5; i++ {
// 			ch <- i
// 			time.Sleep(100 * time.Millisecond)
// 		}
// 		close(ch)
// 	}()
// 	source := FromChanInt(ch).Replay(0, 600*time.Millisecond)
// 	time.Sleep(500 * time.Millisecond)

// 	result, err := source.ToSlice()
// 	expect := []int{0, 1, 2, 3, 4}
// 	assert.NoError(t, err)
// 	assert.Equal(t, expect, result)

// 	time.Sleep(100 * time.Millisecond)

// 	result, err = source.ToSlice()
// 	expect = []int{1, 2, 3, 4}
// 	assert.NoError(t, err)
// 	assert.Equal(t, expect, result)
// }
