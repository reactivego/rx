package test

import (
    "fmt"
    "testing"
)

/*
BDD provides support for writing Behavior Driven Development specs directly in
Go. It is a struct embedding a *testing.T and therefore has the same methods as
a testing.T type.

See also https://inviqa.com/blog/bdd-guide
*/
type BDD struct {
    *testing.T
}

// Describ(e) is Describe from a BDD testing language spec.
// The name is intentionally missing a character at the end so it can be used as follows:
//
//  func TestObservable_ConcatAll(e *testing.T) {
//      Describ(e, "emits", func(t BDD) {
//          Contex(t, "with a fast source", func(t BDD) {
//              Contex(t, "and a slow concat", func(t BDD) {
//                  I(t, "should complete before the first emit to observer", func(t BDD) {
//                      //
//                  })
//              })
//          })
//      })
//  }
// 
func Describ(e *testing.T, name string, f func(BDD)) bool {
    e.Helper()
    return e.Run(name, func(t *testing.T) {
        f(BDD{t})
    })
}

// Contex(t) represents Context from a BDD DSL
func Contex(t BDD, name string, f func(BDD)) bool {
    t.Helper()
    return t.Run(name, func(t *testing.T) {
        f(BDD{t})
    })
}

// Whe(n) represents When from a BDD DSL
func Whe(n BDD, name string, f func(BDD)) bool {
    n.Helper()
    return n.Run(name, func(t *testing.T) {
        f(BDD{t})
    })
}

// I(t) represents It from a BDD DSL
func I(t BDD, name string, f func(BDD)) bool {
    t.Helper()
    return t.Run(name, func(t *testing.T) {
        f(BDD{t})
    })
}

func Asser(t BDD, args ...interface{}) {
    t.Helper()
    if len(args) > 0 {
        if assertion, ok := args[0].(Assertion); ok {
            if !assertion.cond {
                t.Fatal(append([]interface{}{assertion.mesg}, args[1:]...)...)
            }
            return
        }
        if cond, ok := args[0].(bool); ok {
            if !cond {
                t.Fatal(args[1:])
            }
            return
        }
    }
    t.Fatal("Assert error: no args or invalid args")
}

type Assertion struct {
    cond bool
    mesg string
}

func Must(cond bool, args ...interface{}) Assertion {
    if !cond {
        return Assertion{mesg: fmt.Sprint(append([]interface{}{"Must have "}, args...)...)}
    }
    return Assertion{true, "OK"}
}

func Not(cond bool, args ...interface{}) Assertion {
    if cond {
        return Assertion{mesg: fmt.Sprint(append([]interface{}{"Unexpected "}, args...)...)}
    }
    return Assertion{true, "OK"}
}

func NoError(err error, args ...interface{}) Assertion {
    if err == nil {
        return Assertion{true, "OK"}
    }
    return Assertion{mesg: fmt.Sprintf("Unexpected error: %s %s", err.Error(), fmt.Sprint(args...))}
}

func Equal(actual, expect interface{}, args ...interface{}) Assertion {
    tactual := fmt.Sprintf("%T", actual)
    texpect := fmt.Sprintf("%T", expect)
    if tactual != texpect {
        return Assertion{mesg: fmt.Sprintf("Not equal: \n"+
            "expect: %s\n"+
            "actual: %s %s", texpect, tactual, fmt.Sprint(args...))}
    }
    vactual := fmt.Sprintf("%+v", actual)
    vexpect := fmt.Sprintf("%+v", expect)
    if vactual != vexpect {
        return Assertion{mesg: fmt.Sprintf("Not equal: \n"+
            "expect: %s\n"+
            "actual: %s %s", vexpect, vactual, fmt.Sprint(args...))}
    }
    return Assertion{true, "OK"}
}
