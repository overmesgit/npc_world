package pathfinding

import (
    "github.com/beefsack/go-astar"
    "github.com/solarlune/resolv"
    "math"
)

type PathNode struct {
    X, Y  int
    space *resolv.Space
}

func (n PathNode) PathNeighbors() []astar.Pather {
    var neighbors []astar.Pather
    straight := []struct{ dx, dy int }{
        {-1, 0}, {1, 0}, {0, -1}, {0, 1},
    }
    checks := [4]bool{}

    for i, dir := range straight {
        newX, newY := n.X+dir.dx, n.Y+dir.dy
        if newCell := n.space.Cell(newX, newY); newCell != nil && !newCell.ContainsTags("goblin_den", "mountain") {
            neighbors = append(neighbors, PathNode{X: newX, Y: newY, space: n.space})
            checks[i] = true
        }
    }

    diagonal := []struct {
        dx, dy int
        req    []int
    }{
        {-1, -1, []int{0, 2}},
        {1, -1, []int{2, 1}},
        {-1, 1, []int{0, 3}},
        {1, 1, []int{1, 3}},
    }
outer:
    for _, dir := range diagonal {
        newX, newY := n.X+dir.dx, n.Y+dir.dy
        for _, r := range dir.req {
            if !checks[r] {
                continue outer
            }
        }
        if newCell := n.space.Cell(newX, newY); newCell != nil && !newCell.ContainsTags("goblin_den", "mountain") {
            neighbors = append(neighbors, PathNode{X: newX, Y: newY, space: n.space})
        }
    }

    return neighbors
}

func (n PathNode) PathNeighborCost(to astar.Pather) float64 {
    toNode := to.(PathNode)
    diagonal := toNode.X != n.X && toNode.Y != n.Y
    if diagonal {
        return 1.414
    }
    return 1
}

func (n PathNode) PathEstimatedCost(to astar.Pather) float64 {
    toNode := to.(PathNode)
    dx := float64(abs(toNode.X - n.X))
    dy := float64(abs(toNode.Y - n.Y))
    return math.Sqrt(dx*dx + dy*dy)
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func FindPath(space *resolv.Space, startX, startY, endX, endY int) ([]astar.Pather, float64, bool) {
    startCell := space.Cell(startX, startY)
    endCell := space.Cell(endX, endY)

    if startCell == nil || endCell == nil {
        return nil, 0, false
    }

    startNode := PathNode{X: startX, Y: startY, space: space}
    endNode := PathNode{X: endX, Y: endY, space: space}

    return astar.Path(startNode, endNode)
}
