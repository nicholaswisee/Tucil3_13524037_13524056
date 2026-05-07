package gui

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"

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

	leftPanel.AlgorithmSelect.OnChanged = func(alg string) {
		mw.LeftPanel.SetHeuristicEnabled(alg != "UCS")
	}

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

			state.SetMap(m)
			boardRenderer.Refresh()
			mw.LeftPanel.SetStepLabel(0, 0)
			mw.LeftPanel.SetSearchStepLabel(0, 0)
			mw.LeftPanel.SetStats(0, 0)
			mw.LeftPanel.SetSolution("-")
			mw.LeftPanel.SetCost(0)
			mw.LeftPanel.VisitedLabel.SetText("Visited: -")
			mw.LeftPanel.FrontierLabel.SetText("Frontier: -")
		}, w)

		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		absInput, _ := filepath.Abs(defaultInputDir)
		if uri, err := storage.ParseURI("file://" + absInput); err == nil {
			listable, err := storage.ListerForURI(uri)
			if err == nil {
				fd.SetLocation(listable)
			}
		}
		fd.Show()
	}

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

			content := formatMapData(state.MapData)
			_, writeErr := writer.Write([]byte(content))
			if writeErr != nil {
				dialog.ShowError(writeErr, w)
			}
		}, w)

		sd.SetFileName("board.txt")
		absOutput, _ := filepath.Abs(defaultOutputDir)
		if uri, err := storage.ParseURI("file://" + absOutput); err == nil {
			listable, err := storage.ListerForURI(uri)
			if err == nil {
				sd.SetLocation(listable)
			}
		}
		sd.Show()
	}

	leftPanel.RunBtn.OnTapped = func() {
		if state.MapData == nil {
			dialog.ShowInformation("Run", "Please import a board configuration first.", w)
			return
		}

		alg := leftPanel.AlgorithmSelect.Selected
		heurStr := leftPanel.HeuristicSelect.Selected
		heurID := heuristicID(heurStr)

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
			dialog.ShowError(fmt.Errorf("unknown algorithm: %s", alg), w)
			return
		}
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
		} else {
			mw.LeftPanel.SetNoSolution()
			mw.LeftPanel.SetStepLabel(0, 0)
		}

		mw.LeftPanel.SetStats(result.TimeMs, result.NodesEval)
		mw.refreshPlayback()
		boardRenderer.Refresh()
	}

	leftPanel.PhaseToggleBtn.OnTapped = func() {
		state.TogglePhase()
		mw.LeftPanel.SetPhaseToggleLabel(state.SearchPhase)
		mw.refreshPlayback()
		boardRenderer.Refresh()
	}

	leftPanel.FirstStepBtn.OnTapped = func() {
		if state.SearchPhase {
			state.JumpToSearchStart()
		} else {
			state.JumpToStart()
		}
		mw.refreshPlayback()
	}
	leftPanel.PrevStepBtn.OnTapped = func() {
		if state.SearchPhase {
			state.SearchBackward()
		} else {
			state.StepBackward()
		}
		mw.refreshPlayback()
	}
	leftPanel.NextStepBtn.OnTapped = func() {
		if state.SearchPhase {
			state.SearchForward()
		} else {
			state.StepForward()
		}
		mw.refreshPlayback()
	}
	leftPanel.LastStepBtn.OnTapped = func() {
		if state.SearchPhase {
			state.JumpToSearchEnd()
		} else {
			state.JumpToEnd()
		}
		mw.refreshPlayback()
	}

	content := container.NewHSplit(leftPanel.Object(), boardRenderer.Object())
	content.SetOffset(0.28)

	w.SetContent(content)
	return mw
}

func (mw *MainWindow) refreshPlayback() {
	if mw.State.SearchPhase && mw.State.Result != nil {
		total := len(mw.State.Result.SearchFrames)
		if total > 0 {
			mw.LeftPanel.SetSearchStepLabel(mw.State.SearchStep, total-1)
			frame := mw.State.Result.SearchFrames[mw.State.SearchStep]
			mw.LeftPanel.SetSearchFrameStats(frame)
		} else {
			mw.LeftPanel.SetSearchStepLabel(0, 0)
			mw.LeftPanel.VisitedLabel.SetText("Visited: -")
			mw.LeftPanel.FrontierLabel.SetText("Frontier: -")
		}
	} else {
		total := 0
		if mw.State.Result != nil && mw.State.Result.Success {
			total = len(mw.State.Result.PathHistory) - 1
		}
		mw.LeftPanel.SetStepLabel(mw.State.CurrentStep, total)
	}
	mw.BoardRenderer.Refresh()
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
