package solver

import "github.com/nicholaswisee/Tucil3_13524037_13524056/core/models"

func extractPositionsFromMap(m map[models.StateKey]bool) []models.Position {
	res := make([]models.Position, 0, len(m))
	for k := range m {
		res = append(res, k.Pos)
	}
	return res
}

func extractPositionsFromCostMap(m map[models.StateKey]int) []models.Position {
	res := make([]models.Position, 0, len(m))
	for k := range m {
		res = append(res, k.Pos)
	}
	return res
}

func extractPositionsFromPQ(pq PriorityQueue) []models.Position {
	res := make([]models.Position, 0, len(pq))
	for _, node := range pq {
		res = append(res, node.State.Pos)
	}
	return res
}
