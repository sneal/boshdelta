package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("NAME:")
	fmt.Println("  boshdelta - a command line tool to compare two different BOSH releases")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  boshdelta release1 release2")
	fmt.Println()
}
