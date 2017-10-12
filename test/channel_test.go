package test

//jig:file {{.package}}.go

import (
	"fmt"
	"testing"
	"time"

	// "github.com/stretchr/testify/assert"

	_ "github.com/reactivego/rx"
)

const MAXPAR = 3
const BUFSIZE = 10

func TestReplayChan(t *testing.T) {
	ch := NewReplayChanInt(BUFSIZE, 0)

	drain := func(ep *ReplayEndpointInt, count *int32, wait chan struct{}) {
		ep.Range(func(value int, err error, closed bool) bool {
			switch {
			case !closed:
				(*count)++
				// println(value)
			case err != nil:
				// println(err.Error())
			default:
				// println("complete")
			}
			return true
		})
		close(wait)
	}

	var counters []*int32
	var waits []chan struct{}

	for i := 0; i < MAXPAR; i++ {
		count := int32(0)
		wait := make(chan struct{})
		go drain(ch.NewEndpoint(), &count, wait)
		counters = append(counters, &count)
		waits = append(waits, wait)
	}

	count := int32(0)
	end := time.Now().Add(time.Second)
	for i := 0; time.Now().Before(end); i++ {
		ch.Send(i)
		count++
	}
	ch.Close(nil)

	for p := 0; p < MAXPAR; p++ {
		<-waits[p]
	}

	tps := float64(1000000000) / float64(count)
	fmt.Printf("<: %.0f ns/send\n", tps)
	for p := 0; p < MAXPAR; p++ {
		tps := float64(1000000000) / float64(*counters[p])
		fmt.Printf("%d: %.0f ns/send\n", p, tps)
	}
}

func TestNativeChan(t *testing.T) {
	drain := func(ch chan interface{}, count *int32, wait chan struct{}) {
		for range ch {
			(*count)++
		}
		close(wait)
	}

	var counters []*int32
	var channels []chan interface{}
	var waits []chan struct{}

	for i := 0; i < MAXPAR; i++ {
		ch := make(chan interface{}, BUFSIZE)
		count := int32(0)
		wait := make(chan struct{})
		go drain(ch, &count, wait)
		counters = append(counters, &count)
		channels = append(channels, ch)
		waits = append(waits, wait)
	}

	count := int32(0)
	end := time.Now().Add(time.Second)
	for i := 0; time.Now().Before(end); i++ {
		for p := 0; p < MAXPAR; p++ {
			channels[p] <- i
		}
		count++
	}
	for p := 0; p < MAXPAR; p++ {
		close(channels[p])
	}

	for p := 0; p < MAXPAR; p++ {
		<-waits[p]
	}

	tps := float64(1000000000) / float64(count)
	fmt.Printf("<: %.0f ns/send\n", tps)
	for p := 0; p < MAXPAR; p++ {
		tps := float64(1000000000) / float64(*counters[p])
		fmt.Printf("%d: %.0f ns/send\n", p, tps)
	}

}
