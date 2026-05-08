package gui

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

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
	colBg       = color.NRGBA{R: 30, G: 30, B: 46, A: 255}    // #1E1E2E
	colWall     = color.NRGBA{R: 60, G: 60, B: 70, A: 255}    // wall fill
	colWallBdr  = color.NRGBA{R: 100, G: 100, B: 110, A: 255} // wall border
	colPath     = color.NRGBA{R: 184, G: 208, B: 232, A: 255} // #B8D0E8
	colStart    = color.NRGBA{R: 152, G: 195, B: 121, A: 255} // #98C379
	colGoal     = color.NRGBA{R: 229, G: 192, B: 123, A: 255} // #E5C07B
	colNumber   = color.NRGBA{R: 229, G: 192, B: 123, A: 255} // checkpoint yellow
	colNumPast  = color.NRGBA{R: 34, G: 197, B: 94, A: 255}   // checkpoint green (passed)
	colNumberTx = color.NRGBA{R: 30, G: 30, B: 46, A: 255}    // #1E1E2E
	colPlayer   = color.NRGBA{R: 209, G: 154, B: 102, A: 255} // #D19A66

	colTrail       = color.NRGBA{R: 209, G: 154, B: 102, A: 190} // orange solution path
	colSearchTrail = color.NRGBA{R: 198, G: 120, B: 221, A: 160} // purple search path

	colVisitedDot = color.NRGBA{R: 200, G: 200, B: 220, A: 180}
	colFrontier   = color.NRGBA{R: 97, G: 175, B: 239, A: 200}
	colCurrent    = color.NRGBA{R: 97, G: 175, B: 239, A: 255}
	colVisitedOvl = color.NRGBA{R: 0, G: 0, B: 30, A: 60} // dark overlay for visited tiles
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

	animFracRow float64
	animFracCol float64
	animating   bool
	currentAnim *fyne.Animation

	lavaAlpha uint8
	lavaAnim  *fyne.Animation

	OnAnimStart func()
	OnAnimEnd   func()
}

func NewBoardRenderer(state *ViewState) *BoardRenderer {
	b := &BoardRenderer{
		state:     state,
		lavaAlpha: 255,
	}
	b.raster = canvas.NewRaster(b.draw)
	b.obj = container.NewMax(b.raster)

	b.lavaAnim = fyne.NewAnimation(time.Second, func(f float32) {
		b.lavaAlpha = uint8(255 - int(75*f))
		b.raster.Refresh()
	})
	b.lavaAnim.RepeatCount = fyne.AnimationRepeatForever
	b.lavaAnim.AutoReverse = true
	b.lavaAnim.Start()

	return b
}

func (b *BoardRenderer) Object() fyne.CanvasObject {
	return b.obj
}

func (b *BoardRenderer) Refresh() {
	b.raster.Refresh()
}

func (b *BoardRenderer) AnimateSlide(positions []models.Position, duration time.Duration, onComplete func()) {
	if b.currentAnim != nil {
		b.currentAnim.Stop()
		b.currentAnim = nil
	}
	n := len(positions)
	if n < 2 {
		b.animating = false
		if onComplete != nil {
			onComplete()
		}
		return
	}

	b.animating = true
	b.animFracRow = float64(positions[0].X)
	b.animFracCol = float64(positions[0].Y)

	if b.OnAnimStart != nil {
		b.OnAnimStart()
	}

	total := float64(n - 1)
	b.currentAnim = fyne.NewAnimation(duration, func(f float32) {
		ft := float64(f) * total
		idx := int(ft)
		if idx >= n-1 {
			idx = n - 2
		}
		frac := ft - float64(idx)
		r1 := float64(positions[idx].X)
		c1 := float64(positions[idx].Y)
		r2 := float64(positions[idx+1].X)
		c2 := float64(positions[idx+1].Y)
		b.animFracRow = r1 + frac*(r2-r1)
		b.animFracCol = c1 + frac*(c2-c1)
		b.raster.Refresh()
	})
	b.currentAnim.Curve = fyne.AnimationEaseInOut
	b.currentAnim.Start()

	go func() {
		time.Sleep(duration + 20*time.Millisecond)
		b.animating = false
		b.animFracRow = float64(positions[n-1].X)
		b.animFracCol = float64(positions[n-1].Y)
		fyne.Do(func() {
			if b.OnAnimEnd != nil {
				b.OnAnimEnd()
			}
			if onComplete != nil {
				onComplete()
			}
			b.raster.Refresh()
		})
	}()
}

func (b *BoardRenderer) StopAnimation() {
	if b.currentAnim != nil {
		b.currentAnim.Stop()
		b.currentAnim = nil
	}
	b.animating = false
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

	cellCenter := func(p models.Position) (int, int) {
		x := offsetX + padding + p.Y*(cellSize+padding) + cellSize/2
		y := offsetY + padding + p.X*(cellSize+padding) + cellSize/2
		return x, y
	}

	fracCenter := func(row, col float64) (int, int) {
		x := float64(offsetX+padding) + col*(float64(cellSize)+float64(padding)) + float64(cellSize)/2
		y := float64(offsetY+padding) + row*(float64(cellSize)+float64(padding)) + float64(cellSize)/2
		return int(x), int(y)
	}

	var visited map[models.Position]bool
	var checkpointsPassed map[int]bool
	if b.state.SearchPhase {
		visited = b.state.VisitedSet()
	} else {
		checkpointsPassed = b.state.CheckpointsPassed()
	}

	// 1. Draw all tiles
	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			x := offsetX + padding + j*(cellSize+padding)
			y := offsetY + padding + i*(cellSize+padding)
			pos := models.Position{X: i, Y: j}
			tile := m.TileAt(pos)

			col := b.tileColorAnimated(tile, pos, checkpointsPassed)

			if b.state.SearchPhase && visited != nil && visited[pos] {
				col = blendNRGBA(col, colVisitedOvl)
			}

			if radius > 0 {
				drawRoundedRect(img, x, y, cellSize, cellSize, radius, col)
			} else {
				rect := image.Rect(x, y, x+cellSize, y+cellSize)
				draw.Draw(img, rect, image.NewUniform(col), image.Point{}, draw.Src)
			}

			if tile == models.TileWall && cellSize >= 6 {
				drawRectBorder(img, x, y, cellSize, cellSize, 2, colWallBdr)
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

	// 2. SEARCH PHASE
	if b.state.SearchPhase && b.state.Result != nil && len(b.state.Result.SearchFrames) > 0 {

		frontier := b.state.FrontierSet()
		for p := range frontier {
			cx, cy := cellCenter(p)
			circleR := cellSize / 5
			if circleR < 2 {
				circleR = 2
			}
			drawCircleOutline(img, cx, cy, circleR, colFrontier)
		}

		for p := range visited {
			cx, cy := cellCenter(p)
			dotR := cellSize / 8
			if dotR < 1 {
				dotR = 1
			}
			drawCircle(img, cx, cy, dotR, colVisitedDot)
		}

		frame := b.state.Result.SearchFrames[b.state.SearchStep]

		if len(frame.PathToNode) > 1 {
			for s := 1; s < len(frame.PathToNode); s++ {
				p1 := frame.PathToNode[s-1]
				p2 := frame.PathToNode[s]
				cx1, cy1 := cellCenter(p1)
				cx2, cy2 := cellCenter(p2)
				drawThickLine(img, cx1, cy1, cx2, cy2, 2, colSearchTrail)
				drawArrowhead(img, cx1, cy1, cx2, cy2, 8, colSearchTrail)
			}
		}

		cx, cy := cellCenter(frame.Current)
		r := cellSize / 3
		if r < 2 {
			r = 2
		}
		drawCircle(img, cx, cy, r, colCurrent)
	}

	// 3. Goal star (always drawn)
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

	// 4. SOLUTION PHASE path
	if !b.state.SearchPhase && b.state.Result != nil && b.state.Result.Success {
		history := b.state.Result.PathHistory
		steps := b.state.CurrentStep
		if steps > 0 && len(history) > 1 {
			for s := 1; s <= steps && s < len(history); s++ {
				p1 := history[s-1]
				p2 := history[s]
				cx1, cy1 := cellCenter(p1)
				cx2, cy2 := cellCenter(p2)
				drawThickLine(img, cx1, cy1, cx2, cy2, 4, colTrail)
				drawArrowhead(img, cx1, cy1, cx2, cy2, 10, colTrail)
			}
		}
	}

	// 5. Player token (solution phase)
	if !b.state.SearchPhase {
		var pcx, pcy int
		if b.animating {
			pcx, pcy = fracCenter(b.animFracRow, b.animFracCol)
		} else {
			pos := b.state.CurrentPos()
			pcx, pcy = cellCenter(pos)
		}
		radiusP := cellSize / 3
		if radiusP < 2 {
			radiusP = 2
		}
		drawCircle(img, pcx, pcy, radiusP, colPlayer)
		outlineColor := color.NRGBA{R: 229, G: 192, B: 123, A: 255}
		drawCircleOutline(img, pcx, pcy, radiusP+2, outlineColor)
	}

	return img
}

func (b *BoardRenderer) tileColorAnimated(t models.TileType, pos models.Position, checkpointsPassed map[int]bool) color.NRGBA {
	switch t {
	case models.TileWall:
		return colWall
	case models.TilePath:
		return colPath
	case models.TileLava:
		return color.NRGBA{R: 220, G: 60, B: 60, A: b.lavaAlpha}
	case models.TileStart:
		return colStart
	case models.TileGoal:
		return colPath
	case models.TileNumber:
		if checkpointsPassed != nil {
			m := b.state.MapData
			if m != nil {
				numIdx := int(m.Grid[pos.X][pos.Y] - '0')
				if checkpointsPassed[numIdx] {
					return colNumPast // green badge
				}
			}
		}
		return colNumber
	}
	return colPath
}

func blendNRGBA(base, overlay color.NRGBA) color.NRGBA {
	a := float64(overlay.A) / 255.0
	return color.NRGBA{
		R: uint8(float64(overlay.R)*a + float64(base.R)*(1-a)),
		G: uint8(float64(overlay.G)*a + float64(base.G)*(1-a)),
		B: uint8(float64(overlay.B)*a + float64(base.B)*(1-a)),
		A: base.A,
	}
}

func tileColor(t models.TileType) color.Color {
	switch t {
	case models.TileWall:
		return colWall
	case models.TilePath:
		return colPath
	case models.TileLava:
		return color.NRGBA{R: 220, G: 60, B: 60, A: 255}
	case models.TileStart:
		return colStart
	case models.TileGoal:
		return colPath
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

func drawRectBorder(img *image.RGBA, x, y, w, h, bw int, col color.Color) {
	for dy := 0; dy < bw && dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
	for dy := h - bw; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
	for dy := bw; dy < h-bw; dy++ {
		for dx := 0; dx < bw && dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
	for dy := bw; dy < h-bw; dy++ {
		for dx := w - bw; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
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

func drawThinLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
	dx := x1 - x0
	dy := y1 - y0
	steps := abs(dx)
	if abs(dy) > steps {
		steps = abs(dy)
	}
	if steps == 0 {
		img.Set(x0, y0, col)
		return
	}
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x0 + int(float64(dx)*t)
		y := y0 + int(float64(dy)*t)
		img.Set(x, y, col)
	}
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

func drawArrowhead(img *image.RGBA, x0, y0, x1, y1, size int, col color.Color) {
	dx := float64(x1 - x0)
	dy := float64(y1 - y0)
	length := math.Sqrt(dx*dx + dy*dy)
	if length < 1 {
		return
	}
	ux := dx / length
	uy := dy / length

	cos30 := math.Cos(math.Pi / 6)
	sin30 := math.Sin(math.Pi / 6)

	lx := float64(x1) + float64(size)*(-ux*cos30+uy*sin30)
	ly := float64(y1) + float64(size)*(-ux*sin30-uy*cos30)
	rx := float64(x1) + float64(size)*(-ux*cos30-uy*sin30)
	ry := float64(y1) + float64(size)*(ux*sin30-uy*cos30)

	drawThickLine(img, x1, y1, int(lx), int(ly), 2, col)
	drawThickLine(img, x1, y1, int(rx), int(ry), 2, col)
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
			if d2 <= r*r && d2 > (r-1)*(r-1) {
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
