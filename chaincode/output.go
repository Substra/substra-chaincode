package main

// Struct use as return representation of ledger data

type outputObjective struct {
	Key         string         `json:"key"`
	Name        string         `json:"name"`
	Description HashDress      `json:"description"`
	Metrics     *HashDressName `json:"metrics"`
	Owner       string         `json:"owner"`
	TestDataset *Dataset       `json:"testDataset"`
	Permissions string         `json:"permissions"`
}

func (out *outputObjective) Fill(key string, in Objective) {
	out.Key = key
	out.Name = in.Name
	out.Description.StorageAddress = in.DescriptionStorageAddress
	out.Description.Hash = key
	out.Metrics = in.Metrics
	out.Owner = in.Owner
	out.TestDataset = in.TestDataset
	out.Permissions = in.Permissions
}

// outputDataManager is the return representation of the DataManager type stored in the ledger
type outputDataManager struct {
	ObjectiveKey string     `json:"objectiveKey"`
	Description  *HashDress `json:"description"`
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Opener       HashDress  `json:"opener"`
	Owner        string     `json:"owner"`
	Permissions  string     `json:"permissions"`
	Type         string     `json:"type"`
}

func (out *outputDataManager) Fill(key string, in DataManager) {
	out.ObjectiveKey = in.ObjectiveKey
	out.Description = in.Description
	out.Key = key
	out.Name = in.Name
	out.Opener.Hash = key
	out.Opener.StorageAddress = in.OpenerStorageAddress
	out.Owner = in.Owner
	out.Permissions = in.Permissions
	out.Type = in.Type
}

type outputAlgo struct {
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	Content      HashDress  `json:"content"`
	Description  *HashDress `json:"description"`
	Owner        string     `json:"owner"`
	ObjectiveKey string     `json:"objectiveKey"`
	Permissions  string     `json:"permissions"`
}

func (out *outputAlgo) Fill(key string, in Algo) {
	out.Key = key
	out.Name = in.Name
	out.Content.Hash = key
	out.Content.StorageAddress = in.StorageAddress
	out.Description = in.Description
	out.Owner = in.Owner
	out.ObjectiveKey = in.ObjectiveKey
	out.Permissions = in.Permissions
}
