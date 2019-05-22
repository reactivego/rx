package Map

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	var result []string
	FromInts(1, 2, 3, 4).MapString(func(i int) string {
		return fmt.Sprintf("%d!", i)
	}).SubscribeNext(func(next string) {
		result = append(result, next)
	})
	expect := []string{"1!", "2!", "3!", "4!"}
	assert.Equal(t, expect, result)
}
