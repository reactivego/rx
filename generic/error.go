package rx

//jig:template RxError

type RxError string

func (e RxError) Error() string { return string(e) }