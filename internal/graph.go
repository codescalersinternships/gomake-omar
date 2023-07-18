package internal

type graph struct {
	adjacencyList map[string][]string
}

func newGraph(adjacencyList map[string][]string) graph {
	return graph{
		adjacencyList: adjacencyList,
	}
}

func (g *graph) isCyclic(currentNode string, dfsStack *[]string, isVisited, isExploring *map[string]bool) bool {
	(*isVisited)[currentNode] = true
	(*isExploring)[currentNode] = true
	(*dfsStack) = append((*dfsStack), currentNode)

	for _, nextNode := range g.adjacencyList[currentNode] {
		if (*isExploring)[nextNode] ||
			(!(*isVisited)[nextNode] && g.isCyclic(nextNode, dfsStack, isVisited, isExploring)) {
			return true
		}
	}

	(*dfsStack) = (*dfsStack)[:len((*dfsStack))-1]
	(*isExploring)[currentNode] = false

	return false
}

func (g *graph) topologicalSort(currentNode string, isVisited *map[string]bool) []string {
	(*isVisited)[currentNode] = true

	targetsOrder := []string{}
	for _, nextNode := range g.adjacencyList[currentNode] {
		if !(*isVisited)[nextNode] {
			targetsOrder = append(targetsOrder, g.topologicalSort(nextNode, isVisited)...)
		}
	}

	targetsOrder = append(targetsOrder, currentNode)
	return targetsOrder
}

func (g *graph) getOrderedDependencies(currentTarget string) []string {
	isVisited := make(map[string]bool)
	return g.topologicalSort(currentTarget, &isVisited)
}

func (g *graph) getCycle() []string {
	isExploring := make(map[string]bool)
	isVisited := make(map[string]bool)
	dfsStack := []string{}

	for k := range g.adjacencyList {
		if isVisited[k] {
			continue
		}

		if g.isCyclic(k, &dfsStack, &isVisited, &isExploring) {
			return dfsStack
		}
	}

	return []string{}
}
