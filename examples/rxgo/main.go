package main

import (
	"os"
	"strings"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		println(err.Error())
	} else {
		println("PWD=", dir)
	}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "GO") {
			println(env)

		}
	}
}
