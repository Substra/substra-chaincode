package main

import "fmt"

// TrainingTask is a node of a ComputeDAG. It represents a training task
// (i.e. a Traintuple, a CompositeTraintuple or an Aggregatetuple)
type TrainingTask struct {
	ID          string
	InModelsIDs []string
	InputIndex  int
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
	IDPresents := map[string]bool{}
	for i := 0; len(current) != 0; {
		ready := true
		for _, ID := range current[i].InModelsIDs {
			if ID == "" {
				continue
			}
			_, ok := IDPresents[ID]
			ready = ready && ok
		}
		if ready {
			final = append(final, current[i])
			if _, ok := IDPresents[current[i].ID]; ok {
				return fmt.Errorf("Compute plan error: Duplicate training task ID: %s", current[i].ID)
			}
			IDPresents[current[i].ID] = true
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
			return fmt.Errorf("Compute plan error: Cyclic or missing dependency among inModels IDs: %v", errorIDs)
		}
		i = 0
		current = temp
		temp = []TrainingTask{}
	}
	dag.OrderTasks = final
	return nil
}
