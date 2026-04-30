package models

type Position struct{ X, Y int }

// SolverResult adalah output dari Algoritma yang akan dikonsumsi oleh GUI
type SolverResult struct {
	PathHistory []Position // Array history pergerakan (untuk animasi playback)
	TotalCost   int
	TimeMs      int64
	NodesEval   int
}

// MapData adalah representasi papan
type MapData struct {
	Grid          [][]rune
	Costs         [][]int
	Width, Height int
}
