package gui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/parser"
	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/solver"
)

var (
	defaultInputDir  = "../test/input"
	defaultOutputDir = "../test/output"
)

type MainWindow struct {
	App           fyne.App
	Window        fyne.Window
	LeftPanel     *LeftPanel
	BoardRenderer *BoardRenderer
	State         *ViewState

	// Step-info label shown above the board (centred)
	stepOverlay *widget.Label

	// Playback goroutine control
	playDoneCh chan struct{}
	isPlaying  bool

	// Loaded filename for window title
	loadedFile string

	// Slider OnChanged guard: skip re-entrant updates triggered by UpdateSlider
	sliderUpdating bool
}

func NewMainWindow() *MainWindow {
	a := app.New()
	w := a.NewWindow("Ice Sliding Puzzle Solver")
	w.Resize(fyne.NewSize(900, 600))

	state := NewViewState()
	leftPanel := NewLeftPanel()
	boardRenderer := NewBoardRenderer(state)

	mw := &MainWindow{
		App:           a,
		Window:        w,
		LeftPanel:     leftPanel,
		BoardRenderer: boardRenderer,
		State:         state,
	}

	// Step-info overlay (centred above the board)
	mw.stepOverlay = widget.NewLabelWithStyle(
		"Step 0 / 0",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	// Board area: step label on top, canvas below
	boardArea := container.NewBorder(
		container.NewCenter(mw.stepOverlay),
		nil, nil, nil,
		boardRenderer.Object(),
	)

	// Animation callbacks
	boardRenderer.OnAnimStart = func() {
		fyne.Do(func() { leftPanel.SetNavEnabled(false) })
	}
	boardRenderer.OnAnimEnd = func() {
		fyne.Do(func() { leftPanel.SetNavEnabled(true) })
	}

	// Algorithm selection
	leftPanel.AlgorithmSelect.OnChanged = func(alg string) {
		enabled := alg != "UCS"
		mw.LeftPanel.SetHeuristicEnabled(enabled)
		if !enabled {
			mw.LeftPanel.HeuristicSelect.PlaceHolder = "Not required for UCS"
			mw.LeftPanel.HeuristicSelect.Refresh()
		}
	}

	// Import
	leftPanel.ImportBtn.OnTapped = func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			m, err := parser.ParseFile(reader.URI().Path())
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			mw.stopPlayback()
			state.SetMap(m)
			boardRenderer.StopAnimation()
			boardRenderer.Refresh()

			mw.loadedFile = filepath.Base(reader.URI().Path())
			w.SetTitle("Ice Sliding Puzzle Solver — " + mw.loadedFile)

			mw.LeftPanel.SetStepLabel(0, 0)
			mw.LeftPanel.SetSearchStepLabel(0, 0)
			mw.LeftPanel.SetStats(0, 0)
			mw.LeftPanel.SetSolution("-")
			mw.LeftPanel.SetCost(0)
			mw.LeftPanel.ClearSearchStats()
			mw.LeftPanel.UpdateSlider(0, 0)
			mw.stepOverlay.SetText("Step 0 / 0")
		}, w)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		absInput, _ := filepath.Abs(defaultInputDir)
		if uri, err := storage.ParseURI("file://" + absInput); err == nil {
			if listable, err := storage.ListerForURI(uri); err == nil {
				fd.SetLocation(listable)
			}
		}
		fd.Show()
	}

	// Export
	leftPanel.ExportBtn.OnTapped = func() {
		if state.MapData == nil {
			dialog.ShowInformation("Export", "No board loaded to export.", w)
			return
		}
		sd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			content := formatExportResult(state.MapData, state.Result)
			if _, writeErr := writer.Write([]byte(content)); writeErr != nil {
				dialog.ShowError(writeErr, w)
			}
		}, w)
		sd.SetFileName("result.txt")
		absOutput, _ := filepath.Abs(defaultOutputDir)
		if uri, err := storage.ParseURI("file://" + absOutput); err == nil {
			if listable, err := storage.ListerForURI(uri); err == nil {
				sd.SetLocation(listable)
			}
		}
		sd.Show()
	}

	// Run / Solve
	leftPanel.RunBtn.OnTapped = func() {
		if state.MapData == nil {
			dialog.ShowInformation("Run", "Please import a board configuration first.", w)
			return
		}

		mw.stopPlayback()
		boardRenderer.StopAnimation()
		leftPanel.RunBtn.Disable()
		leftPanel.ProgressBar.Show()

		alg := leftPanel.AlgorithmSelect.Selected
		heurStr := leftPanel.HeuristicSelect.Selected
		heurID := heuristicID(heurStr)

		go func() {
			var result *models.SolverResult
			var err error
			switch alg {
			case "UCS":
				s := &solver.UCSSolver{}
				result, err = s.Solve(state.MapData)
			case "GBFS":
				s := &solver.GBFSSolver{HeuristicID: heurID}
				result, err = s.Solve(state.MapData)
			case "A*":
				s := &solver.AStarSolver{HeuristicID: heurID}
				result, err = s.Solve(state.MapData)
			case "IDA*":
				s := &solver.IDAStarSolver{HeuristicID: heurID}
				result, err = s.Solve(state.MapData)
			default:
				err = fmt.Errorf("unknown algorithm: %s", alg)
			}

			fyne.Do(func() {
				leftPanel.ProgressBar.Hide()
				leftPanel.RunBtn.Enable()

				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				state.SetResult(result)
				mw.LeftPanel.SetPhaseToggleLabel(true)

				if result.Success {
					solutionStr := pathToString(result.Path)
					mw.LeftPanel.SetSolution(solutionStr)
					mw.LeftPanel.SetCost(result.TotalCost)
					totalSteps := len(result.PathHistory) - 1
					mw.LeftPanel.SetStepLabel(0, totalSteps)
					mw.LeftPanel.UpdateSlider(0, totalSteps)
				} else {
					mw.LeftPanel.SetNoSolution()
					mw.LeftPanel.SetStepLabel(0, 0)
					mw.LeftPanel.UpdateSlider(0, 0)
					dialog.ShowInformation(
						"No Solution",
						"Tidak ada solusi yang ditemukan untuk konfigurasi ini.",
						w,
					)
				}

				mw.LeftPanel.SetStats(result.TimeMs, result.NodesEval)
				mw.refreshPlayback()
				boardRenderer.Refresh()
			})
		}()
	}

	// Phase toggle
	leftPanel.PhaseToggleBtn.OnTapped = func() {
		mw.stopPlayback()
		boardRenderer.StopAnimation()
		state.TogglePhase()
		mw.LeftPanel.SetPhaseToggleLabel(state.SearchPhase)
		mw.refreshPlayback()
		boardRenderer.Refresh()
	}

	// Step buttons
	leftPanel.FirstStepBtn.OnTapped = func() {
		mw.stopPlayback()
		if state.SearchPhase {
			state.JumpToSearchStart()
		} else {
			state.JumpToStart()
		}
		mw.refreshPlayback()
	}
	leftPanel.PrevStepBtn.OnTapped = func() {
		mw.stopPlayback()
		if state.SearchPhase {
			state.SearchBackward()
		} else {
			state.StepBackward()
		}
		mw.refreshPlayback()
	}
	leftPanel.NextStepBtn.OnTapped = func() {
		mw.stopPlayback()
		prevPos := state.CurrentPos()
		var advanced bool
		if state.SearchPhase {
			advanced = state.SearchForward()
		} else {
			advanced = state.StepForward()
		}
		mw.refreshPlayback()
		if advanced && !state.SearchPhase {
			mw.triggerSlideAnimation(prevPos, state.CurrentPos(), 1.0)
		}
	}
	leftPanel.LastStepBtn.OnTapped = func() {
		mw.stopPlayback()
		if state.SearchPhase {
			state.JumpToSearchEnd()
		} else {
			state.JumpToEnd()
		}
		mw.refreshPlayback()
	}

	// Slider
	leftPanel.PlaybackSlider.OnChanged = func(v float64) {
		if mw.sliderUpdating {
			return
		}
		mw.stopPlayback()
		boardRenderer.StopAnimation()
		step := int(v)
		if state.SearchPhase {
			state.JumpToSearchFrame(step)
		} else {
			state.JumpToStep(step)
		}
		mw.refreshPlayback()
	}

	// Play / Pause
	leftPanel.PlayPauseBtn.OnTapped = func() {
		if mw.isPlaying {
			mw.stopPlayback()
		} else {
			mw.startPlayback()
		}
	}

	// Speed select
	leftPanel.SpeedSelect.OnChanged = func(_ string) {
		if mw.isPlaying {
			// Restart ticker with new interval
			mw.stopPlayback()
			mw.startPlayback()
		}
	}

	content := container.NewHSplit(leftPanel.Object(), boardArea)
	content.SetOffset(0.28)
	w.SetContent(content)

	return mw
}

// refreshPlayback updates all playback UI elements to reflect the current state.
func (mw *MainWindow) refreshPlayback() {
	if mw.State.SearchPhase && mw.State.Result != nil {
		total := len(mw.State.Result.SearchFrames)
		if total > 0 {
			cur := mw.State.SearchStep
			mw.LeftPanel.SetSearchStepLabel(cur, total-1)
			mw.LeftPanel.SetSearchStatsFromState(mw.State)
			mw.updateSlider(cur, total-1)
			mw.stepOverlay.SetText(fmt.Sprintf("Search Step %d / %d", cur, total-1))
		} else {
			mw.LeftPanel.SetSearchStepLabel(0, 0)
			mw.LeftPanel.ClearSearchStats()
			mw.updateSlider(0, 0)
			mw.stepOverlay.SetText("Search Step 0 / 0")
		}
	} else {
		total := 0
		if mw.State.Result != nil && mw.State.Result.Success {
			total = len(mw.State.Result.PathHistory) - 1
		}
		cur := mw.State.CurrentStep
		mw.LeftPanel.SetStepLabel(cur, total)
		mw.updateSlider(cur, total)
		mw.stepOverlay.SetText(fmt.Sprintf("Step %d / %d", cur, total))
	}
	mw.BoardRenderer.Refresh()
}

// updateSlider sets the slider value without firing OnChanged.
func (mw *MainWindow) updateSlider(current, max int) {
	mw.sliderUpdating = true
	mw.LeftPanel.UpdateSlider(current, max)
	mw.sliderUpdating = false
}

// triggerSlideAnimation computes the intermediate slide path and starts the animation.
func (mw *MainWindow) triggerSlideAnimation(from, to models.Position, speedFactor float64) {
	if from == to {
		return
	}
	base := 250 * time.Millisecond
	dur := time.Duration(float64(base) / speedFactor)
	if dur < 30*time.Millisecond {
		dur = 30 * time.Millisecond
	}
	path := append([]models.Position{from}, slidePositions(from, to)...)
	mw.BoardRenderer.AnimateSlide(path, dur, nil)
}

func speedToInterval(s string) time.Duration {
	switch s {
	case "0.5×":
		return 1200 * time.Millisecond
	case "2×":
		return 300 * time.Millisecond
	case "4×":
		return 150 * time.Millisecond
	default: // "1×"
		return 600 * time.Millisecond
	}
}

func speedToAnimFactor(s string) float64 {
	switch s {
	case "0.5×":
		return 0.5
	case "2×":
		return 2.0
	case "4×":
		return 4.0
	default:
		return 1.0
	}
}

func (mw *MainWindow) startPlayback() {
	if mw.isPlaying {
		return
	}
	mw.isPlaying = true
	mw.LeftPanel.PlayPauseBtn.SetText("⏸ Pause")

	doneCh := make(chan struct{})
	mw.playDoneCh = doneCh

	interval := speedToInterval(mw.LeftPanel.SpeedSelect.Selected)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-doneCh:
				return
			case <-ticker.C:
				if mw.BoardRenderer.animating {
					continue
				}
				fyne.Do(func() {
					prevPos := mw.State.CurrentPos()
					var advanced bool
					if mw.State.SearchPhase {
						advanced = mw.State.SearchForward()
					} else {
						advanced = mw.State.StepForward()
					}
					if !advanced {
						mw.stopPlayback()
						return
					}
					mw.refreshPlayback()
					if !mw.State.SearchPhase {
						factor := speedToAnimFactor(mw.LeftPanel.SpeedSelect.Selected)
						mw.triggerSlideAnimation(prevPos, mw.State.CurrentPos(), factor)
					}
				})
			}
		}
	}()
}

func (mw *MainWindow) stopPlayback() {
	if !mw.isPlaying {
		return
	}
	mw.isPlaying = false
	mw.LeftPanel.PlayPauseBtn.SetText("▶ Play")
	if mw.playDoneCh != nil {
		close(mw.playDoneCh)
		mw.playDoneCh = nil
	}
}

func (mw *MainWindow) ShowAndRun() {
	mw.Window.ShowAndRun()
}

func heuristicID(heurStr string) int {
	switch heurStr {
	case "Heuristic 2":
		return 2
	case "Heuristic 3":
		return 3
	default:
		return 1
	}
}

func pathToString(path []models.MoveRecord) string {
	var sb strings.Builder
	for _, mr := range path {
		sb.WriteString(mr.Direction.DirectionName())
	}
	return sb.String()
}

func formatMapData(m *models.MapData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", m.Height, m.Width))
	for _, row := range m.Grid {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
	}
	for _, row := range m.Costs {
		for j, c := range row {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", c))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// formatExportResult formats the output file with board configuration at final position,
// execution time, and number of iterations.
func formatExportResult(m *models.MapData, result *models.SolverResult) string {
	var sb strings.Builder
	
	// Write board dimensions
	sb.WriteString(fmt.Sprintf("%d %d\n", m.Height, m.Width))
	
	// Create a copy of the grid to show final position
	finalGrid := make([][]rune, m.Height)
	for i := range m.Grid {
		finalGrid[i] = make([]rune, m.Width)
		copy(finalGrid[i], m.Grid[i])
	}
	
	// If there's a solution, mark the final position with 'A' (Actor)
	if result != nil && result.Success && len(result.PathHistory) > 0 {
		finalPos := result.PathHistory[len(result.PathHistory)-1]
		// Only mark if it's not already the goal
		if finalGrid[finalPos.X][finalPos.Y] != 'O' {
			finalGrid[finalPos.X][finalPos.Y] = 'A'
		}
	}
	
	// Write the grid
	for _, row := range finalGrid {
		sb.WriteString(string(row))
		sb.WriteByte('\n')
	}
	
	// Write costs
	for _, row := range m.Costs {
		for j, c := range row {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", c))
		}
		sb.WriteByte('\n')
	}
	
	// Add empty line before statistics
	sb.WriteByte('\n')
	
	// Write execution time and iterations
	if result != nil {
		sb.WriteString(fmt.Sprintf("Waktu eksekusi: %d ms\n", result.TimeMs))
		sb.WriteString(fmt.Sprintf("Banyak iterasi yang dilakukan: %d\n", result.NodesEval))
	} else {
		sb.WriteString("Waktu eksekusi: -\n")
		sb.WriteString("Banyak iterasi yang dilakukan: -\n")
	}
	
	return sb.String()
}
