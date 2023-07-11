package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// define error
var (
	ErrInvalidMakefileFormat = errors.New("invalid makefile format") // 1
	ErrTargetMustBeSpecified = errors.New("target must be specified") // 5
	ErrCyclicDependency      = errors.New("there is a cyclic dependency") // 5
	ErrCouldntExecuteCommand = errors.New("couldn't execute command") // 2
	ErrDependencyNotFound    = errors.New("dependency rule is not found") // 5
)

type target string
type parent string

type gomake struct {
	dependencyGraph map[target][]parent
	executer        executer
}

func NewGomake() gomake {
	return gomake{
		dependencyGraph: map[target][]parent{},
		executer:        NewExecuter(),
	}
}

func (gomake *gomake) addActionLine(target target, actionLine string) error {
	if target == "" {
		return ErrInvalidMakefileFormat
	}

	action := action(strings.TrimSpace(actionLine))
	if action != "" {
		gomake.executer.setAction(target, action)
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

	if gomake.dependencyGraph[target] == nil {
		gomake.dependencyGraph[target] = []parent{}
	}

	for _, p := range parents {
		if p == "" {
			continue
		}

		parent := parent(strings.TrimSpace(p))
		if parent != "" {
			gomake.dependencyGraph[target] = append(gomake.dependencyGraph[target], parent)
		}
	}

	return target, nil
}

func (gomake *gomake) loadFromFile(f io.Reader) error {
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
		} else if strings.Contains(line, ":") {
			// this is rule
			target, err := gomake.addTargetLine(line)
			if err != nil {
				return err
			}

			currentTargetName = target
		} else if strings.TrimSpace(line) != "" {
			return fmt.Errorf("%w, at line %q", ErrInvalidMakefileFormat, line)
		}
	}
	return nil
}

// RunGoMake checks if there's cyclic dependency within makefile
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

	err = gomake.loadFromFile(file)
	if err != nil {
		return err
	}

	graph := NewGraph(gomake.dependencyGraph)
	cycle := graph.getCycle()

	if len(cycle) != 0 {
		reversedCycle := []string{}
		for i := len(cycle) - 1; i >= 0; i-- {
			reversedCycle = append(reversedCycle, string(cycle[i]))
		}

		cycleStr := strings.Join(reversedCycle, "->")
		return fmt.Errorf("%w, cycle: %q", ErrCyclicDependency, cycleStr)
	}

	dependencies := graph.getDependency(target(targetToExecute))
	return gomake.executer.execute(dependencies)
}
