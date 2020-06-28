package test

import (
	"fmt"
	"sync"
	"testing"
)

/*
T is a struct embedding a *testing.T and therefore has the same methods as a
testing.T type. It provides support for writing Test Driven Development (TDD)
specs directly in Go. Is provides a domain-specific language patterned after
RSpec https://en.wikipedia.org/wiki/RSpec.

Use it as follows:

	func TestObservable_ConcatAll(e *testing.T) {
	     Describ(e, "emits", func(t T) {
	         Contex(t, "with a fast source", func(t T) {
	             Contex(t, "and a slow concat", func(t T) {
	                 I(t, "should complete before the first emit to observer", func(t T) {
	                     Asser(t).NoError(nil)
	                 })
	             })
	         })
	     })
	}
*/
type T struct {
	*testing.T
	observations map[string]*observation
}

type observation struct {
	T
	sync.Mutex
	actual []interface{}
	err    error
	done   bool
}

// Describ(e) is used to describe a data type, method or an example group. This
// is the outer block which actually contains the test code and it depicts the
// characteristics of the code enclosed in it.  The name argument is a reference
// to the module or data type being tested.
func Describ(e *testing.T, name string, f func(T)) bool {
	e.Helper()
	return e.Run(name, func(t *testing.T) { f(T{T: t}) })
}

// Contex(t) block is used to describe the context in which the data type or
// method mentioned in the describe block is being used. The context description
// argument is used to fence off a specific context in which specific examples
// are to be performed.
func Contex(t T, description string, f func(T)) bool {
	t.Helper()
	return t.Run(description, func(t *testing.T) { f(T{T: t}) })
}

// I(t) describes the specification of the example in the context. The
// description argument can be considered as the function that the block is
// expected to perform or in other words it can be considered as a test case.
func I(t T, description string, f func(T)) bool {
	t.Helper()
	return t.Run(description, func(t *testing.T) {
		f(T{T: t, observations: make(map[string]*observation)})
	})
}

// Assert returns an assertion to me made.
func Asser(t T) Assertion {
	return Assertion{t.T}
}

// Assertion represents all assertians that can be made.
type Assertion struct {
	*testing.T
}

// Must asserts that condition is true, when it is false it calls Fatal with:
//    Must have {args}
func (a Assertion) Must(condition bool, args ...interface{}) {
	if !condition {
		a.Helper()
		a.Fatal("Must have " + fmt.Sprint(args...))
	}
}

// Not asserts that condition is false, when condition is true it calls Fatal with:
//  Must not have {args}
func (a Assertion) Not(condition bool, args ...interface{}) {
	if condition {
		a.Helper()
		a.Fatal("Must not have " + fmt.Sprint(args...))
	}
}

// True asserts that condition is true, when it is false it calls Fatal with:
//  Expected {args}
func (a Assertion) True(condition bool, args ...interface{}) {
	if !condition {
		a.Helper()
		a.Fatal("Expected " + fmt.Sprint(args...))
	}
}

// False asserts that condition is false, when condition is true it calls Fatal with:
//  Unexpected {args}
func (a Assertion) False(condition bool, args ...interface{}) {
	if condition {
		a.Helper()
		a.Fatal("Unexpected " + fmt.Sprint(args...))
	}
}

// Error asserts that err is not nil, but if it is nil, it calls Fatal with:
//  An error is expected but got nil. {args}
func (a Assertion) Error(err error, args ...interface{}) {
	if err == nil {
		a.Helper()
		a.Fatal("An error is expected but got nil. " + fmt.Sprint(args...))
	}
}

// NoError asserts that err must be nil, when it isn't it calls Fatal with:
//    No error was expected but got: {err}
//    {args}
func (a Assertion) NoError(err error, args ...interface{}) {
	if err != nil {
		a.Helper()
		a.Fatal(fmt.Sprintf("No error was expected but got: %s\n", err.Error()) + fmt.Sprint(args...))
	}
}

// Equal assertian compares the type and value of actual and expect and when
// not equal will call fatal with:
//	Not equal:
//	  expect: {expect}
//	  actual: {actual} {args}
func (a Assertion) Equal(actual, expect interface{}, args ...interface{}) {
	a.Helper()
	tactual := fmt.Sprintf("%T", actual)
	texpect := fmt.Sprintf("%T", expect)
	if tactual != texpect {
		a.Fatal(fmt.Sprintf("Not equal: \n"+
			"expect: %s\n"+
			"actual: %s %s", texpect, tactual, fmt.Sprint(args...)))
	}
	vactual := fmt.Sprintf("%+v", actual)
	vexpect := fmt.Sprintf("%+v", expect)
	if vactual != vexpect {
		a.Fatal(fmt.Sprintf("Not equal: \n"+
			"expect: %s\n"+
			"actual: %s %s", vexpect, vactual, fmt.Sprint(args...)))
	}
}

type expect struct{ T }

func Expec(t T) expect {
	return expect{t}
}

type observer = func(next interface{}, err error, done bool)
type intobserver = func(next int, err error, done bool)

func (e expect) MakeObservation(name string) observer {
	observation := &observation{T: e.T}
	e.T.observations[name] = observation
	observer := func(next interface{}, err error, done bool) {
		observation.Lock()
		if !done {
			observation.actual = append(observation.actual, next)
		} else {
			observation.err = err
			observation.done = true
		}
		observation.Unlock()
	}
	return observer
}

func (e expect) MakeIntObservation(name string) intobserver {
	observation := &observation{T: e.T}
	e.observations[name] = observation
	observer := func(next int, err error, done bool) {
		observation.Lock()
		if !done {
			observation.actual = append(observation.actual, next)
		} else {
			observation.err = err
			observation.done = true
		}
		observation.Unlock()
	}
	return observer
}

func (e expect) Observation(name string) *observation {
	return e.observations[name]
}

func (e expect) IntObservation(name string) *observation {
	return e.observations[name]
}

func (o *observation) ToBe(expect ...interface{}) {
	o.Helper()
	o.Lock()
	defer o.Unlock()
	Asser(o.T).Equal(o.actual, expect)
	Asser(o.T).Equal(o.err, nil, "for error")
	Asser(o.T).Equal(o.done, true, "for done")
}

func (o *observation) ToBeError(expect error) {
	o.Helper()
	Asser(o.T).Equal(o.err, expect, "for error")
}

func (o *observation) ToComplete() {
	o.Helper()
	o.Lock()
	defer o.Unlock()
	Asser(o.T).Equal(o.err, nil, "for error, expected to complete")
	Asser(o.T).Equal(o.done, true, "for done, expected to complete")
}

func (o *observation) ToBeActive() {
	o.Helper()
	o.Lock()
	defer o.Unlock()
	Asser(o.T).Equal(o.done, false, "for done, expected to be active")
}
