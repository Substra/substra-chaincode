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
