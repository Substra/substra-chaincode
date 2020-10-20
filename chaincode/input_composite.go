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
	AlgoKey                  string            `validate:"required,len=64,hexadecimal" json:"algo_key"`
	InHeadModelKey           string            `validate:"required_with=InTrunkModelKey,omitempty,len=64,hexadecimal" json:"in_head_model_key"`
	InTrunkModelKey          string            `validate:"required_with=InHeadModelKey,omitempty,len=64,hexadecimal" json:"in_trunk_model_key"`
	OutTrunkModelPermissions inputPermissions  `validate:"required" json:"out_trunk_model_permissions"`
	DataManagerKey           string            `validate:"required,len=36" json:"data_manager_key"`
	DataSampleKeys           []string          `validate:"required,unique,gt=0,dive,len=64,hexadecimal" json:"data_sample_keys"`
	ComputePlanID            string            `validate:"omitempty" json:"compute_plan_id"`
	Rank                     string            `validate:"omitempty" json:"rank"`
	Tag                      string            `validate:"omitempty,lte=64" json:"tag"`
	Metadata                 map[string]string `validate:"lte=100,dive,keys,lte=50,endkeys,lte=100" json:"metadata"`
}

type inputCompositeAlgo struct {
	inputAlgo
}

type inputLogSuccessCompositeTrain struct {
	inputLog
	OutHeadModel  inputHash      `validate:"required" json:"out_head_model"`
	OutTrunkModel inputHashDress `validate:"required" json:"out_trunk_model"`
}
