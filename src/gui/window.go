package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	content := container.NewHSplit(leftPanel.Object(), boardRenderer.Object())
	content.SetOffset(0.28)

	w.SetContent(content)

	return &MainWindow{
		App:           a,
		Window:        w,
		LeftPanel:     leftPanel,
		BoardRenderer: boardRenderer,
		State:         state,
	}
}

func (mw *MainWindow) ShowAndRun() {
	mw.Window.ShowAndRun()
}
