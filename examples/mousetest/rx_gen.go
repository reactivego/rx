// +build generate

package main

//go:generate rxgo

//rx:package observable

//rx:import
//  "rxgo/examples/mousetest/mouse"

//rx:type
//  observable int
//    flatmap[int]float64
//    map[int]float64
//    flatmap[int]*mouse.Move
//    map[int]*mouse.Move

//  observable float64
//    flatmap[float64]int
//    map[float64]int
//    flatmap[float64]*mouse.Move
//    map[float64]*mouse.Move

//  observable string

//  observable bool

//  observable *mouse.Move

//rx:end
