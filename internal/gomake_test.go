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

publish: test gendocs
	echo 'executing publish'

gendocs: build
	echo 'executing gendocs'
`

func graphSample1() map[target][]parent {
	return map[target][]parent{
		"build":   {},
		"test":    {},
		"publish": {"test", "gendocs"},
		"gendocs": {"build"},
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

		assertEqualGraphs(t, gomake.dependencyGraph, graphSample1())
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

}

func assertErr(t testing.TB, got, want error) {
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

func assertEqualGraphs(t testing.TB, got, want map[target][]parent) {
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
