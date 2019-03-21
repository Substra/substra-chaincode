package main

// Struct use as return representation of ledger data

type outputChallenge struct {
	Key         string               `json:"key"`
	Name        string               `json:"name"`
	Description challengeDescription `json:"description"`
	Metrics     *HashDressName       `json:"metrics"`
	Owner       string               `json:"owner"`
	TestData    *DatasetData         `json:"testData"`
	Permissions string               `json:"permissions"`
}

type challengeDescription struct {
	Hash           string `json:"hash"`
	StorageAddress string `json:"storageAddress"`
}

func (out *outputChallenge) Fill(key string, in Challenge) {
	out.Key = key
	out.Name = in.Name
	out.Description.StorageAddress = in.DescriptionStorageAddress
	out.Description.Hash = key
	out.Metrics = in.Metrics
	out.Owner = in.Owner
	out.TestData = in.TestData
	out.Permissions = in.Permissions
}

// outputDataset is the return representation of the Dataset type stored in the ledger
type outputDataset struct {
	ChallengeKey string        `json:"challengeKey"`
	Description  *HashDress    `json:"description"`
	Key          string        `json:"key"`
	Name         string        `json:"name"`
	Opener       datasetOpener `json:"opener"`
	Owner        string        `json:"owner"`
	Permissions  string        `json:"permissions"`
	Type         string        `json:"type"`
}

type datasetOpener struct {
	Hash           string `json:"hash"`
	StorageAddress string `json:"strotageAddress"`
}

func (out *outputDataset) Fill(key string, in Dataset) {
	out.ChallengeKey = in.ChallengeKey
	out.Description = in.Description
	out.Key = key
	out.Name = in.Name
	out.Opener.Hash = key
	out.Opener.StorageAddress = in.OpenerStorageAddress
	out.Owner = in.Owner
	out.Permissions = in.Permissions
	out.Type = in.Type
}
