package solver

import "github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"

// biar ga type casting, klo math.Abs() dia casting ke float64 dulu jd bisa lebih lambat
func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func manhattanDistance(p1, p2 models.Position) int {
	return absInt(p1.X-p2.X) + absInt(p2.X-p2.Y)
}

// Pure Manhattan
func Heuristic1(state *models.GameState, m *models.MapData) int {
	return manhattanDistance(state.Pos, m.GoalPos)
}

// Stop-point
// func Heuristic2

// Manhattan Checkpoint
// Jarak posisi saat ini -> sisa angka berurutan -> goal
func Heuristic3(state *models.GameState, m *models.MapData) int {
	if state.NextNum == -1 {
		return manhattanDistance(state.Pos, m.GoalPos)
	}

	cost := 0

	cost += manhattanDistance(state.Pos, m.NumberPos[state.NextNum])

	for i := state.NextNum; i < m.TotalNumbers-1; i++ {
		cost += manhattanDistance(m.NumberPos[i], m.NumberPos[i+1])
	}

	if m.TotalNumbers > 0 {
		cost += manhattanDistance(m.NumberPos[m.TotalNumbers-1], m.GoalPos)
	}

	return cost
}
