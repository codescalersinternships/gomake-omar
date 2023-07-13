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

func depGraphSample1() map[string][]string {
	return map[string][]string{
		"build":   {},
		"test":    {},
		"publish": {"test", "gendocs"},
		"gendocs": {"build"},
	}
}

func TestReadData(t *testing.T) {
	t.Run("valid makefile format", func(t *testing.T) {
		f, err := os.CreateTemp("", "test_file")
		assertErr(t, err, nil)
		defer os.Remove(f.Name())

		_, err = f.WriteString(makeFileSample1)
		assertErr(t, err, nil)
		_, err = f.Seek(0, io.SeekStart)
		assertErr(t, err, nil)

		gomake := NewGomake()
		err = gomake.readData(f)
		assertErr(t, err, nil)

		assertEqualDepGraphs(t, gomake.getDependencyGraph(), depGraphSample1())
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
			err:      ErrNoTarget,
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
			err = gomake.readData(f)
			assertErr(t, err, tc.err)
		})
	}
}

func TestBuild(t *testing.T) {
	testCases := []struct {
		name     string
		target   string
		makefile string
		err      error
		depWant  map[string][]string
		commands []command
	}{
		{
			name:     "given file path",
			target:   "publish",
			makefile: makeFileSample1,
			err:      nil,
			depWant:  depGraphSample1(),
			commands: []command{{
				cmd:        "echo 'executing publish'",
				suppressed: false,
			}},
		},
		{
			name:     "no file path given",
			target:   "build",
			makefile: makeFileSample1,
			err:      nil,
			depWant:  depGraphSample1(),
			commands: []command{{
				cmd:        "echo 'executing build'",
				suppressed: true,
			}, {
				cmd:        "echo 'cmd2'",
				suppressed: true,
			}},
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
			name:     "invalid file format",
			target:   "a",
			makefile: "invalid format",
			err:      ErrInvalidMakefileFormat,
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
			depWant: map[string][]string{
				"a": {"c", "b", "b"},
				"c": {},
				"b": {},
			},
			commands: []command{
				{cmd: "echo 'newa'", suppressed: false},
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
			err = gomake.Build(f)
			assertErr(t, err, tc.err)

			if tc.err == nil {
				assertEqualDepGraphs(t, gomake.getDependencyGraph(), tc.depWant)

				if !reflect.DeepEqual(gomake.targets[tc.target].commands, tc.commands) {
					t.Errorf("got %v want %v", gomake.targets[tc.target].commands, tc.commands)
				}
			}
		})
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		target   string
		makefile string
		err      error
	}{
		{
			name:     "valid",
			target:   "build",
			makefile: makeFileSample1,
			err:      nil,
		},
		{
			name:     "target is not exist in file",
			target:   "not found",
			makefile: makeFileSample1,
			err:      ErrDependencyNotFound,
		}, {
			name:   "invalid command",
			target: "a",
			makefile: `
a:
	@invalid command
		`,
			err: ErrCouldntExecuteCommand,
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
			err = gomake.Build(f)
			assertErr(t, err, nil)

			err = gomake.Run(tc.target)
			assertErr(t, err, tc.err)
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

func assertEqualDepGraphs(t testing.TB, got, want map[string][]string) {
	t.Helper()
	for k, v := range got {
		if want[k] == nil {
			t.Errorf("got %v want %v", got, want)
			t.Fatalf("key %q exists in got and not exist in want", k)
		}

		sort.Strings(want[k])
		sort.Strings(v)

		if !reflect.DeepEqual(v, want[k]) {
			t.Errorf("in key %q got value %v want value %v", k, v, want[k])
		}
	}

	if len(got) != len(want) {
		t.Errorf("there are keys in want not exist in got\ngot %v want %v", got, want)
	}
}
