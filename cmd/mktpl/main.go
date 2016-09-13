package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	f, err := os.Open("reactivex.go")
	if err != nil {
		println(err.Error())
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		text := scanner.Text()

		if strings.Index(text, "//rxgo:") != -1 {

		} else {
			fmt.Println(text) // Println will add back the final '\n'
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading file %q:%v\n", "reactivex.go", err)
	}

}
