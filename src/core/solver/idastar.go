package solver

import (
	"math"
	"time"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
)

type IDAStarSolver struct {
	HeuristicID int
}

func (i *IDAStarSolver) Name() string { return "IDA*" }

func (i *IDAStarSolver) getH(state *models.GameState, m *models.MapData) int {
	switch i.HeuristicID {
	case 2:
		return Heuristic2(state, m)
	case 3:
		return Heuristic3(state, m)
	default:
		return Heuristic1(state, m)
	}
}

const foundCode = -1

// memori global biar ga allocate memori baru terus pas rekursi
type idaContext struct {
	m            *models.MapData
	getH         func(*models.GameState, *models.MapData) int
	path         []models.MoveRecord
	history      []models.Position
	inPath       map[models.StateKey]bool
	threshold    int
	nodesEval    int
	searchFrames []models.SearchFrame
}

func (ctx *idaContext) search(state models.GameState, g int) int {
	ctx.nodesEval++

	ctx.searchFrames = append(ctx.searchFrames, models.SearchFrame{
		Current:  state.Pos,
		Visited:  extractPositionsFromMap(ctx.inPath),
		Frontier: []models.Position{},
	})

	h := ctx.getH(&state, ctx.m)
	f := g + h

	if f > ctx.threshold {
		return f
	}

	if state.IsGoal(ctx.m) {
		return foundCode
	}

	stateKey := state.GetKey()
	ctx.inPath[stateKey] = true

	minCost := math.MaxInt

	for _, dir := range models.Directions {
		newState, moveCost, isValid := state.Slide(ctx.m, dir)

		if isValid {
			childKey := newState.GetKey()

			if !ctx.inPath[childKey] {
				ctx.path = append(ctx.path, models.MoveRecord{
					Direction: dir,
					NewPos:    newState.Pos,
					MoveCost:  moveCost,
				})
				ctx.history = append(ctx.history, newState.Pos)

				res := ctx.search(newState, g+moveCost)

				if res == foundCode {
					return foundCode
				}

				if res < minCost {
					minCost = res
				}

				ctx.path = ctx.path[:len(ctx.path)-1]
				ctx.history = ctx.history[:len(ctx.history)-1]
			}
		}
	}
	ctx.inPath[stateKey] = false

	return minCost
}

func (i *IDAStarSolver) Solve(m *models.MapData) (*models.SolverResult, error) {
	startTime := time.Now()

	startNum := 0
	if m.TotalNumbers == 0 {
		startNum = -1
	}
	initState := models.GameState{Pos: m.StartPos, NextNum: startNum}

	ctx := &idaContext{
		m:            m,
		getH:         i.getH,
		path:         make([]models.MoveRecord, 0, 100),
		history:      make([]models.Position, 0, 100),
		inPath:       make(map[models.StateKey]bool),
		searchFrames: make([]models.SearchFrame, 0, 1000),
	}

	ctx.history = append(ctx.history, m.StartPos)
	ctx.threshold = ctx.getH(&initState, m)

	for {
		res := ctx.search(initState, 0)

		if res == foundCode {
			finalPath := make([]models.MoveRecord, len(ctx.path))
			copy(finalPath, ctx.path)

			finalHistory := make([]models.Position, len(ctx.history))
			copy(finalHistory, ctx.history)

			totalCost := 0
			for _, mr := range finalPath {
				totalCost += mr.MoveCost
			}

			return &models.SolverResult{
				Path:         finalPath,
				PathHistory:  finalHistory,
				SearchFrames: ctx.searchFrames,
				TotalCost:    totalCost,
				TimeMs:       time.Since(startTime).Milliseconds(),
				NodesEval:    ctx.nodesEval,
				Success:      true,
				Algorithm:    "IDA*",
			}, nil
		}
		if res == math.MaxInt {
			break
		}

		ctx.threshold = res
	}

	return &models.SolverResult{
		Success:      false,
		TimeMs:       time.Since(startTime).Milliseconds(),
		NodesEval:    ctx.nodesEval,
		Algorithm:    "IDA*",
		SearchFrames: ctx.searchFrames,
	}, nil
}
