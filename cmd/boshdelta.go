package main

import (
	"fmt"
	"os"

	"github.com/sneal/bosh-delta/boshdelta"
)

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(fmt.Sprintf("Comparing %s to %s", os.Args[1], os.Args[2]))
	delta, err := boshdelta.Compare(os.Args[1], os.Args[2])
	if err != nil {
		fail(err)
	}
	printDelta(delta)
}

func printDelta(d *boshdelta.Delta) {
	fmt.Println(fmt.Sprintf("Found %d new properties", len(d.DeltaProperties)))
	fmt.Println()
	for di := range d.DeltaProperties {
		p := d.DeltaProperties[di]
		fmt.Println("-------------------------------------------------------------------------------")
		fmt.Println(fmt.Sprintf("%s", p.Name))
		fmt.Println("-------------------------------------------------------------------------------")
		fmt.Println(fmt.Sprintf("%s", p.Description))
		fmt.Println()
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

func fail(err error) {
	fmt.Println("BOSH release comparision failed!")
	fmt.Println(err.Error())
	os.Exit(1)
}
