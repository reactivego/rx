package AsyncSubject

import (
	"time"

	_ "github.com/reactivego/rx"
)

func Example() {
	subject := NewAsyncSubjectString()
	subject.Next("a")
	subscription := subject.Timeout(time.Millisecond).Subscribe(PrintlnString())
	subject.Next("b")
	subject.Next("c")
	subscription.Wait()
	// Output:
}

func Example1() {
	subject := NewAsyncSubjectString()
	subject.Next("a")
	subscription := subject.Subscribe(PrintlnString())
	subject.Next("b")
	subject.Next("c")
	subject.Complete()
	subscription.Wait()
	// Output: c
}
