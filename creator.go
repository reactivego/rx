package observable

type Creator[T any] func() (Next T, Err error, Done bool)
