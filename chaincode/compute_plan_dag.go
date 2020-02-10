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
	OrderTasks    []TrainingTask
	IDToTrainTask map[string]CPTrainTask
}

// Create a Directed Acyclic Graph (DAG) from a compute plan
func createComputeDAG(cp inputComputePlan, IDToTrainTask map[string]CPTrainTask) (ComputeDAG, error) {
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
	DAG.IDToTrainTask = IDToTrainTask
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
	if dag.IDToTrainTask == nil {
		dag.IDToTrainTask = make(map[string]CPTrainTask)
	}
	for i := 0; len(current) != 0; {
		depth := 0
		ready := true
		for _, ID := range current[i].InModelsIDs {
			if ID == "" {
				continue
			}
			parent, ok := dag.IDToTrainTask[ID]
			ready = ready && ok
			if !ok {
				break
			}
			depth = max(depth, parent.Depth+1)
		}
		if ready {
			current[i].Depth = depth
			final = append(final, current[i])
			if _, ok := dag.IDToTrainTask[current[i].ID]; ok {
				return errors.BadRequest("compute plan error: Duplicate training task ID: %s", current[i].ID)
			}
			dag.IDToTrainTask[current[i].ID] = CPTrainTask{Depth: current[i].Depth}
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
			return errors.BadRequest("compute plan error: Cyclic or missing dependency among inModels IDs: %v", errorIDs)
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
