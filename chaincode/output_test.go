// Copyright 2018 Owkin, inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboardSort(t *testing.T) {
	unOrderedPerf := []float64{0.3, 0.5, 0.2, 0.9, 0.4}
	boardTuples := outputBoardTuples{}
	for _, v := range unOrderedPerf {
		boardTuples = append(boardTuples, outputBoardTuple{Perf: float32(v)})
	}

	// Ascending order
	sort.Sort(boardTuples)
	sort.Float64s(unOrderedPerf)
	for i, boardTuple := range boardTuples {
		assert.EqualValues(t, float32(unOrderedPerf[i]), boardTuple.Perf)
	}

	// Reverse order
	sort.Sort(sort.Reverse(boardTuples))
	for i, j := 0, len(unOrderedPerf)-1; i < j; i, j = i+1, j-1 {
		unOrderedPerf[i], unOrderedPerf[j] = unOrderedPerf[j], unOrderedPerf[i]
	}
	for i, boardTuple := range boardTuples {
		assert.EqualValues(t, float32(unOrderedPerf[i]), boardTuple.Perf)
	}
}
