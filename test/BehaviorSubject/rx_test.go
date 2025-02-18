package BehaviorSubject

import (
	"time"

	_ "github.com/reactivego/rx/generic"
)

func Example_behaviorSubject() {
	subject := NewBehaviorSubjectString("a")
	subject.Timeout(1 * time.Millisecond).Println()
	// Output:
	// a
}

func Example_behaviorSubject2() {
	subject := NewBehaviorSubjectString("a")
	subject.Next("b")
	subject.Timeout(1 * time.Millisecond).Println()
	// Output:
	// b
}

func Example_behaviorSubject3() {
	concurrent := GoroutineScheduler()
	subject := NewBehaviorSubjectString("a")
	subject.Next("b")
	subject.Next("c")
	subject.Next("d")
	subject.Next("c")
	subject.Next("a")
	subject.Next("c")
	subject.Next("d")
	subject.Next("b")
	subject.Timeout(1 * time.Millisecond).SubscribeOn(concurrent).Subscribe(PrintlnString())
	subject.Next("c")
	subject.Next("d")
	concurrent.Wait()
	// Output:
	// b
	// c
	// d
}

func Example_behaviorSubjectCompleted() {
	const ms = time.Millisecond
	subject := NewBehaviorSubjectString("a")
	subject.Next("b")
	subject.Next("c")
	subject.Complete()
	subject.Println()
	// Output:
}
