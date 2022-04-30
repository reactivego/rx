package observable

type Creator[T any] func(index int) (Next T, Err error, Done bool)
