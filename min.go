package rx

import "cmp"

func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
