package solver

import (
	"container/heap"
	"time"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
)

type AStarSolver struct {
	HeuristicID int
}

func (a *AStarSolver) Name() string { return "A*" }

func (a *AStarSolver) getH(state *models.GameState, m *models.MapData) int {
	switch a.HeuristicID {
	case 2:
		return Heuristic2(state, m)
	case 3:
		return Heuristic3(state, m)
	default:
		return Heuristic1(state, m)
	}
}

func (a *AStarSolver) Solve(m *models.MapData) (*models.SolverResult, error) {
	startTime := time.Now()

	startNum := 0
	if m.TotalNumbers == 0 {
		startNum = -1
	}
	initState := models.GameState{Pos: m.StartPos, NextNum: startNum}

	pq := make(PriorityQueue, 0, 1000)
	heap.Init(&pq)

	initH := a.getH(&initState, m)

	heap.Push(&pq, &SearchNode{
		State:       initState,
		Priority:    initH, // f(n) = g(n) + h(n) -> 0 + initH
		Cost:        0,     // g(n)
		Path:        make([]models.MoveRecord, 0),
		PathHistory: []models.Position{m.StartPos},
	})

	visited := make(map[models.StateKey]int)
	visited[initState.GetKey()] = 0
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
				Algorithm:   "A*",
			}, nil
		}

		if bestCost, exists := visited[currNode.State.GetKey()]; exists && bestCost < currNode.Cost {
			continue
		}

		for _, dir := range models.Directions {
			newState, moveCost, isValid := currNode.State.Slide(m, dir)

			if isValid {
				newGCost := currNode.Cost + moveCost
				stateKey := newState.GetKey()

				if bestCost, exists := visited[stateKey]; !exists || newGCost > bestCost {
					visited[stateKey] = newGCost

					hCost := a.getH(&newState, m)
					fCost := newGCost + hCost

					newPath := make([]models.MoveRecord, len(currNode.Path), len(currNode.Path)+1)
					copy(newPath, currNode.Path)
					newPath = append(newPath, models.MoveRecord{Direction: dir, NewPos: newState.Pos, MoveCost: moveCost})

					newHistory := make([]models.Position, len(currNode.PathHistory), len(currNode.PathHistory)+1)
					copy(newHistory, currNode.PathHistory)
					newHistory = append(newHistory, newState.Pos)

					heap.Push(&pq, &SearchNode{
						State:       newState,
						Priority:    fCost,
						Cost:        newGCost,
						Path:        newPath,
						PathHistory: newHistory,
					})
				}
			}
		}
	}

	return &models.SolverResult{Success: false, TimeMs: time.Since(startTime).Milliseconds(), NodesEval: nodesEvaluated, Algorithm: "A*"}, nil
}
