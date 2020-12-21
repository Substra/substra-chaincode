package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDAGSort(t *testing.T) {
	ts := []struct {
		name        string
		list        []TrainingTask
		IDToDepth   map[string]TrainTask
		depths      []int
		expectError bool
		errorStr    string
	}{
		{name: "no inModels",
			list: []TrainingTask{
				{ID: "one"},
				{ID: "two"},
				{ID: "three"},
				{ID: "four"},
			},
			depths:      []int{0, 0, 0, 0},
			expectError: false},
		{name: "some inModels",
			list: []TrainingTask{
				{ID: "three", InModelsIDs: []string{"one"}},
				{ID: "one"},
				{ID: "four", InModelsIDs: []string{"two", "three"}},
				{ID: "two"},
			},
			depths:      []int{0, 0, 1, 2},
			expectError: false},
		{name: "all-inModels",
			list: []TrainingTask{
				{ID: "three", InModelsIDs: []string{"two"}},
				{ID: "one"},
				{ID: "four", InModelsIDs: []string{"three", "one"}},
				{ID: "two", InModelsIDs: []string{"one"}},
			},
			depths:      []int{0, 1, 2, 3},
			expectError: false},
		{name: "wrong ID inModels",
			list: []TrainingTask{
				{ID: "one"},
				{ID: "two", InModelsIDs: []string{"one"}},
				{ID: "three", InModelsIDs: []string{"two"}},
				{ID: "four", InModelsIDs: []string{"five"}},
			},
			expectError: true,
			errorStr:    "compute plan error: Cyclic or missing dependency among inModels IDs: [four]"},
		{name: "cyclic inModels",
			list: []TrainingTask{
				{ID: "one"},
				{ID: "two", InModelsIDs: []string{"one", "five"}},
				{ID: "three", InModelsIDs: []string{"two"}},
				{ID: "four", InModelsIDs: []string{"three"}},
				{ID: "five", InModelsIDs: []string{"four"}},
			},
			expectError: true,
			errorStr:    "compute plan error: Cyclic or missing dependency among inModels IDs: [two three four five]"},
		{name: "Same ID twice",
			list: []TrainingTask{
				{ID: "one"},
				{ID: "two", InModelsIDs: []string{"one"}},
				{ID: "three", InModelsIDs: []string{"two"}},
				{ID: "one", InModelsIDs: []string{"three"}},
			},
			expectError: true,
			errorStr:    `compute plan error: Duplicate training task ID: one`},
		{name: "with existing IDs",
			list: []TrainingTask{
				{ID: "three", InModelsIDs: []string{"two", "beta"}},
				{ID: "one", InModelsIDs: []string{"alpha"}},
				{ID: "four", InModelsIDs: []string{"three", "one"}},
				{ID: "two", InModelsIDs: []string{"one"}},
			},
			IDToDepth: map[string]TrainTask{
				"alpha": TrainTask{Depth: 0},
				"beta":  TrainTask{Depth: 4},
			},
			depths:      []int{1, 2, 5, 6},
			expectError: false},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			dag := ComputeDAG{
				OrderTasks:    tc.list,
				IDToTrainTask: tc.IDToDepth,
			}
			err := dag.sort()
			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.errorStr, err.Error())
				return
			}
			assert.NoError(t, err)
			for i, ID := range []string{"one", "two", "three", "four"} {
				assert.Equal(t, ID, dag.OrderTasks[i].ID)
				assert.Equal(t, tc.depths[i], dag.OrderTasks[i].Depth)
			}
		})
	}
}
