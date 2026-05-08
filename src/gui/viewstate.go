package gui

import "github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"

type ViewState struct {
	MapData             *models.MapData
	Result              *models.SolverResult
	CurrentStep         int
	SearchStep          int
	SearchPhase         bool
	accumulatedVisited  map[models.Position]bool
	accumulatedFrontier map[models.Position]bool
}

func NewViewState() *ViewState {
	return &ViewState{}
}

func (vs *ViewState) SetMap(m *models.MapData) {
	vs.MapData = m
	vs.Result = nil
	vs.CurrentStep = 0
	vs.SearchStep = 0
	vs.SearchPhase = false
	vs.accumulatedVisited = nil
	vs.accumulatedFrontier = nil
}

func (vs *ViewState) SetResult(r *models.SolverResult) {
	vs.Result = r
	vs.CurrentStep = 0
	vs.SearchStep = 0
	vs.SearchPhase = true
	vs.rebuildAccumulatedSets()
}

func (vs *ViewState) rebuildAccumulatedSets() {
	vs.accumulatedVisited = make(map[models.Position]bool)
	vs.accumulatedFrontier = make(map[models.Position]bool)
	if vs.Result == nil || len(vs.Result.SearchFrames) == 0 {
		return
	}
	// The first frame's Current starts in the frontier (it's the initial node to expand)
	vs.accumulatedFrontier[vs.Result.SearchFrames[0].Current] = true
	for i := 0; i <= vs.SearchStep && i < len(vs.Result.SearchFrames); i++ {
		frame := vs.Result.SearchFrames[i]
		vs.accumulatedVisited[frame.Current] = true
		delete(vs.accumulatedFrontier, frame.Current)
		for _, child := range frame.Children {
			vs.accumulatedFrontier[child] = true
		}
	}
}

func (vs *ViewState) StepForward() bool {
	if vs.Result == nil || !vs.Result.Success {
		return false
	}
	if vs.CurrentStep < len(vs.Result.PathHistory)-1 {
		vs.CurrentStep++
		return true
	}
	return false
}

func (vs *ViewState) StepBackward() bool {
	if vs.CurrentStep > 0 {
		vs.CurrentStep--
		return true
	}
	return false
}

func (vs *ViewState) JumpToStart() {
	vs.CurrentStep = 0
}

func (vs *ViewState) JumpToEnd() {
	if vs.Result != nil && vs.Result.Success {
		vs.CurrentStep = len(vs.Result.PathHistory) - 1
	}
}

func (vs *ViewState) SearchForward() bool {
	if vs.Result == nil || len(vs.Result.SearchFrames) == 0 {
		return false
	}
	if vs.SearchStep < len(vs.Result.SearchFrames)-1 {
		vs.SearchStep++
		vs.rebuildAccumulatedSets()
		return true
	}
	return false
}

func (vs *ViewState) SearchBackward() bool {
	if vs.SearchStep > 0 {
		vs.SearchStep--
		vs.rebuildAccumulatedSets()
		return true
	}
	return false
}

func (vs *ViewState) JumpToSearchStart() {
	vs.SearchStep = 0
	vs.rebuildAccumulatedSets()
}

func (vs *ViewState) JumpToSearchEnd() {
	if vs.Result != nil && len(vs.Result.SearchFrames) > 0 {
		vs.SearchStep = len(vs.Result.SearchFrames) - 1
		vs.rebuildAccumulatedSets()
	}
}

func (vs *ViewState) TogglePhase() {
	vs.SearchPhase = !vs.SearchPhase
	if !vs.SearchPhase {
		vs.CurrentStep = 0
	} else {
		vs.rebuildAccumulatedSets()
	}
}

func (vs *ViewState) CurrentPos() models.Position {
	if vs.MapData == nil {
		return models.Position{}
	}
	if vs.SearchPhase && vs.Result != nil && vs.SearchStep >= 0 && vs.SearchStep < len(vs.Result.SearchFrames) {
		return vs.Result.SearchFrames[vs.SearchStep].Current
	}
	if vs.Result != nil && vs.Result.Success && vs.CurrentStep >= 0 && vs.CurrentStep < len(vs.Result.PathHistory) {
		return vs.Result.PathHistory[vs.CurrentStep]
	}
	return vs.MapData.StartPos
}

func (vs *ViewState) VisitedSet() map[models.Position]bool {
	return vs.accumulatedVisited
}

func (vs *ViewState) FrontierSet() map[models.Position]bool {
	return vs.accumulatedFrontier
}

func slidePositions(from, to models.Position) []models.Position {
	if from == to {
		return nil
	}
	var dr, dc int
	if to.X > from.X {
		dr = 1
	} else if to.X < from.X {
		dr = -1
	}
	if to.Y > from.Y {
		dc = 1
	} else if to.Y < from.Y {
		dc = -1
	}
	var positions []models.Position
	curr := from
	for curr != to {
		curr.X += dr
		curr.Y += dc
		positions = append(positions, curr)
	}
	return positions
}

func (vs *ViewState) JumpToStep(step int) {
	if vs.Result == nil || !vs.Result.Success {
		return
	}
	maxStep := len(vs.Result.PathHistory) - 1
	if step < 0 {
		step = 0
	}
	if step > maxStep {
		step = maxStep
	}
	vs.CurrentStep = step
}

func (vs *ViewState) JumpToSearchFrame(step int) {
	if vs.Result == nil || len(vs.Result.SearchFrames) == 0 {
		return
	}
	maxStep := len(vs.Result.SearchFrames) - 1
	if step < 0 {
		step = 0
	}
	if step > maxStep {
		step = maxStep
	}
	vs.SearchStep = step
	vs.rebuildAccumulatedSets()
}

func (vs *ViewState) CheckpointsPassed() map[int]bool {
	passed := make(map[int]bool)
	if vs.MapData == nil || vs.Result == nil || !vs.Result.Success || vs.MapData.TotalNumbers == 0 {
		return passed
	}
	history := vs.Result.PathHistory
	if len(history) < 2 {
		return passed
	}
	nextNum := 0
	for step := 0; step < vs.CurrentStep && step < len(history)-1; step++ {
		for _, p := range slidePositions(history[step], history[step+1]) {
			if vs.MapData.TileAt(p) == models.TileNumber {
				numIdx := int(vs.MapData.Grid[p.X][p.Y] - '0')
				if numIdx == nextNum {
					passed[numIdx] = true
					nextNum++
				}
			}
		}
		if nextNum >= vs.MapData.TotalNumbers {
			break
		}
	}
	return passed
}
