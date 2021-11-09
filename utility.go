package main

import (
	"fmt"
)

func pr(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}
