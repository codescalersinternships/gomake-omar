package internal

import (
	"reflect"
	"testing"
)

func TestGetCycle(t *testing.T) {
	t.Run("no cycle", func(t *testing.T) {
		graph := newGraph(map[string][]string{
			"a": {"b", "c"},
			"b": {"c"},
			"x": {"y"},
		})
		cycleGot := graph.getCycle()

		assertEqualCycles(t, cycleGot, []string{})
	})

	t.Run("there is a cycle", func(t *testing.T) {
		graph := newGraph(map[string][]string{
			"a": {"c", "x"},
			"c": {"d"},
			"d": {"a"},
			"x": {"y"},
		})
		cycleGot := graph.getCycle()
		cycleWant := []string{"a", "c", "d", "a"}

		assertEqualCycles(t, cycleGot, cycleWant)
	})
}

func TestGetDependency(t *testing.T) {
	t.Run("no cycle", func(t *testing.T) {
		graph := newGraph(map[string][]string{
			"r": {"a", "o"},
			"a": {"o", "m"},
			"c": {},
			"d": {},
			"x": {"y"},
		})
		depGot := graph.getOrderedDependencies("r")
		depWant := []string{"o", "m", "a", "r"}

		if !reflect.DeepEqual(depGot, depWant) {
			t.Errorf("got %v want %v", depGot, depWant)
		}
	})
}

func assertEqualCycles(t testing.TB, cycleGot, cycleWant []string) {
	t.Helper()
	if len(cycleGot) == 0 && len(cycleWant) == 0 {
		return
	}
	if len(cycleGot) != len(cycleWant) {
		t.Fatalf("got %v want %v", cycleGot, cycleWant)
	}

	// to check for slices have the same relative order
	// make them have the same value at index 0
	// then check if they are identical

	// try to make 'cycleWant' have the same value as 'cycleGot' at index 0
	// by cyclic shift 'cycleWant' till values at index 0 equalize
	for i := 0; i < len(cycleWant) && cycleWant[0] != cycleGot[0]; i++ {
		cycleWant = append(cycleWant[1:], cycleWant[0])
	}

	if !reflect.DeepEqual(cycleGot, cycleWant) {
		t.Errorf("got %v want %v", cycleGot, cycleWant)
	}
}
