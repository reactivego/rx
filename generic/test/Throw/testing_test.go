package Throw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_println(t *testing.T) {
	err := ErrorInt(RxError("throw")).Println()
	assert.EqualError(t, err, "throw")
}
