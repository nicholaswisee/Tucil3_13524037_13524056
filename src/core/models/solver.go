package models

type SolverResult struct {
	Path        []MoveRecord // urutan gerakan beserta posisi akhir & cost per gerakan
	PathHistory []Position   // posisi berhenti setiap step (termasuk posisi awal untuk playback)
	TotalCost   int
	TimeMs      int64
	NodesEval   int
	Success     bool
	Algorithm   string
	Heuristic   string
}

type Solver interface {
	Solve(m *MapData) (*SolverResult, error)
	Name() string
}
