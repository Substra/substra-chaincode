package main

import "fmt"

func getCPWorkerState(db *LedgerDB, wStateKey string, worker string) (*ComputePlanWorkerState, error) {
	wState := ComputePlanWorkerState{}
	err := db.Get(wStateKey, &wState)
	if err != nil {
		return nil, err
	}
	return &wState, nil
}

// TryAddIntermediaryModel will reference the model key if the compute plan key
// is not empty and if it's an intermediary model meaning without any children
func TryAddIntermediaryModel(db *LedgerDB, ComputePlanKey, worker, tupleKey, modelKey string) error {
	if ComputePlanKey == "" {
		return nil
	}
	cp, err := db.GetComputePlan(ComputePlanKey)
	if err != nil {
		return err
	}
	if !cp.CleanModels {
		return nil
	}

	// Check for previously added intermediary models that might now be ready for deletion
	wStateKey := cp.getCPWorkerStateKey(worker)
	wState, err := getCPWorkerState(db, wStateKey, worker)
	if err != nil {
		return err
	}

	modelsUsed, modelsUnused, err := updateIntermediaryModelsInuse(db, wState.IntermediaryModelsInUse)

	wState.IntermediaryModelsInUse = modelsUsed

	if len(modelsUnused) != 0 {
		db.AddComputePlanEvent(ComputePlanKey, cp.State.Status, modelsUnused)
	}

	children, err := getTupleChildren(db, tupleKey)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		// If a tuple has no children it's considered final and should not be
		// listed in the index. Here we're in the other case: the tuple *does* have children.
		wState.IntermediaryModelsInUse = append(wState.IntermediaryModelsInUse, modelKey)
	}

	return db.Put(wStateKey, wState)
}

// incrementDoneCount increases the count of tuples in the "done" state for
// a given compute plan and worker
func (cp *ComputePlan) incrementDoneCount(db *LedgerDB, tupleWorker string) error {
	key := cp.getCPWorkerStateKey(tupleWorker)
	count := ComputePlanWorkerState{}
	err := db.Get(key, &count)
	if err != nil {
		return err
	}
	count.DoneCount++
	return db.Put(key, count)
}

// incrementTupleCount increases the total number of tuples for
// a given compute plan and worker
func (cp *ComputePlan) incrementTupleCount(db *LedgerDB, tupleWorker string) error {

	// Add the worker to the list of workers, if missing
	found := false
	for _, worker := range cp.Workers {
		if worker == tupleWorker {
			found = true
			break
		}
	}
	if !found {
		cp.Workers = append(cp.Workers, tupleWorker)
	}

	// Create or update the done count
	key := cp.getCPWorkerStateKey(tupleWorker)
	count := ComputePlanWorkerState{}
	err := db.Get(key, &count)

	if err != nil {
		return db.Add(key, ComputePlanWorkerState{TupleCount: 1})
	}

	count.TupleCount++
	return db.Put(key, count)
}

// getWorkerCount returns the count of tuples in the "done" state for a given
// compute plan and worker
func (cp *ComputePlan) getWorkerCount(db *LedgerDB, tupleWorker string) (ComputePlanWorkerState, error) {
	key := cp.getCPWorkerStateKey(tupleWorker)
	doneCount := ComputePlanWorkerState{}

	err := db.Get(key, &doneCount)
	return doneCount, err
}

func (cp *ComputePlan) getTotalCount(db *LedgerDB) (count ComputePlanWorkerState, err error) {
	for _, worker := range cp.Workers {
		wCount, err := cp.getWorkerCount(db, worker)
		if err != nil {
			return count, err
		}
		count.TupleCount += wCount.TupleCount
		count.DoneCount += wCount.DoneCount
	}
	return count, nil
}

func (cp *ComputePlan) getCPWorkerStateKey(tupleWorker string) string {
	return fmt.Sprintf("computePlan~%v~doneCountByWorker~%v", cp.Key, tupleWorker)
}

func modelIsInUse(db *LedgerDB, modelKey string) (bool, error) {
	keys, err := db.GetIndexKeys("tuple~modelKey~key", []string{"tuple", modelKey})
	if err != nil {
		return false, err
	}
	if len(keys) == 0 {
		// This occurs for the keys added during the same transaction. But
		// thoses models can just be considered in use
		return true, nil
	}
	children, err := getTupleChildren(db, keys[0])
	if err != nil {
		return false, nil
	}
	for _, tupleKey := range children {
		tuple, err := db.GetGenericTuple(tupleKey)
		if err != nil {
			return false, nil
		}
		if tuple.Status != StatusDone {
			return true, nil
		}
	}
	return false, nil
}

func updateIntermediaryModelsInuse(db *LedgerDB, oldModelsUsed []string) (modelsUsed []string, modelsUnused []string, err error) {
	modelsUsed = []string{}
	modelsUnused = []string{}
	for _, modelKey := range oldModelsUsed {
		inUse, err := modelIsInUse(db, modelKey)
		if err != nil {
			return []string{}, []string{}, err
		}
		if inUse {
			modelsUsed = append(modelsUsed, modelKey)
		} else {
			modelsUnused = append(modelsUnused, modelKey)
		}
	}
	return modelsUsed, modelsUnused, nil
}

func removeAllIntermediaryModels(db *LedgerDB, cp *ComputePlan) ([]string, error) {
	res := []string{}
	for _, worker := range cp.Workers {
		wStateKey := cp.getCPWorkerStateKey(worker)
		wState, err := getCPWorkerState(db, wStateKey, worker)
		if err != nil {
			return []string{}, err
		}
		res = append(res, wState.IntermediaryModelsInUse...)
		wState.IntermediaryModelsInUse = []string{}

		// clear
		err = db.Put(wStateKey, wState)
		if err != nil {
			return []string{}, err
		}
	}
	return res, nil
}
