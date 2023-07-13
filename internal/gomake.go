package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrInvalidMakefileFormat = errors.New("invalid makefile format")
	ErrNoTarget              = errors.New("target must be specified")
	ErrCyclicDependency      = errors.New("there is a cyclic dependency")
	ErrCouldntExecuteCommand = errors.New("could not execute command")
	ErrDependencyNotFound    = errors.New("dependency rule is not found")
)

type target struct {
	dependencies []string
	commands     []command
}

type Make struct {
	targets map[string]target
}

// NewGomake is a factory of gomake struct
func NewGomake() Make {
	return Make{
		targets: map[string]target{},
	}
}

func (gomake *Make) getDependencyGraph() map[string][]string {
	g := map[string][]string{}

	for targetName, target := range gomake.targets {
		g[targetName] = target.dependencies
	}

	return g
}

func (mk *Make) addCommandLine(targetName, commandLine string) error {
	if targetName == "" { // try to add command before writing a target
		return ErrInvalidMakefileFormat
	}

	cmd := strings.TrimSpace(commandLine)
	if cmd == "" || cmd == "@" {
		return nil
	}

	// take a copy
	entry := mk.targets[targetName]

	if cmd[0] == '@' {
		entry.commands = append(entry.commands, command{cmd[1:], true})
	} else {
		entry.commands = append(entry.commands, command{cmd, false})
	}

	// update
	mk.targets[targetName] = entry

	return nil
}

func (mk *Make) addTargetLine(targetLine string) (string, error) {
	lineParts := strings.SplitN(targetLine, ":", 2)
	targetName := strings.TrimSpace(lineParts[0])
	dependencies := strings.Fields(lineParts[1])

	if targetName == "" {
		return "", ErrNoTarget
	}

	// take a copy
	entry := mk.targets[targetName]

	if entry.dependencies != nil {
		// target name appeared before
		fmt.Printf("warning: overriding recipe for target %q\n", targetName)
		fmt.Printf("warning: ignoring old recipe for target %q\n", targetName)

		// remove last added commands
		entry.commands = []command{}
		// append at the front to keep dependency order as makefile doing
		entry.dependencies = append(dependencies, entry.dependencies...)
	} else {
		entry.dependencies = dependencies
	}

	// update
	mk.targets[targetName] = entry

	return targetName, nil
}

func (mk *Make) readData(r io.Reader) error {
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanLines)

	currentTargetName := ""
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if strings.HasPrefix(line, "\t") {
			// this is action

			if err := mk.addCommandLine(currentTargetName, line); err != nil {
				return err
			}
			continue
		}

		if strings.Contains(line, ":") {
			// this is rule
			target, err := mk.addTargetLine(line)
			if err != nil {
				return err
			}

			currentTargetName = target
			continue
		}

		if strings.TrimSpace(line) != "" {
			return fmt.Errorf("%w, at line %q", ErrInvalidMakefileFormat, line)
		}
	}

	return nil
}

// Build analyze given data and prepare it to execute
func (mk *Make) Build(r io.Reader) error {
	err := mk.readData(r)
	if err != nil {
		return err
	}

	// checkCyclicDependency
	graph := newGraph(mk.getDependencyGraph())
	cycle := graph.getCycle()

	if len(cycle) != 0 {
		cycleStr := fmt.Sprintf("%v -> %v", strings.Join(cycle, " -> "), cycle[0])
		return fmt.Errorf("%w, cycle: %q", ErrCyclicDependency, cycleStr)
	}

	return nil
}

// Run executes target
func (mk *Make) Run(targetName string) error {
	graph := newGraph(mk.getDependencyGraph())
	orderedDep := graph.getOrderedDependencies(targetName)

	for _, dep := range orderedDep {
		if _, ok := mk.targets[targetName]; !ok {
			return fmt.Errorf("%w, target %q dependant on %q", ErrDependencyNotFound, targetName, dep)
		}

		for _, command := range mk.targets[dep].commands {
			if err := command.execute(); err != nil {
				return err
			}
		}
	}

	return nil
}
