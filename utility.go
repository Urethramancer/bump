package main

import (
	"fmt"
)

func pr(f string, v ...any) {
	fmt.Printf(f+"\n", v...)
}
