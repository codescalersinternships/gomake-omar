package internal

import (
	"errors"
	"io"
	"os"
	"reflect"
	"sort"
	"testing"
)

var makeFileSample1 = `
build:
	@echo 'executing build'
	@echo 'cmd2'
test:
		echo 'executing test'

publish : test gendocs
	echo 'executing publish'

 gendocs  : build
	echo 'executing gendocs'
`

func depGraphSample1() map[target][]parent {
	return map[target][]parent{
		"build":   {},
		"test":    {},
		"publish": {"test", "gendocs"},
		"gendocs": {"build"},
	}
}

func actionGraphSample1() map[target][]action {
	return map[target][]action{
		"build":   {"@echo 'executing build'", "@echo 'cmd2'"},
		"test":    {"echo 'executing test'"},
		"publish": {"echo 'executing publish'"},
		"gendocs": {"echo 'executing gendocs'"},
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Run("valid makefile format", func(t *testing.T) {
		f, err := os.CreateTemp("", "test_file")
		assertErr(t, err, nil)
		defer os.Remove(f.Name())

		_, err = f.WriteString(makeFileSample1)
		assertErr(t, err, nil)
		_, err = f.Seek(0, io.SeekStart)
		assertErr(t, err, nil)

		gomake := NewGomake()
		err = gomake.loadFromFile(f)
		assertErr(t, err, nil)

		assertEqualDepGraphs(t, gomake.dependencyGraph, depGraphSample1())
	})

	invalidTestCases := []struct {
		name     string
		makefile string
		err      error
	}{
		{
			name:     "invalid format: no target exist",
			makefile: "\t@echo 'too sad'",
			err:      ErrInvalidMakefileFormat,
		},
		{
			name:     "invalid format: invalid target line",
			makefile: "build dep",
			err:      ErrInvalidMakefileFormat,
		}, {
			name:     "invalid format: invalid target line",
			makefile: ":dep",
			err:      ErrTargetMustBeSpecified,
		}, {
			name:     "invalid format: invalid action line",
			makefile: "build: dep\n@echo 'too sad'",
			err:      ErrInvalidMakefileFormat,
		},
	}

	for _, tc := range invalidTestCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.CreateTemp("", "test_file")
			assertErr(t, err, nil)
			defer os.Remove(f.Name())

			_, err = f.WriteString(tc.makefile)
			assertErr(t, err, nil)
			_, err = f.Seek(0, io.SeekStart)
			assertErr(t, err, nil)

			gomake := NewGomake()
			err = gomake.loadFromFile(f)
			assertErr(t, err, tc.err)
		})
	}
}

func TestRunGoMake(t *testing.T) {
	testCases := []struct {
		name       string
		target     target
		makefile   string
		err        error
		depWant    map[target][]parent
		actionWant map[target][]action
	}{
		{
			name:       "given file path",
			target:     "publish",
			makefile:   makeFileSample1,
			err:        nil,
			depWant:    depGraphSample1(),
			actionWant: actionGraphSample1(),
		},
		{
			name:       "no file path given",
			target:     "publish",
			makefile:   makeFileSample1,
			err:        nil,
			depWant:    depGraphSample1(),
			actionWant: actionGraphSample1(),
		}, {
			name:     "no target given",
			target:   "",
			makefile: makeFileSample1,
			err:      ErrTargetMustBeSpecified,
		}, {
			name:   "cyclic dependency",
			target: "a",
			makefile: `
a: b
b: c
c: a
`,
			err: ErrCyclicDependency,
		}, {
			name:   "invalid command",
			target: "a",
			makefile: `
a: 
	@invalid command
`,
			err: ErrCouldntExecuteCommand,
		}, {
			name:   "target repeats",
			target: "a",
			makefile: `
a: b b
	echo 'a'
a: c
	echo 'newa'

c:
	echo 'c'
b:
	echo 'b'
`,
			err: nil,
			depWant: map[target][]parent{
				"a": {"c", "b", "b"},
				"c": {},
				"b": {},
			},
			actionWant: map[target][]action{
				"a": {"echo 'newa'"},
				"c": {"echo 'c'"},
				"b": {"echo 'b'"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.CreateTemp("", "Makefile")
			assertErr(t, err, nil)
			defer os.Remove(f.Name())

			_, err = f.WriteString(tc.makefile)
			assertErr(t, err, nil)
			_, err = f.Seek(0, io.SeekStart)
			assertErr(t, err, nil)

			gomake := NewGomake()
			err = gomake.RunGoMake(f.Name(), string(tc.target))
			assertErr(t, err, tc.err)

			if tc.err == nil {
				assertEqualDepGraphs(t, gomake.dependencyGraph, tc.depWant)
				assertEqualActionGraphs(t, gomake.actionExecuter.targetActions, tc.actionWant)
			}
		})
	}
}

func assertErr(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		if got == nil {
			t.Errorf("got nil want %q", want.Error())
		} else if want == nil {
			t.Errorf("got %q want nil", got.Error())
		} else {
			t.Errorf("got %q want %q", got.Error(), want.Error())
		}
	}
}

func assertEqualDepGraphs(t testing.TB, got, want map[target][]parent) {
	t.Helper()
	for k, v := range got {
		if want[k] == nil {
			t.Errorf("got %v want %v", got, want)
			t.Fatalf("key %q exists in got and not exist in want", k)
		}

		sort.Slice(want[k], func(i, j int) bool {
			return want[k][i] < want[k][j]
		})
		sort.Slice(v, func(i, j int) bool {
			return v[i] < v[j]
		})

		if !reflect.DeepEqual(v, want[k]) {
			t.Errorf("in key %q got value %v want value %v", k, v, want[k])
		}
	}

	if len(got) != len(want) {
		t.Errorf("there is keys in want not exist in got\ngot %v want %v", got, want)
	}
}

func assertEqualActionGraphs(t testing.TB, got, want map[target][]action) {
	t.Helper()
	for k, v := range got {
		if want[k] == nil {
			t.Errorf("got %v want %v", got, want)
			t.Fatalf("key %q exists in got and not exist in want", k)
		}

		sort.Slice(want[k], func(i, j int) bool {
			return want[k][i] < want[k][j]
		})
		sort.Slice(v, func(i, j int) bool {
			return v[i] < v[j]
		})

		if !reflect.DeepEqual(v, want[k]) {
			t.Errorf("in key %q got value %v want value %v", k, v, want[k])
		}
	}

	if len(got) != len(want) {
		t.Errorf("there is keys in want not exist in got\ngot %v want %v", got, want)
	}
}
