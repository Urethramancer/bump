package main

import (
	"fmt"
	"os"
)

func pr(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}

func showHelp() {
	pr("%s increases the version number of a semver-conformant git repository, based on the most recent annotated tag.\n", os.Args[0])
	pr("Usage:")
	pr("%s <segment>", os.Args[0])
	pr("Possible keywords are 'major', 'minor' and 'patch' with the shortcuts 'ma', 'mi' or 'p'.")
}
