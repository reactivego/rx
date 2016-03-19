// +build special

package rxgo

import (
	"errors"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"rxgo/observable"
)

func TestNilRange(t *testing.T) {
	var items []observable.Unsubscriber
	assert.NotPanics(t, func() {
		for _, v := range items {
			_ = v
		}
	}, "Calling range on nil slice should not panic")
}

func TestConcat(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5}
	c := []int{6, 7}
	oa := observable.FromIntArray(a)
	ob := observable.FromIntArray(b)
	oc := observable.FromIntArray(c)
	s := oa.Concat(ob).Concat(oc).ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
	s = oa.Concat(ob, oc).ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
}

func TestMerge(t *testing.T) {
	sa := observable.CreateInt(func(observer observable.IntSubscriber) {
		time.Sleep(10 * time.Millisecond)
		observer.Next(1)
		time.Sleep(10 * time.Millisecond)
		observer.Next(3)
		observer.Complete()
	})
	sb := observable.CreateInt(func(observer observable.IntSubscriber) {
		time.Sleep(5 * time.Millisecond)
		observer.Next(0)
		time.Sleep(10 * time.Millisecond)
		observer.Next(2)
		observer.Complete()
	})
	a := sa.Merge(sb).ToArray()
	assert.Equal(t, []int{0, 1, 2, 3}, a)
}

func TestMergeDelayError(t *testing.T) {
	sa := observable.CreateInt(func(observer observable.IntSubscriber) {
		time.Sleep(10 * time.Millisecond)
		observer.Next(1)
		observer.Error(errors.New("error"))
	})
	sb := observable.CreateInt(func(observer observable.IntSubscriber) {
		time.Sleep(5 * time.Millisecond)
		observer.Next(0)
		time.Sleep(10 * time.Millisecond)
		observer.Next(2)
		observer.Complete()
	})
	a := sa.MergeDelayError(sb).ToArray()
	assert.Equal(t, []int{0, 1, 2}, a)
}

func TestRecover(t *testing.T) {
	// Create three different observables.
	// They are all passive until someone subscribes to them...
	o123 := observable.FromInts(1, 2, 3)
	o45 := observable.FromInts(4, 5)
	oThrowError := observable.ThrowInt(errors.New("error"))

	a, err := o123.Concat(oThrowError).Catch(o45).ToArrayWithError()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}

func TestRetry(t *testing.T) {
	errored := false
	a := observable.CreateInt(func(observer observable.IntSubscriber) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		if errored {
			observer.Complete()
		} else {
			observer.Error(errors.New("error"))
			errored = true
		}
	})
	b := a.Retry().ToArray()
	assert.Equal(t, []int{1, 2, 3, 1, 2, 3}, b)
	assert.True(t, errored)
}
