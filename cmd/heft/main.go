package main

import (
	"fmt"
	"os"

	"github.com/google/skylark"
)

func main() {
	execfile(os.Args[1])
}

// Makes a new globals dict with our favorite custom bits in it.
func newGlobals() skylark.StringDict {
	return skylark.StringDict{}
}

func execfile(filename string) {
	thread := &skylark.Thread{}
	globals := newGlobals()
	if err := skylark.ExecFile(thread, filename, nil, globals); err != nil {
		fmt.Fprintf(os.Stderr, "larking: %s", err)
		os.Exit(4)
	}
}
