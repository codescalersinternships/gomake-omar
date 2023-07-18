package internal

import (
	"fmt"
	"os/exec"
)

type command struct {
	cmdName    string
	cmdArgs    []string
	suppressed bool // in case of @ prefixing the command
}

func (c *command) execute() error {
	if !c.suppressed {
		fmt.Printf("%s %v\n", c.cmdName, c.cmdArgs)
	}

	cmd := exec.Command(c.cmdName, c.cmdArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("%w, error message: %q", ErrCouldntExecuteCommand, err)
	}

	fmt.Print(string(output))

	return nil
}
