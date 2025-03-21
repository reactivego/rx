package rx

import "cmp"

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
