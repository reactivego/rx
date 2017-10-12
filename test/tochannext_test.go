package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ToChanNext is used to convert the observable to a channel of Next struct
// values. ToChanNext changes the scheduler (it uses for subscribing) to an
// asynchronous one, so the subscription runs concurrently with the channel
// reading for loop. Errors are emitted in-band with the data.
func Example_toChanNext() {
	for item := range Range(1, 9).ToChanNext() {
		if item.Err == nil {
			fmt.Print(item.Next)
		} else {
			fmt.Println(item.Err)
		}
	}

	// Output:
	// 123456789
}

func TestToChanNext(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5, 4, 3, 2, 1}
	a := FromSliceInt(expected).ToChanNext()
	b := []int{}
	for i := range a {
		b = append(b, i.Next)
	}
	assert.Equal(t, expected, b)
}

func Test_toChan