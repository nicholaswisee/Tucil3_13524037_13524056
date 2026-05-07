package gui

import "github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"

type ViewState struct {
	MapData     *models.MapData
	Result      *models.SolverResult
	CurrentStep int
	SearchStep  int
	SearchPhase bool
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
}

func (vs *ViewState) SetResult(r *models.SolverResult) {
	vs.Result = r
	vs.CurrentStep = 0
	vs.SearchStep = 0
	vs.SearchPhase = true
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
		return true
	}
	return false
}

func (vs *ViewState) SearchBackward() bool {
	if vs.SearchStep > 0 {
		vs.SearchStep--
		return true
	}
	return false
}

func (vs *ViewState) JumpToSearchStart() {
	vs.SearchStep = 0
}

func (vs *ViewState) JumpToSearchEnd() {
	if vs.Result != nil {
		vs.SearchStep = len(vs.Result.SearchFrames) - 1
		if vs.SearchStep < 0 {
			vs.SearchStep = 0
		}
	}
}

func (vs *ViewState) TogglePhase() {
	vs.SearchPhase = !vs.SearchPhase
	if !vs.SearchPhase {
		vs.CurrentStep = 0
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
