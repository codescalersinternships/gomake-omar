package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

type command struct {
	cmd        string
	suppressed bool // in case of @ prefixing the command
}

func (c *command) execute() error {
	if c.suppressed {
		fmt.Println(c.cmd)
	}

	parts := strings.Fields(string(c.cmd))
	cmdName := parts[0]
	cmdArgs := parts[1:]

	cmd := exec.Command(cmdName, cmdArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return ErrCouldntExecuteCommand
	}

	fmt.Print(string(output))

	return nil
}
