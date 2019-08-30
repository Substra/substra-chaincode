package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboardSort(t *testing.T) {
	outLeaderboard := outputLeaderboard{}
	unOrderedPerf := []float64{0.3, 0.5, 0.2, 0.9, 0.4}
	boardTuples := []outputBoardTuple{}
	for _, v := range unOrderedPerf {
		boardTuples = append(boardTuples, outputBoardTuple{Perf: float32(v)})
	}
	outLeaderboard.Testtuples = boardTuples

	// Ascending order
	outLeaderboard.Sort(true)
	sort.Float64s(unOrderedPerf)
	for i, boardTuple := range outLeaderboard.Testtuples {
		assert.EqualValues(t, float32(unOrderedPerf[i]), boardTuple.Perf)
	}

	// Reverse order
	outLeaderboard.Sort(false)
	for i, j := 0, len(unOrderedPerf)-1; i < j; i, j = i+1, j-1 {
		unOrderedPerf[i], unOrderedPerf[j] = unOrderedPerf[j], unOrderedPerf[i]
	}
	for i, boardTuple := range outLeaderboard.Testtuples {
		assert.EqualValues(t, float32(unOrderedPerf[i]), boardTuple.Perf)
	}

}
