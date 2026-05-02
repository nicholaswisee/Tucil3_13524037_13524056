package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"
)

func ParseFile(path string) (*models.MapData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, fmt.Errorf("file kosong atau gagal membaca baris pertama")
	}
	firstLine := strings.Fields(scanner.Text())
	if len(firstLine) != 2 {
		return nil, fmt.Errorf("baris pertama harus berisi tepat dua angka (N M), ditemukan %d token", len(firstLine))
	}
	n, err := strconv.Atoi(firstLine[0])
	if err != nil {
		return nil, fmt.Errorf("N bukan angka valid: %w", err)
	}
	m, err := strconv.Atoi(firstLine[1])
	if err != nil {
		return nil, fmt.Errorf("M bukan angka valid: %w", err)
	}
	if n <= 0 || m <= 0 {
		return nil, fmt.Errorf("dimensi papan harus positif, ditemukan N=%d M=%d", n, m)
	}

	grid := make([][]rune, n)
	tileTypes := make([][]models.TileType, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("kekurangan baris grid: diharapkan %d baris, hanya ditemukan %d", n, i)
		}
		line := scanner.Text()
		if len(line) != m {
			return nil, fmt.Errorf("baris grid ke-%d panjangnya %d, diharapkan %d", i+1, len(line), m)
		}
		grid[i] = make([]rune, m)
		tileTypes[i] = make([]models.TileType, m)
		for j, ch := range line {
			grid[i][j] = ch
			tt, ok := parseTileType(ch)
			if !ok {
				return nil, fmt.Errorf("karakter tidak valid '%c' pada baris %d kolom %d", ch, i+1, j+1)
			}
			tileTypes[i][j] = tt
		}
	}

	// N baris cost
	costs := make([][]int, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("kekurangan baris cost: diharapkan %d baris, hanya ditemukan %d", n, i)
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) != m {
			return nil, fmt.Errorf("baris cost ke-%d memiliki %d token, diharapkan %d", i+1, len(fields), m)
		}
		costs[i] = make([]int, m)
		for j, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("cost pada baris %d kolom %d bukan angka valid: %w", i+1, j+1, err)
			}
			costs[i][j] = val
		}
	}

	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) != "" {
			return nil, fmt.Errorf("ditemukan baris tambahan setelah %d baris cost", n)
		}
	}

	mapData := &models.MapData{
		Grid:      grid,
		TileTypes: tileTypes,
		Costs:     costs,
		Width:     m,
		Height:    n,
		NumberPos: make(map[int]models.Position),
	}

	if err := validateBoard(mapData); err != nil {
		return nil, err
	}

	return mapData, nil
}

func parseTileType(ch rune) (models.TileType, bool) {
	switch ch {
	case '*':
		return models.TilePath, true
	case 'X':
		return models.TileWall, true
	case 'L':
		return models.TileLava, true
	case 'Z':
		return models.TileStart, true
	case 'O':
		return models.TileGoal, true
	default:
		if unicode.IsDigit(ch) {
			return models.TileNumber, true
		}
		return models.TilePath, false
	}
}

func validateBoard(m *models.MapData) error {
	startCount := 0
	goalCount := 0
	maxNum := -1

	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			ch := m.Grid[i][j]
			pos := models.Position{X: i, Y: j}
			switch ch {
			case 'Z':
				startCount++
				m.StartPos = pos
			case 'O':
				goalCount++
				m.GoalPos = pos
			default:
				if unicode.IsDigit(ch) {
					num := int(ch - '0')
					m.NumberPos[num] = pos
					if num > maxNum {
						maxNum = num
					}
				}
			}
		}
	}

	if startCount != 1 {
		return fmt.Errorf("papan harus memiliki tepat satu 'Z' (start), ditemukan %d", startCount)
	}
	if goalCount != 1 {
		return fmt.Errorf("papan harus memiliki tepat satu 'O' (goal), ditemukan %d", goalCount)
	}

	if maxNum >= 0 {
		m.TotalNumbers = maxNum + 1
		for i := 0; i <= maxNum; i++ {
			if _, ok := m.NumberPos[i]; !ok {
				return fmt.Errorf("angka pada papan harus berurutan mulai dari 0, namun angka %d tidak ditemukan", i)
			}
		}
	} else {
		m.TotalNumbers = 0
	}

	return nil
}
