package models

type TileType int

const (
	TilePath   TileType = iota // * (jalan biasa, termasuk setelah angka dikumpulkan)
	TileWall                   // X (rintangan/batu)
	TileLava                   // L (lava, game over jika dilewati)
	TileStart                  // Z (posisi awal aktor)
	TileGoal                   // O (titik tujuan)
	TileNumber                 // 0..9 (angka yang harus dikumpulkan berurutan)
)

type MapData struct {
	Grid         [][]rune         // data input
	TileTypes    [][]TileType     // Tipe u/ petak
	Costs        [][]int          // Biaya traversal u/ petak
	Width        int              // Jumlah kolom (M)
	Height       int              // Jumlah baris (N)
	StartPos     Position         // Lokasi Z
	GoalPos      Position         // Lokasi O
	NumberPos    map[int]Position // Lokasi setiap angka 0..9 yang ada
	TotalNumbers int              // Banyak angka pada papan
	MinCost      int              // Biaya minimum di seluruh papan (untuk heuristic)
}

func (m *MapData) InBounds(p Position) bool {
	return p.X >= 0 && p.X < m.Height && p.Y >= 0 && p.Y < m.Width
}

func (m *MapData) TileAt(p Position) TileType {
	if !m.InBounds(p) {
		return TileWall
	}
	return m.TileTypes[p.X][p.Y]
}

func (m *MapData) CostAt(p Position) int {
	if !m.InBounds(p) {
		return 0
	}
	return m.Costs[p.X][p.Y]
}
