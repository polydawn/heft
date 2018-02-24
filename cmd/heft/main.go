package main

import (
	"os"

	"go.polydawn.net/heft/interpret"
)

func main() {
	interpret.ExecFile(os.Args[1])
}
