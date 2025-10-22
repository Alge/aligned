package main

import (
	"fmt"
	"io"
)

func printVersion(stdout io.Writer) int {
	fmt.Fprintf(stdout, "align version %s\n", version)
	return 0
}