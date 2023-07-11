package internal

type graph struct {
	adjacencyList map[target][]parent

	isVisited   map[target]bool
	isExploring map[target]bool
	dfsStack    []target
}

func NewGraph(adjacencyList map[target][]parent) graph {
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
		if g.isExploring[nextNode] {
			return true
		} else if !g.isVisited[nextNode] {
			if g.isCyclic(nextNode) {
				return true
			}
		}
	}

	g.dfsStack = g.dfsStack[:len(g.dfsStack)-1]
	g.isExploring[currentNode] = false
	return false
}

func (g *graph) dfs(currentNode target) []target {
	g.isVisited[currentNode] = true

	targetsRet := []target{currentNode}
	for _, nextNode := range g.adjacencyList[currentNode] {
		nextNode := target(nextNode)
		if !g.isVisited[nextNode] {
			targetsRet = append(targetsRet, g.dfs(nextNode)...)
		}
	}

	return targetsRet
}

func (g *graph) getDependency(currentTarget target) []target {
	g.isVisited = map[target]bool{}
	return g.dfs(currentTarget)
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
