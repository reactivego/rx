package x

import "fmt"

type Error string

func (e Error) Error() string { return string(e) }

func Errorf(format string, items ...any) error {
	return Error(fmt.Sprintf(format, items...))
}

func Errorw(err error) error {
	return Error(fmt.Sprintf("%[1]T:%[1]v", err))
}
