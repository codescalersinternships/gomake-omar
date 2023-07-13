package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/codescalersinternships/gomake-omar/internal"
)

func checkErr(err error) bool {
	if err == nil {
		return false
	}

	fmt.Fprintln(os.Stderr, "Error:", err)
	switch {
	case errors.Is(err, internal.ErrInvalidMakefileFormat):
		os.Exit(1)
	case errors.Is(err, internal.ErrCouldntExecuteCommand):
		os.Exit(2)
	default:
		os.Exit(5)
	}

	return true
}

func main() {
	var filepath, targetName string
	flag.StringVar(&filepath, "f", "./Makefile", "Specify the filepath")
	flag.StringVar(&targetName, "t", "", "Specify the target")
	flag.Parse()

	if targetName == "" {
		checkErr(internal.ErrNoTarget)
		return
	}

	gomake := internal.NewGomake()

	f, err := os.Open(filepath)
	if checkErr(err) {
		return
	}
	defer f.Close()

	err = gomake.Build(f)
	if checkErr(err) {
		return
	}

	err = gomake.Run(targetName)
	if checkErr(err) {
		return
	}
}
