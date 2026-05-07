package models

type SearchFrame struct {
	Current  Position
	Visited  []Position
	Frontier []Position
}

type SolverResult struct {
	Path         []MoveRecord // urutan gerakan beserta posisi akhir & cost per gerakan
	PathHistory  []Position   // posisi berhenti setiap step (termasuk posisi awal untuk playback)
	SearchFrames []SearchFrame // NEW: per-expansion snapshots
	TotalCost    int
	TimeMs       int64
	NodesEval    int
	Success      bool
	Algorithm    string
	Heuristic    string
}

type Solver interface {
	Solve(m *MapData) (*SolverResult, error)
	Name() string
}
