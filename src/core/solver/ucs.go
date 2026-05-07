package solver

import (
	"container/heap"
	"time"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
)

type UCSSolver struct{}

func (u *UCSSolver) Name() string {
	return "UCS"
}

func (u *UCSSolver) Solve(m *models.MapData) (*models.SolverResult, error) {
	startTime := time.Now()

	startNum := 0

	if m.TotalNumbers == 0 {
		startNum = -1
	}

	initState := models.GameState{Pos: m.StartPos, NextNum: startNum}

	pq := make(PriorityQueue, 0, 1000)
	heap.Init(&pq)

	heap.Push(&pq, &SearchNode{
		State:       initState,
		Priority:    0,
		Cost:        0,
		Path:        make([]models.MoveRecord, 0),
		PathHistory: []models.Position{m.StartPos},
	})

	visited := make(map[models.StateKey]int)
	visited[initState.GetKey()] = 0

	nodesEvaluated := 0
	var searchFrames []models.SearchFrame

	for pq.Len() > 0 {
		currNode := heap.Pop(&pq).(*SearchNode)
		nodesEvaluated++

		var children []models.Position

		if currNode.State.IsGoal(m) {
			if len(searchFrames) < models.MaxSearchFrames {
				searchFrames = append(searchFrames, models.SearchFrame{
					Current:  currNode.State.Pos,
					Children: children,
				})
			}
			return &models.SolverResult{
				Path:         currNode.Path,
				PathHistory:  currNode.PathHistory,
				SearchFrames: searchFrames,
				TotalCost:    currNode.Cost,
				TimeMs:       time.Since(startTime).Milliseconds(),
				NodesEval:    nodesEvaluated,
				Success:      true,
				Algorithm:    "UCS",
			}, nil
		}

		if bestCost, exists := visited[currNode.State.GetKey()]; exists && bestCost < currNode.Cost {
			continue
		}

		for _, dir := range models.Directions {
			newState, moveCost, isValid := currNode.State.Slide(m, dir)

			if isValid {
				newCost := currNode.Cost + moveCost
				stateKey := newState.GetKey()

				if bestCost, exists := visited[stateKey]; !exists || newCost < bestCost {
					visited[stateKey] = newCost
					children = append(children, newState.Pos)

					newPath := make([]models.MoveRecord, len(currNode.Path), len(currNode.Path)+1)
					copy(newPath, currNode.Path)
					newPath = append(newPath, models.MoveRecord{Direction: dir, NewPos: newState.Pos, MoveCost: moveCost})

					newHistory := make([]models.Position, len(currNode.PathHistory), len(currNode.PathHistory)+1)
					copy(newHistory, currNode.PathHistory)
					newHistory = append(newHistory, newState.Pos)

					heap.Push(&pq, &SearchNode{
						State:       newState,
						Priority:    newCost,
						Cost:        newCost,
						Path:        newPath,
						PathHistory: newHistory,
					})
				}
			}
		}

		if len(searchFrames) < models.MaxSearchFrames {
			searchFrames = append(searchFrames, models.SearchFrame{
				Current:  currNode.State.Pos,
				Children: children,
			})
		}
	}

	return &models.SolverResult{Success: false, TimeMs: time.Since(startTime).Milliseconds(), NodesEval: nodesEvaluated, Algorithm: "UCS", SearchFrames: searchFrames}, nil
}
