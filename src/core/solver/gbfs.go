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
	if g.HeuristicID == 3 {
		return Heuristic3(state, m)
	}
	return Heuristic1(state, m)
}

func (g *GBFSSolver) Solve(m *models.MapData) (*models.SolverResult, error) {
	startTime := time.Now()

	startNum := 0
	if m.TotalNumbers == 0 {
		startNum = -1
	}
	initialState := models.GameState{Pos: m.StartPos, NextNum: startNum}

	pq := make(PriorityQueue, 0, 1000)
	heap.Init(&pq)

	initH := g.getH(&initialState, m)

	heap.Push(&pq, &SearchNode{
		State:       initialState,
		Priority:    initH,
		Cost:        0,
		Path:        make([]models.MoveRecord, 0),
		PathHistory: []models.Position{m.StartPos},
	})

	visited := make(map[models.StateKey]bool)
	visited[initialState.GetKey()] = true
	nodesEvaluated := 0
	var searchFrames []models.SearchFrame

	for pq.Len() > 0 {
		currNode := heap.Pop(&pq).(*SearchNode)
		nodesEvaluated++

		if currNode.State.IsGoal(m) {
			return &models.SolverResult{
				Path:         currNode.Path,
				PathHistory:  currNode.PathHistory,
				SearchFrames: searchFrames,
				TotalCost:    currNode.Cost,
				TimeMs:       time.Since(startTime).Milliseconds(),
				NodesEval:    nodesEvaluated,
				Success:      true,
				Algorithm:    "GBFS",
			}, nil
		}

		var nextStates [4]models.GameState
		var nextCosts [4]int
		var nextDirs [4]models.Direction
		var validChildren []models.Position
		count := 0

		for _, dir := range models.Directions {
			newState, moveCost, isValid := currNode.State.Slide(m, dir)
			if isValid {
				nextStates[count] = newState
				nextCosts[count] = moveCost
				nextDirs[count] = dir
				validChildren = append(validChildren, newState.Pos)
				count++
			}
		}

		if len(searchFrames) < models.MaxSearchFrames {
			pathCopy := make([]models.Position, len(currNode.PathHistory))
			copy(pathCopy, currNode.PathHistory)
			searchFrames = append(searchFrames, models.SearchFrame{
				Current:    currNode.State.Pos,
				Children:   validChildren,
				PathToNode: pathCopy,
			})
		}

		for i := 0; i < count; i++ {
			newState := nextStates[i]
			moveCost := nextCosts[i]
			dir := nextDirs[i]
			stateKey := newState.GetKey()

			if !visited[stateKey] {
				visited[stateKey] = true
				hCost := g.getH(&newState, m)
				newGCost := currNode.Cost + moveCost

				newPath := make([]models.MoveRecord, len(currNode.Path), len(currNode.Path)+1)
				copy(newPath, currNode.Path)
				newPath = append(newPath, models.MoveRecord{Direction: dir, NewPos: newState.Pos, MoveCost: moveCost})

				newHistory := make([]models.Position, len(currNode.PathHistory), len(currNode.PathHistory)+1)
				copy(newHistory, currNode.PathHistory)
				newHistory = append(newHistory, newState.Pos)

				heap.Push(&pq, &SearchNode{
					State:       newState,
					Priority:    hCost,
					Cost:        newGCost,
					Path:        newPath,
					PathHistory: newHistory,
				})
			}
		}
	}

	return &models.SolverResult{Success: false, SearchFrames: searchFrames, TimeMs: time.Since(startTime).Milliseconds(), NodesEval: nodesEvaluated, Algorithm: "GBFS"}, nil
}
