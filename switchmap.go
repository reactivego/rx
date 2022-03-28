package observable

func SwitchMap[T, U any](o Observable[T], project func(T) Observable[U]) Observable[U] {
	return SwitchAll(Map(o, project))
}
