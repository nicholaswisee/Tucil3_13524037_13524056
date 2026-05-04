package models

type GameState struct {
	Pos     Position // posisi aktor setelah slide berhenti
	NextNum int      // angka berikutnya yang harus dilewati, -1 kalo sisa
}

func (s *GameState) IsGoal(m *MapData) bool {
	return s.NextNum == -1 && s.Pos == m.GoalPos
}

type MoveRecord struct {
	Direction Direction
	NewPos    Position
	MoveCost  int
}

type StateKey struct {
	Pos     Position
	NextNum int
}

func (s *GameState) GetKey() StateKey {
	return StateKey{Pos: s.Pos, NextNum: s.NextNum}
}

func (s *GameState) Slide(m *MapData, d Direction) (GameState, int, bool) {
	currPos := s.Pos
	cost := 0
	nextNum := s.NextNum

	for {
		nextPos := currPos.Add(d.DxDy())

		if !m.InBounds(nextPos) {
			return GameState{}, 0, false
		}

		tile := m.TileAt(nextPos)

		if tile == TileWall {
			break
		}

		if tile == TileLava {
			return GameState{}, 0, false
		}

		currPos = nextPos
		cost += m.CostAt(currPos)

		if tile == TileNumber {
			ch := m.Grid[currPos.X][currPos.Y]
			numInjected := int(ch - '0')

			if numInjected == nextNum {
				nextNum++
				if nextNum >= m.TotalNumbers {
					nextNum = -1
				}
			} else if numInjected > nextNum {
				return GameState{}, 0, false
			}
		}
	}

	if currPos == s.Pos {
		return GameState{}, 0, false
	}

	return GameState{Pos: currPos, NextNum: nextNum}, cost, true
}
