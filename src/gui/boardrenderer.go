package gui

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
)

var (
	colWall   = color.NRGBA{R: 64, G: 64, B: 64, A: 255}
	colPath   = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	colLava   = color.NRGBA{R: 220, G: 20, B: 20, A: 255}
	colStart  = color.NRGBA{R: 34, G: 180, B: 34, A: 255}
	colGoal   = color.NRGBA{R: 255, G: 220, B: 50, A: 255} // yellow star
	colNumber = color.NRGBA{R: 255, G: 220, B: 50, A: 255}
	colPlayer = color.NRGBA{R: 255, G: 140, B: 20, A: 255}
	colTrail  = color.NRGBA{R: 255, G: 140, B: 20, A: 120}
	colText   = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	colBg     = color.NRGBA{R: 30, G: 30, B: 30, A: 255}
)

var numberFace font.Face

func init() {
	tt, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return
	}
	numberFace, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

type BoardRenderer struct {
	state  *ViewState
	raster *canvas.Raster
	obj    fyne.CanvasObject
}

func NewBoardRenderer(state *ViewState) *BoardRenderer {
	b := &BoardRenderer{state: state}
	b.raster = canvas.NewRaster(b.draw)
	b.obj = container.NewMax(b.raster)
	return b
}

func (b *BoardRenderer) Object() fyne.CanvasObject {
	return b.obj
}

func (b *BoardRenderer) Refresh() {
	b.raster.Refresh()
}

func (b *BoardRenderer) draw(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.NewUniform(colBg), image.Point{}, draw.Src)

	m := b.state.MapData
	if m == nil || m.Height == 0 || m.Width == 0 {
		return img
	}

	padding := 4
	availableW := w - padding*(m.Width+1)
	availableH := h - padding*(m.Height+1)
	if availableW <= 0 || availableH <= 0 {
		return img
	}

	cellW := availableW / m.Width
	cellH := availableH / m.Height
	if cellW < 2 || cellH < 2 {
		return img
	}
	cellSize := cellW
	if cellH < cellSize {
		cellSize = cellH
	}

	totalGridW := cellSize*m.Width + padding*(m.Width+1)
	totalGridH := cellSize*m.Height + padding*(m.Height+1)
	offsetX := (w - totalGridW) / 2
	offsetY := (h - totalGridH) / 2

	// 1. Draw all tiles
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			x := offsetX + padding + j*(cellSize+padding)
			y := offsetY + padding + i*(cellSize+padding)
			rect := image.Rect(x, y, x+cellSize, y+cellSize)
			tile := m.TileAt(models.Position{X: i, Y: j})
			col := tileColor(tile)
			draw.Draw(img, rect, image.NewUniform(col), image.Point{}, draw.Src)

			if tile == models.TileNumber {
				ch := m.Grid[i][j]
				drawCenteredText(img, string(ch), x+cellSize/2, y+cellSize/2, cellSize)
			}
		}
	}

	// 2. Draw yellow star on goal tile (sits on top of default/path tile)
	goalPos := m.GoalPos
	gx := offsetX + padding + goalPos.Y*(cellSize+padding) + cellSize/2
	gy := offsetY + padding + goalPos.X*(cellSize+padding) + cellSize/2
	starOuter := cellSize / 2
	if starOuter < 4 {
		starOuter = 4
	}
	starInner := starOuter / 2
	if starInner < 2 {
		starInner = 2
	}
	drawStar(img, gx, gy, starOuter, starInner, colGoal)

	// 3. Draw path trail
	if b.state.Result != nil && b.state.Result.Success {
		history := b.state.Result.PathHistory
		steps := b.state.CurrentStep
		if steps > 0 && len(history) > 1 {
			for s := 1; s <= steps && s < len(history); s++ {
				p1 := history[s-1]
				p2 := history[s]
				cx1 := offsetX + padding + p1.Y*(cellSize+padding) + cellSize/2
				cy1 := offsetY + padding + p1.X*(cellSize+padding) + cellSize/2
				cx2 := offsetX + padding + p2.Y*(cellSize+padding) + cellSize/2
				cy2 := offsetY + padding + p2.X*(cellSize+padding) + cellSize/2
				drawLine(img, cx1, cy1, cx2, cy2, colTrail)
			}
		}
	}

	// 4. Draw player circle on top of current tile
	pos := b.state.CurrentPos()
	cx := offsetX + padding + pos.Y*(cellSize+padding) + cellSize/2
	cy := offsetY + padding + pos.X*(cellSize+padding) + cellSize/2
	radius := cellSize / 3
	if radius < 2 {
		radius = 2
	}
	drawCircle(img, cx, cy, radius, colPlayer)

	return img
}

func tileColor(t models.TileType) color.Color {
	switch t {
	case models.TileWall:
		return colWall
	case models.TilePath:
		return colPath
	case models.TileLava:
		return colLava
	case models.TileStart:
		return colStart
	case models.TileGoal:
		return colPath // goal uses default tile; star drawn on top
	case models.TileNumber:
		return colNumber
	}
	return colPath
}

func drawCenteredText(img *image.RGBA, s string, cx, cy, cellSize int) {
	if cellSize < 12 || numberFace == nil {
		return
	}
	metrics := numberFace.Metrics()
	textH := (metrics.Ascent + metrics.Descent).Ceil()
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(colText),
		Face: numberFace,
	}
	bounds := d.MeasureString(s)
	textW := bounds.Ceil()
	x := cx - textW/2
	y := cy + textH/2 - metrics.Descent.Ceil()
	d.Dot = fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
	d.DrawString(s)
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	sy := 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func drawCircle(img *image.RGBA, cx, cy, r int, col color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				img.Set(cx+x, cy+y, col)
			}
		}
	}
}

func drawStar(img *image.RGBA, cx, cy, outerR, innerR int, col color.Color) {
	if outerR < 4 {
		return
	}
	poly := make([]image.Point, 10)
	for i := 0; i < 10; i++ {
		angle := float64(i)*math.Pi/5.0 - math.Pi/2.0
		r := outerR
		if i%2 == 1 {
			r = innerR
		}
		x := cx + int(float64(r)*math.Cos(angle))
		y := cy + int(float64(r)*math.Sin(angle))
		poly[i] = image.Pt(x, y)
	}

	minX, maxX := cx-outerR, cx+outerR
	minY, maxY := cy-outerR, cy+outerR
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX >= img.Bounds().Dx() {
		maxX = img.Bounds().Dx() - 1
	}
	if maxY >= img.Bounds().Dy() {
		maxY = img.Bounds().Dy() - 1
	}

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInPolygon(image.Pt(x, y), poly) {
				img.Set(x, y, col)
			}
		}
	}
}

func pointInPolygon(p image.Point, poly []image.Point) bool {
	inside := false
	for i, j := 0, len(poly)-1; i < len(poly); j, i = i, i+1 {
		pi, pj := poly[i], poly[j]
		if ((pi.Y > p.Y) != (pj.Y > p.Y)) && (p.X < (pj.X-pi.X)*(p.Y-pi.Y)/(pj.Y-pi.Y)+pi.X) {
			inside = !inside
		}
	}
	return inside
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
