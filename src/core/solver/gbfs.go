package solver

import (
	"container/heap"
	"time"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
)

type GBFSSolver struct {
	HeuristicID int
}

func (g *GBFSSolver) Name() string { return "GBFS" }

func (g *GBFSSolver) getH(state *models.GameState, m *models.MapData) int {
	switch g.HeuristicID {
	case 2:
		return Heuristic2(state, m)
	case 3:
		return Heuristic3(state, m)
	default:
		return Heuristic1(state, m)
	}
}

func (g *GBFSSolver) Solve(m *models.MapData) (*models.SolverResult, error) {
	startTime := time.Now()

	startNum := 0
	if m.TotalNumbers == 0 {
		startNum = -1
	}
	initState := models.GameState{Pos: m.StartPos, NextNum: startNum}

	pq := make(PriorityQueue, 0, 1000)
	heap.Init(&pq)

	initH := g.getH(&initState, m)

	heap.Push(&pq, &SearchNode{
		State:       initState,
		Priority:    initH,
		Cost:        0,
		Path:        make([]models.MoveRecord, 0),
		PathHistory: []models.Position{m.StartPos},
	})

	visited := make(map[models.StateKey]bool)
	visited[initState.GetKey()] = true
	nodesEvaluated := 0

	for pq.Len() > 0 {
		currNode := heap.Pop(&pq).(*SearchNode)
		nodesEvaluated++

		if currNode.State.IsGoal(m) {
			return &models.SolverResult{
				Path:        currNode.Path,
				PathHistory: currNode.PathHistory,
				TotalCost:   currNode.Cost,
				TimeMs:      time.Since(startTime).Milliseconds(),
				NodesEval:   nodesEvaluated,
				Success:     true,
				Algorithm:   "GBFS",
			}, nil
		}

		for _, dir := range models.Directions {
			newState, moveCost, isValid := currNode.State.Slide(m, dir)

			if isValid {
				stateKey := newState.GetKey()

				if !visited[stateKey] {
					visited[stateKey] = true

					newCost := currNode.Cost + moveCost
					hCost := g.getH(&newState, m)

					newPath := make([]models.MoveRecord, len(currNode.Path), len(currNode.Path)+1)
					copy(newPath, currNode.Path)
					newPath = append(newPath, models.MoveRecord{Direction: dir, NewPos: newState.Pos, MoveCost: moveCost})

					newHistory := make([]models.Position, len(currNode.PathHistory), len(currNode.PathHistory)+1)
					copy(newHistory, currNode.PathHistory)
					newHistory = append(newHistory, newState.Pos)

					heap.Push(&pq, &SearchNode{
						State:       newState,
						Priority:    hCost,
						Cost:        newCost,
						Path:        newPath,
						PathHistory: newHistory,
					})
				}
			}
		}
	}

	return &models.SolverResult{
		Success:   false,
		TimeMs:    time.Since(startTime).Milliseconds(),
		NodesEval: nodesEvaluated,
		Algorithm: "GBFS",
	}, nil
}
