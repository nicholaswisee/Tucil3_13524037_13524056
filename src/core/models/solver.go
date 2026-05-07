package models

const MaxSearchFrames = 5000

type SearchFrame struct {
	Current  Position   // node being expanded this step
	Children []Position // nodes added to frontier by expanding Current
}

type SolverResult struct {
	Path         []MoveRecord // urutan gerakan beserta posisi akhir & cost per gerakan
	PathHistory  []Position   // posisi berhenti setiap step (termasuk posisi awal untuk playback)
	SearchFrames []SearchFrame // per-expansion snapshots
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
