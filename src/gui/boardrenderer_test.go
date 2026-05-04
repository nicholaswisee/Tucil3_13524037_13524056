package gui

import (
	"testing"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
)

func TestBoardRenderer_Creation(t *testing.T) {
	state := NewViewState()
	br := NewBoardRenderer(state)
	if br == nil {
		t.Fatal("NewBoardRenderer returned nil")
	}
	if br.Object() == nil {
		t.Error("Object() returned nil")
	}
}

func TestBoardRenderer_DrawNilState(t *testing.T) {
	state := NewViewState()
	br := NewBoardRenderer(state)
	img := br.draw(100, 100)
	if img == nil {
		t.Fatal("draw returned nil")
	}
	if img.Bounds().Dx() != 100 || img.Bounds().Dy() != 100 {
		t.Error("draw returned wrong image size")
	}
}

func TestBoardRenderer_DrawWithMap(t *testing.T) {
	state := NewViewState()
	state.MapData = &models.MapData{
		Height: 3, Width: 3,
		Grid:      [][]rune{{'X', '*', 'Z'}, {'*', 'O', '*'}, {'1', 'L', '*'}},
		TileTypes: [][]models.TileType{{models.TileWall, models.TilePath, models.TileStart}, {models.TilePath, models.TileGoal, models.TilePath}, {models.TileNumber, models.TileLava, models.TilePath}},
		StartPos: models.Position{X: 0, Y: 2},
	}
	br := NewBoardRenderer(state)
	img := br.draw(200, 200)
	if img == nil {
		t.Fatal("draw returned nil")
	}
}
