package main

import "fmt"

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
	wState, err := db.GetCPWorkerState(wStateKey)
	if err != nil {
		return err
	}

	modelsUsed, modelsUnused, err := getModelsInUse(db, wState.IntermediaryModelsInUse)

	wState.IntermediaryModelsInUse = modelsUsed

	if len(modelsUnused) != 0 {
		db.AddComputePlanEvent(cp.Key, cp.State.Status, modelsUnused)
	}

	children, err := getTupleChildren(db, tupleKey, false)
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

// incrementWorkerTupleCount increases the total number of tuples for
// a given compute plan and worker
func (cp *ComputePlan) incrementWorkerTupleCount(db *LedgerDB, worker string) error {

	// Add the worker to the list of workers, if missing
	found := false
	for _, w := range cp.Workers {
		if w == worker {
			found = true
			break
		}
	}
	if !found {
		cp.Workers = append(cp.Workers, worker)
	}

	// Create or update the done count
	wStateKey := cp.getCPWorkerStateKey(worker)
	wState, err := db.GetCPWorkerState(wStateKey)
	if err != nil {
		return db.Add(wStateKey, ComputePlanWorkerState{TupleCount: 1})
	}

	wState.TupleCount++
	return db.Put(wStateKey, wState)
}

// incrementWorkerDoneCount increases the count of tuples in the "done" state for
// a given compute plan and worker
func (cp *ComputePlan) incrementWorkerDoneCount(db *LedgerDB, worker string) error {
	wStateKey := cp.getCPWorkerStateKey(worker)
	wState, err := db.GetCPWorkerState(wStateKey)
	if err != nil {
		return nil
	}
	wState.DoneCount++
	return db.Put(wStateKey, wState)
}

// getTupleCounts returns the number of tuples in the "done" state and the total number of tuples
// for a given compute plan
func (cp *ComputePlan) getTupleCounts(db *LedgerDB) (doneCount int, tupleCount int, err error) {
	for _, worker := range cp.Workers {
		wStateKey := cp.getCPWorkerStateKey(worker)
		wState, err := db.GetCPWorkerState(wStateKey)
		if err != nil {
			return doneCount, tupleCount, nil
		}
		doneCount += wState.DoneCount
		tupleCount += wState.TupleCount
	}
	return doneCount, tupleCount, nil
}

// getCPWorkerStateKey returns the worker state key for a given compute plan and worker
func (cp *ComputePlan) getCPWorkerStateKey(worker string) string {
	return fmt.Sprintf("computePlan~%v~doneCountByWorker~%v", cp.Key, worker)
}

// removeAllIntermediaryModels iterates through all the worker states, and clears the lists of
// intermediary models. It returns the concatenated list of all the intermediary models that
// have been removed from the worker states.
func (cp *ComputePlan) removeAllIntermediaryModels(db *LedgerDB) ([]string, error) {
	res := []string{}
	for _, worker := range cp.Workers {
		wStateKey := cp.getCPWorkerStateKey(worker)
		wState, err := db.GetCPWorkerState(wStateKey)
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

// getModelsInUse takes a list of model keys and checks whether these models are still in use or not.
// It returns the initial list split into two sublists: the models in use, and the models unused.
func getModelsInUse(db *LedgerDB, modelKeys []string) (usedModels []string, unusedModels []string, err error) {
	usedModels = []string{}
	unusedModels = []string{}
	for _, modelKey := range modelKeys {
		inUse, err := isModelInUse(db, modelKey)
		if err != nil {
			return []string{}, []string{}, err
		}
		if inUse {
			usedModels = append(usedModels, modelKey)
		} else {
			unusedModels = append(unusedModels, modelKey)
		}
	}
	return usedModels, unusedModels, nil
}

// isModelInUse returns true if the model with the supplied key is the
// in-model of another training task that isn't in the "done" state.
func isModelInUse(db *LedgerDB, modelKey string) (bool, error) {
	keys, err := db.GetIndexKeys("tuple~modelKey~key", []string{"tuple", modelKey})
	if err != nil {
		return false, err
	}
	if len(keys) == 0 {
		// This occurs for the keys added during the same transaction. But
		// thoses models can just be considered in use
		return true, nil
	}
	children, err := getTupleChildren(db, keys[0], true)
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

// getTupleChildren returns the keys of all the tuple which have the supplied tuple as an in-model
// If includeTesttuples is True, return the testtuple children, else omit them,
func getTupleChildren(db *LedgerDB, tupleKey string, includeTesttuples bool) ([]string, error) {
	tupleChildrenKeys, err := db.GetIndexKeys("tuple~inModel~key", []string{"tuple", tupleKey})
	if err != nil {
		return []string{}, err
	}

	if !includeTesttuples {
		return tupleChildrenKeys, nil
	}

	testtupleChildrenKeys, err := db.GetIndexKeys("testtuple~traintuple~certified~key", []string{"testtuple", tupleKey})
	if err != nil {
		return []string{}, err
	}
	return append(tupleChildrenKeys, testtupleChildrenKeys...), nil
}
