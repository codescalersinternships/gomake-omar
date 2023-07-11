package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/codescalersinternships/gomake-omar/internal"
)

func main() {
	filepath := flag.String("f", "", "Specify the filepath")
	target := flag.String("t", "", "Specify the target")
	flag.Parse()

	gomake := internal.NewGomake()
	err := gomake.RunGoMake(*filepath, *target)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		if errors.Is(err, internal.ErrInvalidMakefileFormat) {
			os.Exit(1)
		} else if errors.Is(err, internal.ErrCouldntExecuteCommand) {
			os.Exit(2)
		} else {
			os.Exit(5)
		}
	}
}
