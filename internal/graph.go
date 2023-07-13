package internal

type graph struct {
	adjacencyList map[string][]string

	isVisited   map[string]bool
	isExploring map[string]bool
	dfsStack    []string
}

func newGraph(adjacencyList map[string][]string) graph {
	return graph{
		adjacencyList: adjacencyList,
		isVisited:     map[string]bool{},
		isExploring:   map[string]bool{},
		dfsStack:      []string{},
	}
}

func (g *graph) isCyclic(currentNode string) bool {
	g.isVisited[currentNode] = true
	g.isExploring[currentNode] = true
	g.dfsStack = append(g.dfsStack, currentNode)

	for _, nextNode := range g.adjacencyList[currentNode] {
		if g.isExploring[nextNode] ||
			(!g.isVisited[nextNode] && g.isCyclic(nextNode)) {
			return true
		}
	}

	g.dfsStack = g.dfsStack[:len(g.dfsStack)-1]
	g.isExploring[currentNode] = false

	return false
}

func (g *graph) topologicalSort(currentNode string) []string {
	g.isVisited[currentNode] = true

	targetsOrder := []string{}
	for _, nextNode := range g.adjacencyList[currentNode] {
		if !g.isVisited[nextNode] {
			targetsOrder = append(targetsOrder, g.topologicalSort(nextNode)...)
		}
	}

	targetsOrder = append(targetsOrder, currentNode)
	return targetsOrder
}

func (g *graph) getOrderedDependencies(currentTarget string) []string {
	g.isVisited = map[string]bool{}
	return g.topologicalSort(currentTarget)
}

func (g *graph) getCycle() []string {
	isCyclic := false

	for k := range g.adjacencyList {
		if g.isVisited[k] {
			continue
		}

		isCyclic = isCyclic || g.isCyclic(k)
		if isCyclic {
			break
		}
	}

	if !isCyclic {
		return []string{}
	}

	return g.dfsStack
}
