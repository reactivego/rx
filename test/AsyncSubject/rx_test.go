package AsyncSubject

import (
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_asyncSubject() {
	subject := NewAsyncSubjectString()
	subject.Next("a")
	subscription := subject.Subscribe(PrintlnString())
	subject.Next("b")
	subject.Next("c")
	subject.Complete()
	subscription.Wait()
	// Output: c
}

func Example_asyncSubjectTimeout() {
	subject := NewAsyncSubjectString()
	subject.Next("a")
	subscription := subject.Timeout(time.Millisecond).Subscribe(PrintlnString())
	subject.Next("b")
	subject.Next("c")
	subscription.Wait()
	// Output:
}
