package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

type action string

type executer struct {
	targetActions map[target][]action
}

func NewExecuter() executer {
	return executer{
		targetActions: map[target][]action{},
	}
}

func (exe *executer) setAction(target target, action action) {
	exe.targetActions[target] = append(exe.targetActions[target], action)
}

func (exe *executer) execute(targets []target) error {
	for _, target := range targets {
		if exe.targetActions[target] == nil {
			return fmt.Errorf("%w, rule: %q", ErrDependencyNotFound, target)
		}
		actions := exe.targetActions[target]
		for _, action := range actions {
			if action[0] != '@' {
				fmt.Println(action)
			} else {
				action = action[1:]
			}

			parts := strings.Fields(string(action))
			cmdName := parts[0]
			cmdArgs := parts[1:]
			cmd := exec.Command(cmdName, cmdArgs...)
			output, err := cmd.CombinedOutput()
			
			if err != nil {
				return fmt.Errorf("%w, %q", ErrCouldntExecuteCommand, err)
			}
			fmt.Println(string(output))
		}
	}
	return nil
}
