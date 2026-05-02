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
