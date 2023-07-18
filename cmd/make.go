package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/codescalersinternships/gomake-omar/internal"
)

func exitStatusCode(err error) int {
	if err == nil {
		return 0
	}

	switch {
	case errors.Is(err, internal.ErrInvalidMakefileFormat):
		return 1
	case errors.Is(err, internal.ErrCouldntExecuteCommand):
		return 2
	default:
		return 5
	}
}

func main() {
	var filepath, targetName string
	flag.StringVar(&filepath, "f", "./Makefile", "Specify the filepath")
	flag.StringVar(&targetName, "t", "", "Specify the target")
	flag.Parse()

	if targetName == "" {
		fmt.Fprintln(os.Stderr, "Error:", internal.ErrNoTarget)
		os.Exit(exitStatusCode(internal.ErrNoTarget))
	}

	gomake := internal.NewGomake()

	makefile, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(exitStatusCode(err))
	}
	defer makefile.Close()

	if err = gomake.Build(makefile); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(exitStatusCode(err))
	}

	if err = gomake.Run(targetName); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(exitStatusCode(err))
	}
}
