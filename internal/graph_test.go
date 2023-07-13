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
			"a": {"c", "b"},
			"c": {"d"},
			"d": {"a"},
			"x": {"y"},
		})
		cycleGot := graph.getCycle()
		cycleWant := []string{"a", "c", "d"}

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

	for i := 0; i < len(cycleWant) && cycleWant[0] != cycleGot[0]; i++ {
		cycleWant = append(cycleWant[1:], cycleWant[:1]...)
	}

	if !reflect.DeepEqual(cycleGot, cycleWant) {
		t.Errorf("got %v want %v", cycleGot, cycleWant)
	}
}