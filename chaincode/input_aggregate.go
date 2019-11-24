// Copyright 2018 Owkin, inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// inputAggregateTuple is the representation of input args to register an aggregate Tuple
type inputAggregateTuple struct {
	AlgoKey       string   `validate:"required,len=64,hexadecimal" json:"algoKey"`
	ObjectiveKey  string   `validate:"required,len=64,hexadecimal" json:"objectiveKey"`
	InModels      []string `validate:"omitempty,dive,len=64,hexadecimal" json:"inModels"`
	ComputePlanID string   `validate:"omitempty" json:"computePlanID"`
	Rank          string   `validate:"omitempty" json:"rank"`
	Tag           string   `validate:"omitempty,lte=64" json:"tag"`
	Worker        string   `validate:"required" json:"worker"`
}

type inputAggregateAlgo struct {
	inputAlgo
}
