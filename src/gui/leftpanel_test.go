package gui

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestLeftPanel_Creation(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	lp := NewLeftPanel()
	if lp == nil {
		t.Fatal("NewLeftPanel returned nil")
	}
	if lp.Object() == nil {
		t.Error("Object() returned nil")
	}
	if lp.AlgorithmSelect == nil {
		t.Error("AlgorithmSelect is nil")
	}
	if lp.HeuristicSelect == nil {
		t.Error("HeuristicSelect is nil")
	}
	if lp.ImportBtn == nil {
		t.Error("ImportBtn is nil")
	}
	if lp.ExportBtn == nil {
		t.Error("ExportBtn is nil")
	}
	if lp.RunBtn == nil {
		t.Error("RunBtn is nil")
	}
	if lp.StepLabel == nil {
		t.Error("StepLabel is nil")
	}
	if lp.SolutionLabel == nil {
		t.Error("SolutionLabel is nil")
	}
	if lp.CostLabel == nil {
		t.Error("CostLabel is nil")
	}
	if lp.TimeLabel == nil {
		t.Error("TimeLabel is nil")
	}
	if lp.IterationsLabel == nil {
		t.Error("IterationsLabel is nil")
	}
}

func TestLeftPanel_SetStats(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	lp := NewLeftPanel()
	lp.SetStats(150, 42)
	if lp.TimeLabel.Text != ">> Waktu eksekusi: 150 ms" {
		t.Errorf("unexpected time label: %s", lp.TimeLabel.Text)
	}
	if lp.IterationsLabel.Text != ">> Banyak iterasi yang dilakukan: 42 iterasi" {
		t.Errorf("unexpected iterations label: %s", lp.IterationsLabel.Text)
	}
}

func TestLeftPanel_SetStepLabel(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	lp := NewLeftPanel()
	lp.SetStepLabel(3, 10)
	if lp.StepLabel.Text != "Step 3 / 10" {
		t.Errorf("unexpected step label: %s", lp.StepLabel.Text)
	}
}

func TestLeftPanel_SetHeuristicEnabled(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	lp := NewLeftPanel()

	lp.SetHeuristicEnabled(true)
	if lp.HeuristicSelect.Disabled() {
		t.Error("expected heuristic select to be enabled")
	}

	lp.SetHeuristicEnabled(false)
	if !lp.HeuristicSelect.Disabled() {
		t.Error("expected heuristic select to be disabled")
	}
}

func TestLeftPanel_SetSolution(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	lp := NewLeftPanel()
	lp.SetSolution("RULUDRUR")
	if lp.SolutionLabel.Text != "Solusi Yang Ditemukan: RULUDRUR" {
		t.Errorf("unexpected solution label: %s", lp.SolutionLabel.Text)
	}
}

func TestLeftPanel_SetCost(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	lp := NewLeftPanel()
	lp.SetCost(87)
	if lp.CostLabel.Text != "Cost dari Solusi: 87" {
		t.Errorf("unexpected cost label: %s", lp.CostLabel.Text)
	}
}

func TestLeftPanel_SetNoSolution(t *testing.T) {
	testApp := test.NewApp()
	defer testApp.Quit()
	lp := NewLeftPanel()
	lp.SetNoSolution()
	if lp.SolutionLabel.Text != "Solusi Yang Ditemukan: Tidak ada solusi" {
		t.Errorf("unexpected solution label: %s", lp.SolutionLabel.Text)
	}
}