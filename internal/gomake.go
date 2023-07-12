package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrInvalidMakefileFormat = errors.New("invalid makefile format")
	ErrTargetMustBeSpecified = errors.New("target must be specified")
	ErrCyclicDependency      = errors.New("there is a cyclic dependency")
	ErrCouldntExecuteCommand = errors.New("couldn't execute command")
	ErrDependencyNotFound    = errors.New("dependency rule is not found")
)

type target string
type parent string

type gomake struct {
	dependencyGraph map[target][]parent
	actionExecuter  executer
}

// NewGomake is a factory of gomake struct
func NewGomake() gomake {
	return gomake{
		dependencyGraph: map[target][]parent{},
		actionExecuter:  newExecuter(),
	}
}

func (gomake *gomake) addActionLine(target target, actionLine string) error {
	if target == "" {
		return ErrInvalidMakefileFormat
	}

	action := action(strings.TrimSpace(actionLine))
	if action != "" {
		gomake.actionExecuter.setAction(target, action)
	}

	return nil
}

func (gomake *gomake) addTargetLine(targetLine string) (target, error) {
	lineParts := strings.SplitN(targetLine, ":", 2)
	target := target(strings.TrimSpace(lineParts[0]))
	parents := strings.Split(lineParts[1], " ")

	if target == "" {
		return "", ErrTargetMustBeSpecified
	}

	filteredParents := []parent{}
	for _, p := range parents {
		parent := parent(strings.TrimSpace(p))
		if parent != "" {
			filteredParents = append(filteredParents, parent)
		}
	}

	if gomake.dependencyGraph[target] != nil {
		// target appeared before
		fmt.Printf("warning: overriding recipe for target %q\n", target)
		fmt.Printf("warning: ignoring old recipe for target %q\n", target)

		gomake.actionExecuter.removeActions(target)

		// append at the front to keep dependency order as makefile doing
		gomake.dependencyGraph[target] = append(filteredParents, gomake.dependencyGraph[target]...)
	} else {
		gomake.dependencyGraph[target] = filteredParents
	}

	return target, nil
}

func (gomake *gomake) loadData(f io.Reader) error {
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	currentTargetName := target("")
	for fileScanner.Scan() {
		line := fileScanner.Text()

		if strings.HasPrefix(line, "\t") {
			// this is action
			err := gomake.addActionLine(currentTargetName, line)
			if err != nil {
				return err
			}

			continue
		}

		if strings.Contains(line, ":") {
			// this is rule
			target, err := gomake.addTargetLine(line)
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

func getCycleString(cycle []target) string {
	castedCycle := make([]string, len(cycle))
	for i, v := range cycle {
		castedCycle[i] = string(v)
	}

	return fmt.Sprintf("%v -> %v", strings.Join(castedCycle, " -> "), castedCycle[0])
}

// RunGoMake checks if there's cyclic dependency within makefile then run target
func (gomake *gomake) RunGoMake(filePath, targetToExecute string) error {
	if targetToExecute == "" {
		return fmt.Errorf("%w, use -t to specify it", ErrTargetMustBeSpecified)
	}

	if filePath == "" {
		filePath = "Makefile"
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gomake.loadData(file)
	if err != nil {
		return err
	}

	graph := newGraph(gomake.dependencyGraph)
	cycle := graph.getCycle()

	if len(cycle) != 0 {
		return fmt.Errorf("%w, cycle: %q", ErrCyclicDependency, getCycleString(cycle))
	}

	dependencies := graph.getDependency(target(targetToExecute))
	return gomake.actionExecuter.execute(dependencies)
}
