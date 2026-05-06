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
	return absInt(p1.X-p2.X) + absInt(p1.Y-p2.Y)
}

func checkRock(m *models.MapData, x, y int) bool {
	if !m.InBounds(models.Position{X: x, Y: y}) {
		return true // out of bounds acts as a wall
	}
	return m.TileAt(models.Position{X: x, Y: y}) == models.TileWall
}

// Pure Manhattan
func Heuristic1(state *models.GameState, m *models.MapData) int {
	return manhattanDistance(state.Pos, m.GoalPos)
}

// Overshoot-penalty stop point
func Heuristic2(state *models.GameState, m *models.MapData) int {
	target := m.GoalPos
	if state.NextNum != -1 {
		target = m.NumberPos[state.NextNum]
	}

	currX, currY := state.Pos.X, state.Pos.Y
	targetX, targetY := target.X, target.Y

	minCost := m.MinCost
	if minCost <= 0 {
		minCost = 1
	}

	h := manhattanDistance(state.Pos, target) * minCost

	if currX != targetX && currY != targetY {
		corner1HasRock := checkRock(m, currX, targetY+1) || checkRock(m, currX, targetY-1)
		corner2HasRock := checkRock(m, targetX+1, currY) || checkRock(m, targetX-1, currY)

		if !corner1HasRock && !corner2HasRock {
			h += 2 * minCost
		}
	}

	return h
}

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
