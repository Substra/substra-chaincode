package main

import "chaincode/errors"

// TrainingTask is a node of a ComputeDAG. It represents a training task
// (i.e. a Traintuple, a CompositeTraintuple or an Aggregatetuple)
type TrainingTask struct {
	ID          string
	InModelsIDs []string
	InputIndex  int
	Depth       int
	TaskType    AssetType
}

// ComputeDAG is a Directed Acyclic Graph (DAG)
// used for compute plans
type ComputeDAG struct {
	OrderTasks []TrainingTask
}

// Create a Directed Acyclic Graph (DAG) from a compute plan
func createComputeDAG(cp inputComputePlan) (ComputeDAG, error) {
	DAG := ComputeDAG{}
	for i, traintuple := range cp.Traintuples {
		task := TrainingTask{
			ID:          traintuple.ID,
			InModelsIDs: traintuple.InModelsIDs,
			InputIndex:  i,
			TaskType:    TraintupleType,
		}
		DAG.OrderTasks = append(DAG.OrderTasks, task)
	}
	for i, traintuple := range cp.CompositeTraintuples {
		task := TrainingTask{
			ID:          traintuple.ID,
			InModelsIDs: []string{traintuple.InHeadModelID, traintuple.InTrunkModelID},
			InputIndex:  i,
			TaskType:    CompositeTraintupleType,
		}
		DAG.OrderTasks = append(DAG.OrderTasks, task)
	}
	for i, traintuple := range cp.Aggregatetuples {
		task := TrainingTask{
			ID:          traintuple.ID,
			InModelsIDs: traintuple.InModelsIDs,
			InputIndex:  i,
			TaskType:    AggregatetupleType,
		}
		DAG.OrderTasks = append(DAG.OrderTasks, task)
	}
	err := DAG.sort()
	if err != nil {
		return DAG, err
	}
	return DAG, nil
}

// Sort the DAG's task list, or return an error if there is a cyclic dependency in inModelIDs
func (dag *ComputeDAG) sort() error {
	current := dag.OrderTasks
	var temp, final []TrainingTask
	IDPresents := map[string]int{}
	for i := 0; len(current) != 0; {
		depth := 0
		ready := true
		for _, ID := range current[i].InModelsIDs {
			if ID == "" {
				continue
			}
			parentDepth, ok := IDPresents[ID]
			ready = ready && ok
			depth = max(depth, parentDepth+1)
		}
		if ready {
			current[i].Depth = depth
			final = append(final, current[i])
			if _, ok := IDPresents[current[i].ID]; ok {
				return errors.Conflict("compute plan error: Duplicate training task ID: %s", current[i].ID)
			}
			IDPresents[current[i].ID] = current[i].Depth
		} else {
			temp = append(temp, current[i])
		}
		if i != len(current)-1 {
			i++
			continue
		}
		if len(temp) == len(current) {
			var errorIDs []string
			for _, c := range current {
				errorIDs = append(errorIDs, c.ID)
			}
			return errors.Internal("compute plan error: Cyclic or missing dependency among inModels IDs: %v", errorIDs)
		}
		i = 0
		current = temp
		temp = []TrainingTask{}
	}
	dag.OrderTasks = final
	return nil
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
