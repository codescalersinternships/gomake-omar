package internal

type graph struct {
	adjacencyList map[string][]string
}

func newGraph(adjacencyList map[string][]string) graph {
	return graph{
		adjacencyList: adjacencyList,
	}
}

func (g *graph) getCycleDfs(currentNode string, isVisited, isExploring map[string]bool) ([]string, bool) {
	isVisited[currentNode] = true
	isExploring[currentNode] = true

	for _, nextNode := range g.adjacencyList[currentNode] {
		if isExploring[nextNode] {
			return []string{currentNode}, true
		}

		if !isVisited[nextNode] {
			stk, isCyclic := g.getCycleDfs(nextNode, isVisited, isExploring)

			if isCyclic {
				return append([]string{currentNode}, stk...), true
			}
		}
	}

	isExploring[currentNode] = false
	return []string{}, false
}

func (g *graph) topologicalSort(currentNode string, isVisited map[string]bool) []string {
	isVisited[currentNode] = true

	targetsOrder := []string{}
	for _, nextNode := range g.adjacencyList[currentNode] {
		if !isVisited[nextNode] {
			targetsOrder = append(targetsOrder, g.topologicalSort(nextNode, isVisited)...)
		}
	}

	targetsOrder = append(targetsOrder, currentNode)
	return targetsOrder
}

func (g *graph) getOrderedDependencies(currentTarget string) []string {
	isVisited := make(map[string]bool)
	return g.topologicalSort(currentTarget, isVisited)
}

func (g *graph) getCycle() []string {
	isExploring := make(map[string]bool)
	isVisited := make(map[string]bool)

	for k := range g.adjacencyList {
		if isVisited[k] {
			continue
		}

		if stk, isCyclic := g.getCycleDfs(k, isVisited, isExploring); isCyclic {
			return stk
		}
	}

	return []string{}
}
