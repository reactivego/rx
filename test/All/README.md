# All

## Overview

[All](http://reactivex.io/documentation/operators/all.html) determines whether all items emitted by an Observable meet some
criteria.

Pass a predicate function to the All operator that accepts an item emitted
by the source Observable and returns a boolean value based on an
evaluation of that item. All returns an ObservableBool that emits a single
boolean value: true if and only if the source Observable terminates
normally and every item emitted by the source Observable evaluated as
true according to the predicate; false if any item emitted by the source
Observable evaluates as false according to the predicate.

![All](http://reactivex.io/documentation/operators/images/all.png)

## Example (All)
Code:
```go
source := CreateInt(func(N NextInt, E Error, C Complete, X Canceled) {
	N(1)
	N(2)
	N(6)
	N(2)
	N(1)
	C()
})

// Print the sequence of numbers
source.Println()
fmt.Println("Source observable completed")

// Setup All to produce true only when all source values are less than 5
predicate := func(i int) bool {
	return i < 5
}

// Observer prints a message describing the next value it observes.
observer := func(next bool, err error, done bool) {
	if !done {
		fmt.Println("All values less than 5?", next)
	}
}

source.All(predicate).Subscribe(observer).Wait()
```
Output:
```
1
2
6
2
1
Source observable completed
All values less than 5? false
```