package internal

type graph struct {
	adjacencyList map[target][]parent

	isVisited   map[target]bool
	isExploring map[target]bool
	dfsStack    []target
}

func newGraph(adjacencyList map[target][]parent) graph {
	return graph{
		adjacencyList: adjacencyList,
		isVisited:     map[target]bool{},
		isExploring:   map[target]bool{},
		dfsStack:      []target{},
	}
}

func (g *graph) isCyclic(currentNode target) bool {
	g.isVisited[currentNode] = true
	g.isExploring[currentNode] = true
	g.dfsStack = append(g.dfsStack, currentNode)

	for _, nextNode := range g.adjacencyList[currentNode] {
		nextNode := target(nextNode)

		if g.isExploring[nextNode] ||
			(!g.isVisited[nextNode] && g.isCyclic(nextNode)) {
			return true
		}
	}

	g.dfsStack = g.dfsStack[:len(g.dfsStack)-1]
	g.isExploring[currentNode] = false
	return false
}

func (g *graph) topologicalSort(currentNode target) []target {
	g.isVisited[currentNode] = true

	targetsOrder := []target{}
	for _, nextNode := range g.adjacencyList[currentNode] {
		nextNode := target(nextNode)
		if !g.isVisited[nextNode] {
			targetsOrder = append(targetsOrder, g.topologicalSort(nextNode)...)
		}
	}

	targetsOrder = append(targetsOrder, currentNode)
	return targetsOrder
}

func (g *graph) getDependency(currentTarget target) []target {
	g.isVisited = map[target]bool{}
	return g.topologicalSort(currentTarget)
}

func (g *graph) getCycle() []target {
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
		return []target{}
	}

	return g.dfsStack
}
