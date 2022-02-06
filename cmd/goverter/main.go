package main

import (
	"flag"
	"fmt"
	"os"

	goverter "github.com/FanFani4/goverter"
)

func main() {
	packageName := flag.String("packageName", "generated", "")
	output := flag.String("output", "./generated/generated.go", "")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "expected one argument")
		return
	}
	pattern := args[0]

	err := goverter.GenerateConverterFile(*output, goverter.GenerateConfig{
		PackageName: *packageName,
		ScanDir:     pattern,
	})

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
}
