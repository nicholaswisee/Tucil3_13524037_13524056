package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type LeftPanel struct {
	container fyne.CanvasObject

	// Algorithm card
	AlgorithmSelect *widget.Select
	HeuristicSelect *widget.Select

	// Actions card
	ImportBtn      *widget.Button
	ExportBtn      *widget.Button
	RunBtn         *widget.Button
	PhaseToggleBtn *widget.Button
	ProgressBar    *widget.ProgressBarInfinite

	// Playback card
	FirstStepBtn   *widget.Button
	PrevStepBtn    *widget.Button
	NextStepBtn    *widget.Button
	LastStepBtn    *widget.Button
	PlayPauseBtn   *widget.Button
	PlaybackSlider *widget.Slider
	SpeedSelect    *widget.Select

	// Legacy label fields
	StepLabel       *widget.Label
	SearchStepLabel *widget.Label
	SolutionLabel   *widget.Label
	CostLabel       *widget.Label
	TimeLabel       *widget.Label
	IterationsLabel *widget.Label
	VisitedLabel    *widget.Label
	FrontierLabel   *widget.Label

	// Display value labels shown in the two-column stats grid.
	solValueDisp      *widget.Label
	costValueDisp     *widget.Label
	timeValueDisp     *widget.Label
	iterValueDisp     *widget.Label
	visitedValueDisp  *widget.Label
	frontierValueDisp *widget.Label
}

func NewLeftPanel() *LeftPanel {
	lp := &LeftPanel{}

	// Algorithm card
	lp.AlgorithmSelect = widget.NewSelect([]string{"UCS", "GBFS", "A*", "IDA*"}, func(string) {})
	lp.AlgorithmSelect.SetSelected("UCS")

	lp.HeuristicSelect = widget.NewSelect([]string{"Heuristic 1", "Heuristic 2", "Heuristic 3"}, func(string) {})
	lp.HeuristicSelect.SetSelected("Heuristic 1")
	lp.HeuristicSelect.Disable()

	algCard := widget.NewCard("Algorithm", "", container.NewVBox(
		widget.NewLabel("Algorithm"),
		lp.AlgorithmSelect,
		widget.NewLabel("Heuristic"),
		lp.HeuristicSelect,
	))

	// Actions card
	lp.ImportBtn = widget.NewButton("Import Config", func() {})
	lp.ExportBtn = widget.NewButton("Export Config", func() {})
	lp.RunBtn = widget.NewButton("Run / Solve", func() {})
	lp.RunBtn.Importance = widget.HighImportance
	lp.PhaseToggleBtn = widget.NewButton("Show Solution", func() {})
	lp.ProgressBar = widget.NewProgressBarInfinite()
	lp.ProgressBar.Hide()

	actCard := widget.NewCard("Actions", "", container.NewVBox(
		lp.ImportBtn,
		lp.ExportBtn,
		lp.RunBtn,
		lp.ProgressBar,
		lp.PhaseToggleBtn,
	))

	// Playback card
	lp.FirstStepBtn = widget.NewButton("|<", func() {})
	lp.PrevStepBtn = widget.NewButton("<", func() {})
	lp.NextStepBtn = widget.NewButton(">", func() {})
	lp.LastStepBtn = widget.NewButton(">|", func() {})
	lp.PlayPauseBtn = widget.NewButton("▶ Play", func() {})
	lp.PlaybackSlider = widget.NewSlider(0, 0)
	lp.PlaybackSlider.Step = 1
	lp.SpeedSelect = widget.NewSelect([]string{"0.5×", "1×", "2×", "4×"}, func(string) {})
	lp.SpeedSelect.SetSelected("1×")

	// Legacy labels kept for test compatibility (not placed in container).
	lp.StepLabel = widget.NewLabel("Step 0 / 0")
	lp.SearchStepLabel = widget.NewLabel("Search Step 0 / 0")

	navRow := container.NewHBox(
		lp.FirstStepBtn, lp.PrevStepBtn, lp.NextStepBtn, lp.LastStepBtn,
	)
	speedRow := container.NewHBox(
		widget.NewLabel("Speed:"), lp.SpeedSelect,
	)

	pbCard := widget.NewCard("Playback", "", container.NewVBox(
		lp.PlaybackSlider,
		lp.PlayPauseBtn,
		navRow,
		speedRow,
	))

	// Statistics card
	lp.SolutionLabel = widget.NewLabel("Solusi Yang Ditemukan: -")
	lp.CostLabel = widget.NewLabel("Cost dari Solusi: -")
	lp.TimeLabel = widget.NewLabel("Waktu eksekusi: -")
	lp.IterationsLabel = widget.NewLabel("Banyak iterasi: -")
	lp.VisitedLabel = widget.NewLabel("Visited: -")
	lp.FrontierLabel = widget.NewLabel("Frontier: -")

	bold := fyne.TextStyle{Bold: true}
	lp.solValueDisp = widget.NewLabelWithStyle("-", fyne.TextAlignLeading, bold)
	lp.costValueDisp = widget.NewLabelWithStyle("-", fyne.TextAlignLeading, bold)
	lp.timeValueDisp = widget.NewLabelWithStyle("-", fyne.TextAlignLeading, bold)
	lp.iterValueDisp = widget.NewLabelWithStyle("-", fyne.TextAlignLeading, bold)
	lp.visitedValueDisp = widget.NewLabelWithStyle("-", fyne.TextAlignLeading, bold)
	lp.frontierValueDisp = widget.NewLabelWithStyle("-", fyne.TextAlignLeading, bold)

	statsGrid := container.New(layout.NewGridLayout(2),
		widget.NewLabel("Solusi"), lp.solValueDisp,
		widget.NewLabel("Cost"), lp.costValueDisp,
		widget.NewLabel("Waktu"), lp.timeValueDisp,
		widget.NewLabel("Iterasi"), lp.iterValueDisp,
		widget.NewLabel("Visited"), lp.visitedValueDisp,
		widget.NewLabel("Frontier"), lp.frontierValueDisp,
	)
	statsCard := widget.NewCard("Statistics", "", statsGrid)

	// Assemble with scroll
	lp.container = container.NewVScroll(container.NewVBox(
		algCard,
		actCard,
		pbCard,
		statsCard,
	))

	return lp
}

func (lp *LeftPanel) Object() fyne.CanvasObject {
	return lp.container
}

func (lp *LeftPanel) SetStats(timeMs int64, iterations int) {
	tStr := fmt.Sprintf(">> Waktu eksekusi: %d ms", timeMs)
	iStr := fmt.Sprintf(">> Banyak iterasi yang dilakukan: %d iterasi", iterations)
	lp.TimeLabel.SetText(tStr)
	lp.IterationsLabel.SetText(iStr)
	lp.timeValueDisp.SetText(fmt.Sprintf("%d ms", timeMs))
	lp.iterValueDisp.SetText(fmt.Sprintf("%d", iterations))
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
	lp.visitedValueDisp.SetText(fmt.Sprintf("%d", visitedCount))
	lp.frontierValueDisp.SetText("0")
}

func (lp *LeftPanel) SetSearchStatsFromState(vs *ViewState) {
	if vs.Result == nil || len(vs.Result.SearchFrames) == 0 {
		lp.ClearSearchStats()
		return
	}
	visited := vs.VisitedSet()
	frontier := vs.FrontierSet()
	vStr := fmt.Sprintf("Visited: %d", len(visited))
	fStr := fmt.Sprintf("Frontier: %d", len(frontier))
	lp.VisitedLabel.SetText(vStr)
	lp.FrontierLabel.SetText(fStr)
	lp.visitedValueDisp.SetText(fmt.Sprintf("%d", len(visited)))
	lp.frontierValueDisp.SetText(fmt.Sprintf("%d", len(frontier)))
}

func (lp *LeftPanel) ClearSearchStats() {
	lp.VisitedLabel.SetText("Visited: -")
	lp.FrontierLabel.SetText("Frontier: -")
	lp.visitedValueDisp.SetText("-")
	lp.frontierValueDisp.SetText("-")
}

func (lp *LeftPanel) SetSolution(solution string) {
	lp.SolutionLabel.SetText(fmt.Sprintf("Solusi Yang Ditemukan: %s", solution))
	lp.solValueDisp.SetText(solution)
}

func (lp *LeftPanel) SetCost(cost int) {
	lp.CostLabel.SetText(fmt.Sprintf("Cost dari Solusi: %d", cost))
	lp.costValueDisp.SetText(fmt.Sprintf("%d", cost))
}

func (lp *LeftPanel) SetNoSolution() {
	lp.SolutionLabel.SetText("Solusi Yang Ditemukan: Tidak ada solusi")
	lp.CostLabel.SetText("Cost dari Solusi: -")
	lp.TimeLabel.SetText("Waktu eksekusi: -")
	lp.IterationsLabel.SetText("Banyak iterasi: -")
	lp.solValueDisp.SetText("Tidak ditemukan")
	lp.costValueDisp.SetText("-")
	lp.timeValueDisp.SetText("-")
	lp.iterValueDisp.SetText("-")
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

// UpdateSlider adjusts the slider range and position without triggering OnChanged.
func (lp *LeftPanel) UpdateSlider(current, max int) {
	lp.PlaybackSlider.Max = float64(max)
	lp.PlaybackSlider.SetValue(float64(current))
}

// SetNavEnabled enables or disables the step navigation buttons and slider.
func (lp *LeftPanel) SetNavEnabled(enabled bool) {
	if enabled {
		lp.FirstStepBtn.Enable()
		lp.PrevStepBtn.Enable()
		lp.NextStepBtn.Enable()
		lp.LastStepBtn.Enable()
		lp.PlaybackSlider.Enable()
	} else {
		lp.FirstStepBtn.Disable()
		lp.PrevStepBtn.Disable()
		lp.NextStepBtn.Disable()
		lp.LastStepBtn.Disable()
		lp.PlaybackSlider.Disable()
	}
}
