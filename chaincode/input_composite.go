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

// inputCompositeTraintuple is the representation of input args to register a composite Traintuple
type inputCompositeTraintuple struct {
	AlgoKey                  string           `validate:"required,len=64,hexadecimal" json:"algoKey"`
	InHeadModelKey           string           `validate:"required_with=InTrunkModelKey,omitempty,len=64,hexadecimal" json:"inHeadModelKey"`
	InTrunkModelKey          string           `validate:"required_with=InHeadModelKey,omitempty,len=64,hexadecimal" json:"inTrunkModelKey"`
	OutTrunkModelPermissions inputPermissions `validate:"required" json:"OutTrunkModelPermissions"`
	DataManagerKey           string           `validate:"required,len=64,hexadecimal" json:"dataManagerKey"`
	DataSampleKeys           []string         `validate:"required,unique,gt=0,dive,len=64,hexadecimal" json:"dataSampleKeys"`
	ComputePlanID            string           `validate:"omitempty" json:"computePlanID"`
	Rank                     string           `validate:"omitempty" json:"rank"`
	Tag                      string           `validate:"omitempty,lte=64" json:"tag"`
}

type inputCompositeAlgo struct {
	inputAlgo
}

type inputLogSuccessCompositeTrain struct {
	inputLog
	OutHeadModel  inputHashDress `validate:"required" json:"outHeadModel"`
	OutTrunkModel inputHashDress `validate:"required" json:"outTrunkModel"`
	Perf          float32        `validate:"omitempty" json:"perf"`
}
