package Publish

import (
	"fmt"
)

func Example_publish() {
	pub := FromInts(1, 2).Publish()

	gr := NewGoroutine()

	sub0 := pub.SubscribeOn(gr).
		MapString(func(v int) string {
			return fmt.Sprintf("value is %d", v)
		}).
		Subscribe(func(next string, err error, done bool) {
			if !done {
				fmt.Println(next)
			}
		})

	sub1 := pub.SubscribeOn(gr).
		MapBool(func(v int) bool {
			return v > 1
		}).
		Subscribe(func(next bool, err error, done bool) {
			if !done {
				fmt.Println("value is ", next)
			}
		})

	sub2 := pub.SubscribeOn(gr).
		Subscribe(func(next int, err error, done bool) {
			if !done {
				fmt.Println(next)
			}
		})

	pub.Connect()

	sub0.Wait()
	sub1.Wait()
	sub2.Wait()

	//Output:
	// 1
	// value is 1
	// value is 2
	// value is  false
	// value is  true
	// 2
}
