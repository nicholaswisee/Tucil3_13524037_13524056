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
