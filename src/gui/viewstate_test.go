package gui

import (
	"testing"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
)

func TestViewState_SetMap(t *testing.T) {
	vs := NewViewState()
	m := &models.MapData{Height: 3, Width: 3}
	vs.SetMap(m)
	if vs.MapData != m {
		t.Error("MapData not set")
	}
	if vs.CurrentStep != 0 {
		t.Error("CurrentStep should reset to 0")
	}
}

func TestViewState_Playback(t *testing.T) {
	vs := NewViewState()
	vs.MapData = &models.MapData{Height: 2, Width: 2, StartPos: models.Position{X: 0, Y: 0}}
	vs.Result = &models.SolverResult{
		Success:     true,
		PathHistory: []models.Position{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
	}

	if !vs.StepForward() {
		t.Error("expected StepForward to succeed")
	}
	if vs.CurrentStep != 1 {
		t.Errorf("expected step 1, got %d", vs.CurrentStep)
	}

	vs.JumpToEnd()
	if vs.CurrentStep != 2 {
		t.Errorf("expected step 2, got %d", vs.CurrentStep)
	}

	vs.JumpToStart()
	if vs.CurrentStep != 0 {
		t.Errorf("expected step 0, got %d", vs.CurrentStep)
	}

	if vs.StepBackward() {
		t.Error("expected StepBackward to fail at step 0")
	}
}

func TestViewState_SearchPlayback(t *testing.T) {
	vs := NewViewState()
	vs.MapData = &models.MapData{Height: 2, Width: 2, StartPos: models.Position{X: 0, Y: 0}}
	vs.Result = &models.SolverResult{
		Success: true,
		SearchFrames: []models.SearchFrame{
			{Current: models.Position{X: 0, Y: 0}, Children: []models.Position{{X: 0, Y: 1}}},
			{Current: models.Position{X: 0, Y: 1}, Children: []models.Position{{X: 1, Y: 1}}},
			{Current: models.Position{X: 1, Y: 1}, Children: []models.Position{}},
		},
		PathHistory: []models.Position{{X: 0, Y: 0}, {X: 1, Y: 1}},
	}

	vs.SetResult(vs.Result)
	if !vs.SearchPhase {
		t.Error("expected SearchPhase to be true after SetResult")
	}
	if vs.SearchStep != 0 {
		t.Errorf("expected SearchStep 0, got %d", vs.SearchStep)
	}

	// After applying frame 0: visited={0,0}, frontier={0,1}
	if !vs.VisitedSet()[models.Position{X: 0, Y: 0}] {
		t.Error("expected (0,0) to be visited after SetResult")
	}
	if !vs.FrontierSet()[models.Position{X: 0, Y: 1}] {
		t.Error("expected (0,1) to be in frontier after SetResult")
	}

	if !vs.SearchForward() {
		t.Error("expected SearchForward to succeed")
	}
	if vs.SearchStep != 1 {
		t.Errorf("expected SearchStep 1, got %d", vs.SearchStep)
	}

	// After applying frame 1: visited={0,0; 0,1}, frontier={1,1}
	if !vs.VisitedSet()[models.Position{X: 0, Y: 1}] {
		t.Error("expected (0,1) to be visited after SearchForward")
	}
	if !vs.FrontierSet()[models.Position{X: 1, Y: 1}] {
		t.Error("expected (1,1) to be in frontier after SearchForward")
	}
	if vs.FrontierSet()[models.Position{X: 0, Y: 1}] {
		t.Error("expected (0,1) to NOT be in frontier after being visited")
	}

	pos := vs.CurrentPos()
	if pos != (models.Position{X: 0, Y: 1}) {
		t.Errorf("expected current pos {0,1}, got %+v", pos)
	}

	vs.TogglePhase()
	if vs.SearchPhase {
		t.Error("expected SearchPhase to be false after TogglePhase")
	}
	if vs.CurrentStep != 0 {
		t.Errorf("expected CurrentStep 0 after toggle, got %d", vs.CurrentStep)
	}

	vs.TogglePhase()
	if !vs.SearchPhase {
		t.Error("expected SearchPhase to be true after second toggle")
	}
	if vs.SearchStep != 1 {
		t.Errorf("expected SearchStep restored to 1, got %d", vs.SearchStep)
	}

	vs.JumpToSearchEnd()
	if vs.SearchStep != 2 {
		t.Errorf("expected SearchStep 2, got %d", vs.SearchStep)
	}

	vs.SearchBackward()
	if vs.SearchStep != 1 {
		t.Errorf("expected SearchStep 1, got %d", vs.SearchStep)
	}

	vs.JumpToSearchStart()
	if vs.SearchStep != 0 {
		t.Errorf("expected SearchStep 0, got %d", vs.SearchStep)
	}
}
