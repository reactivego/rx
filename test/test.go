package test

import (
	"fmt"
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
type T struct{ *testing.T }

// Describ(e) is used to describe a data type, method or an example group. This
// is the outer block which actually contains the test code and it depicts the
// characteristics of the code enclosed in it.  The name argument is a reference
// to the module or data type being tested.
func Describ(e *testing.T, name string, f func(T)) bool {
	e.Helper()
	return e.Run(name, func(t *testing.T) { f(T{t}) })
}

// Contex(t) block is used to describe the context in which the data type or
// method mentioned in the describe block is being used. The context description
// argument is used to fence off a specific context in which specific examples
// are to be performed.
func Contex(t T, description string, f func(T)) bool {
	t.Helper()
	return t.Run(description, func(t *testing.T) { f(T{t}) })
}

// I(t) describes the specification of the example in the context. The
// description argument can be considered as the function that the block is
// expected to perform or in other words it can be considered as a test case.
func I(t T, description string, f func(T)) bool {
	t.Helper()
	return t.Run(description, func(t *testing.T) { f(T{t}) })
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

type observation struct {
	t      T
	actual []interface{}
	err    error
}

type expect struct {
	t            T
	observations map[string]*observation
}

func Expec(t T) expect {
	return expect{
		t: t,
		observations: make(map[string]*observation),
	}
}

type observer = func(next interface{}, err error, done bool)

func (e expect) MakeObservation(name string) observer {
	observation := &observation{t: e.t}
	e.observations[name] = observation
	observer := func(next interface{}, err error, done bool) {
		if !done {
			observation.actual = append(observation.actual, next)
		} else {
			observation.err = err
		}
	}
	return observer
}

func (e expect) Observation(name string) *observation {
	return e.observations[name]
}

func (o *observation) ToBe(expect []interface{}) {
	o.t.Helper()
	Asser(o.t).Equal(o.actual,expect)
}

func (e expect) Error(expect error) observer {
	observer := func(next interface{}, err error, done bool) {
		if done {
			Asser(e.t).Equal(err, expect)
		}
	}
	return observer
}

func (e expect) Completion() observer {
	observer := func(next interface{}, err error, done bool) {
		if done {
			Asser(e.t).NoError(err)
		}
	}
	return observer
}

func (e expect) Never() observer {
	observer := func(next interface{}, err error, done bool) {
		if done {
			Asser(e.t).NoError(err)
		}
	}
	return observer
}

func (e expect) Slice(expect []interface{}) observer {
	i := 0
	actual := []interface{}(nil)
	observer := func(next interface{}, err error, done bool) {
		failed := false
		if !done {
			actual = append(actual, next)
			failed = (i >= len(expect) || next != expect[i])
			i++
		} else {
			failed = (i != len(expect))
		}
		if failed {
			Asser(e.t).Equal(actual, expect)
		}
	}
	return observer
}

func (e expect) Value(expect interface{}) observer {
	initial := true
	observer := func(actual interface{}, err error, done bool) {
		switch {
		case !done:
			if initial {
				initial = false
				vactual := fmt.Sprintf("%+v", actual)
				vexpect := fmt.Sprintf("%+v", expect)
				if vactual != vexpect {
					e.t.T.Fatal(fmt.Sprintf("Not equal: \n"+
						"expect: %s\n"+
						"actual: %s", vexpect, vactual))
				}
			} else {
				e.t.T.Fatal("Observable emitted multiple values, expected exactly one")
			}
		case err != nil:
			e.t.T.Fatalf("No error expected, but got: %s", err.Error())
		default:
			if initial {
				e.t.T.Fatal("Observable completed without emitting a value, expected exactly one")
			}
		}
	}
	return observer
}

type intobserver = func(next int, err error, done bool)

func (e expect) IntValue(expect int) intobserver {
	initial := true
	observer := func(actual int, err error, done bool) {
		switch {
		case !done:
			if initial {
				initial = false
				if actual != expect {
					e.t.T.Fatal(fmt.Sprintf("Not equal: \n"+
						"expect: %d\n"+
						"actual: %d", expect, actual))
				}
			} else {
				e.t.T.Fatal("Observable emitted multiple values, expected exactly one")
			}
		case err != nil:
			e.t.T.Fatalf("No error expected, but got: %s", err.Error())
		default:
			if initial {
				e.t.T.Fatal("Observable completed without emitting a value, expected exactly one")
			}
		}
	}
	return observer
}
