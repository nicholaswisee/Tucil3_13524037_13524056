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
	colBg       = color.NRGBA{R: 30, G: 30, B: 46, A: 255}   // #1E1E2E
	colWall     = color.NRGBA{R: 59, G: 59, B: 79, A: 255}   // #3B3B4F
	colPath     = color.NRGBA{R: 184, G: 208, B: 232, A: 255} // #B8D0E8
	colLava     = color.NRGBA{R: 224, G: 108, B: 117, A: 255} // #E06C75
	colStart    = color.NRGBA{R: 152, G: 195, B: 121, A: 255} // #98C379
	colGoal     = color.NRGBA{R: 229, G: 192, B: 123, A: 255} // #E5C07B
	colNumber   = color.NRGBA{R: 229, G: 192, B: 123, A: 255} // #E5C07B
	colNumberTx = color.NRGBA{R: 30, G: 30, B: 46, A: 255}   // #1E1E2E
	colPlayer   = color.NRGBA{R: 209, G: 154, B: 102, A: 255} // #D19A66
	colTrail    = color.NRGBA{R: 209, G: 154, B: 102, A: 153} // #D19A66 @ 60%
	colVisited  = color.NRGBA{R: 198, G: 120, B: 221, A: 102} // #C678DD @ 40%
	colFrontier = color.NRGBA{R: 97, G: 175, B: 239, A: 77}   // #61AFEF @ 30%
	colCurrent  = color.NRGBA{R: 97, G: 175, B: 239, A: 255}  // #61AFEF
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

	padding := 8
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

	radius := 4
	if cellSize < 10 {
		radius = 2
	}
	if cellSize < 6 {
		radius = 0
	}

	// 1. Draw all tiles
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			x := offsetX + padding + j*(cellSize+padding)
			y := offsetY + padding + i*(cellSize+padding)
			tile := m.TileAt(models.Position{X: i, Y: j})
			col := tileColor(tile)
			if radius > 0 {
				drawRoundedRect(img, x, y, cellSize, cellSize, radius, col)
			} else {
				rect := image.Rect(x, y, x+cellSize, y+cellSize)
				draw.Draw(img, rect, image.NewUniform(col), image.Point{}, draw.Src)
			}

			if tile == models.TileNumber {
				ch := m.Grid[i][j]
				size := cellSize * 60 / 100
				if size < 8 {
					size = 8
				}
				drawCenteredText(img, string(ch), x+cellSize/2, y+cellSize/2, size)
			}
		}
	}

	// 2. Search phase overlays
	if b.state.SearchPhase && b.state.Result != nil && len(b.state.Result.SearchFrames) > 0 {
		frame := b.state.Result.SearchFrames[b.state.SearchStep]
		for _, p := range frame.Visited {
			x := offsetX + padding + p.Y*(cellSize+padding)
			y := offsetY + padding + p.X*(cellSize+padding)
			if radius > 0 {
				drawRoundedRect(img, x, y, cellSize, cellSize, radius, colVisited)
			} else {
				rect := image.Rect(x, y, x+cellSize, y+cellSize)
				draw.Draw(img, rect, image.NewUniform(colVisited), image.Point{}, draw.Src)
			}
		}
		for _, p := range frame.Frontier {
			x := offsetX + padding + p.Y*(cellSize+padding)
			y := offsetY + padding + p.X*(cellSize+padding)
			if radius > 0 {
				drawRoundedRect(img, x, y, cellSize, cellSize, radius, colFrontier)
			} else {
				rect := image.Rect(x, y, x+cellSize, y+cellSize)
				draw.Draw(img, rect, image.NewUniform(colFrontier), image.Point{}, draw.Src)
			}
		}
		cx := offsetX + padding + frame.Current.Y*(cellSize+padding) + cellSize/2
		cy := offsetY + padding + frame.Current.X*(cellSize+padding) + cellSize/2
		r := cellSize / 2
		if r < 2 {
			r = 2
		}
		drawCircle(img, cx, cy, r, colCurrent)
	}

	// 3. Goal star
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

	// 4. Path trail (solution phase only)
	if !b.state.SearchPhase && b.state.Result != nil && b.state.Result.Success {
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
				drawThickLine(img, cx1, cy1, cx2, cy2, 3, colTrail)
			}
		}
	}

	// 5. Player token (solution phase only)
	if !b.state.SearchPhase {
		pos := b.state.CurrentPos()
		cx := offsetX + padding + pos.Y*(cellSize+padding) + cellSize/2
		cy := offsetY + padding + pos.X*(cellSize+padding) + cellSize/2
		radiusP := cellSize / 3
		if radiusP < 2 {
			radiusP = 2
		}
		drawCircle(img, cx, cy, radiusP, colPlayer)
		outlineColor := color.NRGBA{R: 229, G: 192, B: 123, A: 255}
		drawCircleOutline(img, cx, cy, radiusP+2, outlineColor)
	}

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
		return colPath // goal uses path tile; star drawn on top
	case models.TileNumber:
		return colNumber
	}
	return colPath
}

func drawRoundedRect(img *image.RGBA, x, y, w, h, r int, col color.Color) {
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			px := x + dx
			py := y + dy
			corner := false
			cx, cy := -1, -1
			if dx < r && dy < r {
				corner = true
				cx, cy = x+r, y+r
			} else if dx >= w-r && dy < r {
				corner = true
				cx, cy = x+w-r-1, y+r
			} else if dx < r && dy >= h-r {
				corner = true
				cx, cy = x+r, y+h-r-1
			} else if dx >= w-r && dy >= h-r {
				corner = true
				cx, cy = x+w-r-1, y+h-r-1
			}
			if corner {
				d2 := (px-cx)*(px-cx) + (py-cy)*(py-cy)
				if d2 > r*r {
					continue
				}
			}
			img.Set(px, py, col)
		}
	}
}

func drawCenteredText(img *image.RGBA, s string, cx, cy, fontSize int) {
	if fontSize < 8 || numberFace == nil {
		return
	}
	var face font.Face = numberFace
	if fontSize != 20 {
		tt, err := opentype.Parse(gomono.TTF)
		if err == nil {
			f, err := opentype.NewFace(tt, &opentype.FaceOptions{
				Size:    float64(fontSize),
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err == nil {
				face = f
			}
		}
	}
	metrics := face.Metrics()
	textH := (metrics.Ascent + metrics.Descent).Ceil()
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(colNumberTx),
		Face: face,
	}
	bounds := d.MeasureString(s)
	textW := bounds.Ceil()
	x := cx - textW/2
	y := cy + textH/2 - metrics.Descent.Ceil()
	d.Dot = fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
	d.DrawString(s)
}

func drawThickLine(img *image.RGBA, x0, y0, x1, y1, thickness int, col color.Color) {
	dx := x1 - x0
	dy := y1 - y0
	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}
	if steps == 0 {
		drawCircle(img, x0, y0, thickness/2, col)
		return
	}
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x0 + int(float64(dx)*t)
		y := y0 + int(float64(dy)*t)
		drawCircle(img, x, y, thickness/2, col)
	}
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

func drawCircleOutline(img *image.RGBA, cx, cy, r int, col color.Color) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			d2 := x*x + y*y
			if d2 <= r*r && d2 > (r-2)*(r-2) {
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
