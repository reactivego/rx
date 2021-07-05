package DistinctUntilChanged

import (
	"image/color"

	_ "github.com/reactivego/rx"
)

func Example_distinctUntilChanged() {
	FromInt(1, 2, 2, 1, 3).DistinctUntilChanged().Println()
	// Output:
	// 1
	// 2
	// 1
	// 3
}

func Example_distinctUntilChangedRGBA() {
	//jig:type RGBA = color.RGBA
	c := []color.RGBA{
		{R: 1, G: 5, A: 1},
		{R: 2, G: 6, A: 2},
		{R: 2, G: 6, A: 3}, // this one is NOT distinct => skip
		{R: 1, G: 5, A: 4},
		{R: 3, G: 3, A: 5},
	}
	red := func(next color.RGBA) interface{} { return next.R }
	green := func(next color.RGBA) interface{} { return next.G }
	FromRGBA(c...).DistinctUntilChanged(red, green).Println()
	// Output:
	// {1 5 0 1}
	// {2 6 0 2}
	// {1 5 0 4}
	// {3 3 0 5}
}
