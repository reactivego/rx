package rx

type Tuple2[T, U any] struct {
	First  T
	Second U
}

type Tuple3[T, U, V any] struct {
	First  T
	Second U
	Third  V
}

type Tuple4[T, U, V, W any] struct {
	First  T
	Second U
	Third  V
	Fourth W
}

type Tuple5[T, U, V, W, X any] struct {
	First  T
	Second U
	Third  V
	Fourth W
	Fifth  X
}
