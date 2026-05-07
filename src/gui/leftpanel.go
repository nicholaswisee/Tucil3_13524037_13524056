package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LeftPanel struct {
	container fyne.CanvasObject

	AlgorithmSelect *widget.Select
	HeuristicSelect *widget.Select
	ImportBtn       *widget.Button
	ExportBtn       *widget.Button
	RunBtn          *widget.Button
	FirstStepBtn    *widget.Button
	PrevStepBtn     *widget.Button
	NextStepBtn     *widget.Button
	LastStepBtn     *widget.Button
	PhaseToggleBtn  *widget.Button
	StepLabel       *widget.Label
	SearchStepLabel *widget.Label
	SolutionLabel   *widget.Label
	CostLabel       *widget.Label
	TimeLabel       *widget.Label
	IterationsLabel *widget.Label
	VisitedLabel    *widget.Label
	FrontierLabel   *widget.Label
}

func NewLeftPanel() *LeftPanel {
	lp := &LeftPanel{}

	lp.AlgorithmSelect = widget.NewSelect([]string{"UCS", "GBFS", "A*", "IDA*"}, func(string) {})
	lp.AlgorithmSelect.SetSelected("UCS")

	lp.HeuristicSelect = widget.NewSelect([]string{"Heuristic 1", "Heuristic 2", "Heuristic 3"}, func(string) {})
	lp.HeuristicSelect.SetSelected("Heuristic 1")
	lp.HeuristicSelect.Disable()

	lp.ImportBtn = widget.NewButton("Import Config", func() {})
	lp.ExportBtn = widget.NewButton("Export Config", func() {})
	lp.RunBtn = widget.NewButton("Run / Solve", func() {})

	lp.FirstStepBtn = widget.NewButton("|<", func() {})
	lp.PrevStepBtn = widget.NewButton("<", func() {})
	lp.NextStepBtn = widget.NewButton(">", func() {})
	lp.LastStepBtn = widget.NewButton(">|", func() {})

	lp.PhaseToggleBtn = widget.NewButton("Show Solution", func() {})

	lp.StepLabel = widget.NewLabel("Step 0 / 0")
	lp.SearchStepLabel = widget.NewLabel("Search Step 0 / 0")
	lp.SolutionLabel = widget.NewLabel("Solusi: -")
	lp.CostLabel = widget.NewLabel("Cost: -")
	lp.TimeLabel = widget.NewLabel("Waktu eksekusi: -")
	lp.IterationsLabel = widget.NewLabel("Banyak iterasi: -")
	lp.VisitedLabel = widget.NewLabel("Visited: -")
	lp.FrontierLabel = widget.NewLabel("Frontier: -")

	playbackRow := container.NewHBox(lp.FirstStepBtn, lp.PrevStepBtn, lp.NextStepBtn, lp.LastStepBtn, lp.StepLabel)

	lp.container = container.NewVBox(
		widget.NewLabel("Algorithm"),
		lp.AlgorithmSelect,
		widget.NewLabel("Heuristic"),
		lp.HeuristicSelect,
		widget.NewSeparator(),
		lp.ImportBtn,
		lp.ExportBtn,
		widget.NewSeparator(),
		lp.RunBtn,
		widget.NewSeparator(),
		lp.PhaseToggleBtn,
		widget.NewLabel("Playback"),
		playbackRow,
		lp.SearchStepLabel,
		widget.NewSeparator(),
		widget.NewLabel("Stats"),
		lp.SolutionLabel,
		lp.CostLabel,
		lp.TimeLabel,
		lp.IterationsLabel,
		lp.VisitedLabel,
		lp.FrontierLabel,
	)

	return lp
}

func (lp *LeftPanel) Object() fyne.CanvasObject {
	return lp.container
}

func (lp *LeftPanel) SetStats(timeMs int64, iterations int) {
	lp.TimeLabel.SetText(fmt.Sprintf(">> Waktu eksekusi: %d ms", timeMs))
	lp.IterationsLabel.SetText(fmt.Sprintf(">> Banyak iterasi yang dilakukan: %d iterasi", iterations))
}

func (lp *LeftPanel) SetStepLabel(current, total int) {
	lp.StepLabel.SetText(fmt.Sprintf("Step %d / %d", current, total))
}

func (lp *LeftPanel) SetSearchStepLabel(current, total int) {
	lp.SearchStepLabel.SetText(fmt.Sprintf("Search Step %d / %d", current, total))
}

func (lp *LeftPanel) SetSearchStats(totalSteps, currentStep, visitedCount int) {
	lp.SetSearchStepLabel(currentStep, totalSteps)
	lp.VisitedLabel.SetText(fmt.Sprintf("Visited: %d", visitedCount))
	lp.FrontierLabel.SetText(fmt.Sprintf("Frontier: %d", 0))
}

func (lp *LeftPanel) SetSearchStatsFromState(vs *ViewState) {
	if vs.Result == nil || len(vs.Result.SearchFrames) == 0 {
		lp.VisitedLabel.SetText("Visited: -")
		lp.FrontierLabel.SetText("Frontier: -")
		return
	}
	visited := vs.VisitedSet()
	frontier := vs.FrontierSet()
	lp.VisitedLabel.SetText(fmt.Sprintf("Visited: %d", len(visited)))
	lp.FrontierLabel.SetText(fmt.Sprintf("Frontier: %d", len(frontier)))
}

func (lp *LeftPanel) SetSolution(solution string) {
	lp.SolutionLabel.SetText(fmt.Sprintf("Solusi Yang Ditemukan: %s", solution))
}

func (lp *LeftPanel) SetCost(cost int) {
	lp.CostLabel.SetText(fmt.Sprintf("Cost dari Solusi: %d", cost))
}

func (lp *LeftPanel) SetNoSolution() {
	lp.SolutionLabel.SetText("Solusi Yang Ditemukan: Tidak ada solusi")
	lp.CostLabel.SetText("Cost dari Solusi: -")
	lp.TimeLabel.SetText("Waktu eksekusi: -")
	lp.IterationsLabel.SetText("Banyak iterasi: -")
}

func (lp *LeftPanel) SetHeuristicEnabled(enabled bool) {
	if enabled {
		lp.HeuristicSelect.Enable()
	} else {
		lp.HeuristicSelect.Disable()
	}
}

func (lp *LeftPanel) SetPhaseToggleLabel(searchPhase bool) {
	if searchPhase {
		lp.PhaseToggleBtn.SetText("Show Solution")
	} else {
		lp.PhaseToggleBtn.SetText("Show Search")
	}
}
