package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	ts := []struct {
		name        string
		list        []TrainingTask
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
			expectError: false},
		{name: "some inModels",
			list: []TrainingTask{
				{ID: "three", InModelsIDs: []string{"one"}},
				{ID: "one"},
				{ID: "four", InModelsIDs: []string{"two", "three"}},
				{ID: "two"},
			},
			expectError: false},
		{name: "all-inModels",
			list: []TrainingTask{
				{ID: "three", InModelsIDs: []string{"two"}},
				{ID: "one"},
				{ID: "four", InModelsIDs: []string{"three", "one"}},
				{ID: "two", InModelsIDs: []string{"one"}},
			},
			expectError: false},
		{name: "wrong ID inModels",
			list: []TrainingTask{
				{ID: "one"},
				{ID: "two", InModelsIDs: []string{"one"}},
				{ID: "three", InModelsIDs: []string{"two"}},
				{ID: "four", InModelsIDs: []string{"five"}},
			},
			expectError: true,
			errorStr:    "compute plan error, either cyclic or missing dep among those IDs'inModels: [four]"},
		{name: "cyclic inModels",
			list: []TrainingTask{
				{ID: "one"},
				{ID: "two", InModelsIDs: []string{"one", "five"}},
				{ID: "three", InModelsIDs: []string{"two"}},
				{ID: "four", InModelsIDs: []string{"three"}},
				{ID: "five", InModelsIDs: []string{"four"}},
			},
			expectError: true,
			errorStr:    "compute plan error, either cyclic or missing dep among those IDs'inModels: [two three four five]"},
		{name: "Same ID twice",
			list: []TrainingTask{
				{ID: "one"},
				{ID: "two", InModelsIDs: []string{"one"}},
				{ID: "three", InModelsIDs: []string{"two"}},
				{ID: "one", InModelsIDs: []string{"three"}},
			},
			expectError: true,
			errorStr:    `compute plan error, ID use twice: one`},
	}
	for _, tc := range ts {
		t.Run(tc.name, func(t *testing.T) {
			dag := ComputeDAG{OrderTasks: tc.list}
			err := dag.sort()
			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.errorStr, err.Error())
				return
			}
			assert.NoError(t, err)
			for i, ID := range []string{"one", "two", "three", "four"} {
				assert.Equal(t, ID, dag.OrderTasks[i].ID)
			}
		})
	}
}
