package gui

import "github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"

type ViewState struct {
	MapData     *models.MapData
	Result      *models.SolverResult
	CurrentStep int
}

func NewViewState() *ViewState {
	return &ViewState{}
}

func (vs *ViewState) SetMap(m *models.MapData) {
	vs.MapData = m
	vs.Result = nil
	vs.CurrentStep = 0
}

func (vs *ViewState) SetResult(r *models.SolverResult) {
	vs.Result = r
	vs.CurrentStep = 0
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

func (vs *ViewState) CurrentPos() models.Position {
	if vs.MapData == nil {
		return models.Position{}
	}
	if vs.Result != nil && vs.Result.Success && vs.CurrentStep >= 0 && vs.CurrentStep < len(vs.Result.PathHistory) {
		return vs.Result.PathHistory[vs.CurrentStep]
	}
	return vs.MapData.StartPos
}
