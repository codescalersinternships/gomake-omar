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
	name         string
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

func (mk *Make) parseTargetLine(targetLine string) (target, error) {
	lineParts := strings.SplitN(targetLine, ":", 2)
	targetName := strings.TrimSpace(lineParts[0])
	dependencies := strings.Fields(lineParts[1])

	if targetName == "" {
		return target{}, ErrNoTarget
	}

	t := target{name: targetName, dependencies: dependencies}
	return t, nil
}

func (mk *Make) parseCommandLine(commandLine string) command {
	// 'commandLine' definitely has at least one non-whitespace character

	commandLine = strings.TrimSpace(commandLine)
	if commandLine[0] == '@' {
		return command{cmd: commandLine[1:], suppressed: true}
	}
	return command{cmd: commandLine, suppressed: false}
}

func (mk *Make) setTarget(t target) {
	if _, ok := mk.targets[t.name]; !ok {
		// this target name haven't been added before
		mk.targets[t.name] = t
		return
	}

	fmt.Printf("warning: overriding recipe for target %q\n", t.name)
	fmt.Printf("warning: ignoring old recipe for target %q\n", t.name)

	// take a copy
	entry := mk.targets[t.name]

	// remove last added commands
	entry.commands = []command{}
	// append at the front to keep dependency order as makefile doing
	entry.dependencies = append(t.dependencies, entry.dependencies...)

	// update
	mk.targets[t.name] = entry
}

func (mk *Make) setCommand(t target, c command) error {
	if t.name == "" {
		// try to add command before initialize target
		return ErrInvalidMakefileFormat
	}

	// take a copy
	entry := mk.targets[t.name]

	entry.commands = append(entry.commands, c)

	// update
	mk.targets[t.name] = entry

	return nil
}

func (mk *Make) readData(r io.Reader) error {
	fileScanner := bufio.NewScanner(r)
	fileScanner.Split(bufio.ScanLines)

	currentTarget := target{}
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if strings.TrimSpace(line) == "" {
			// empty line
			continue
		}

		if strings.HasPrefix(line, "\t") {
			// this is command line

			c := mk.parseCommandLine(line)
			if err := mk.setCommand(currentTarget, c); err != nil {
				return err
			}
			continue
		}

		if strings.Contains(line, ":") {
			// this is target line

			t, err := mk.parseTargetLine(line)
			if err != nil {
				return err
			}
			mk.setTarget(t)
			currentTarget = t
			continue
		}

		return fmt.Errorf("%w, at line %q", ErrInvalidMakefileFormat, line)
	}

	return nil
}

// Build analyze given data and prepare it to execute
func (mk *Make) Build(r io.Reader) error {
	if err := mk.readData(r); err != nil {
		return err
	}

	// check cyclic dependency
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
	if _, ok := mk.targets[targetName]; !ok {
		return fmt.Errorf("%w, target %q", ErrDependencyNotFound, targetName)
	}

	graph := newGraph(mk.getDependencyGraph())
	orderedDep := graph.getOrderedDependencies(targetName)

	for _, dep := range orderedDep {
		if _, ok := mk.targets[dep]; !ok {
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
