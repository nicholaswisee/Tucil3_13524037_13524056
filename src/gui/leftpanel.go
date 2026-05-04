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
	StepLabel       *widget.Label
	TimeLabel       *widget.Label
	IterationsLabel *widget.Label
}

func NewLeftPanel() *LeftPanel {
	lp := &LeftPanel{}

	lp.AlgorithmSelect = widget.NewSelect([]string{"UCS", "GBFS", "A*"}, func(string) {})
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

	lp.StepLabel = widget.NewLabel("Step 0 / 0")
	lp.TimeLabel = widget.NewLabel("Time: 0 ms")
	lp.IterationsLabel = widget.NewLabel("Iterations: 0")

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
		widget.NewLabel("Playback"),
		playbackRow,
		widget.NewSeparator(),
		widget.NewLabel("Stats"),
		lp.TimeLabel,
		lp.IterationsLabel,
	)

	return lp
}

func (lp *LeftPanel) Object() fyne.CanvasObject {
	return lp.container
}

func (lp *LeftPanel) SetStats(timeMs int64, iterations int) {
	lp.TimeLabel.SetText(fmt.Sprintf("Time: %d ms", timeMs))
	lp.IterationsLabel.SetText(fmt.Sprintf("Iterations: %d", iterations))
}

func (lp *LeftPanel) SetStepLabel(current, total int) {
	lp.StepLabel.SetText(fmt.Sprintf("Step %d / %d", current, total))
}

func (lp *LeftPanel) SetHeuristicEnabled(enabled bool) {
	if enabled {
		lp.HeuristicSelect.Enable()
	} else {
		lp.HeuristicSelect.Disable()
	}
}
