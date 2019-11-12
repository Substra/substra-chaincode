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

import (
	"fmt"
	"strings"
)

// queryFilter returns all elements of the ledger matching some filters
// For now, ok for everything. Later returns if the requester has permission to see it
func queryFilter(db LedgerDB, args []string) (elements interface{}, err error) {
	inp := inputQueryFilter{}
	err = AssetFromJSON(args, &inp)
	if err != nil {
		return
	}
	// check validity of inputs
	validIndexNames := []string{
		"traintuple~worker~status",
		"testtuple~worker~status",
		"testtuple~tag",
		"traintuple~tag",
		"compositeTraintuple~worker~status",
		"compositeTraintuple~tag",
	}
	if !stringInSlice(inp.IndexName, validIndexNames) {
		err = fmt.Errorf("invalid indexName filter query: %s", inp.IndexName)
		return
	}
	indexName := inp.IndexName + "~key"
	attributes := strings.Split(strings.Replace(inp.Attributes, " ", "", -1), ",")
	attributes = append([]string{strings.Split(indexName, "~")[0]}, attributes...)

	filteredKeys, err := db.GetIndexKeys(indexName, attributes)
	if err != nil {
		return
	}
	// get elements with filtererd keys
	switch indexName {
	case "testtuple~worker~status~key", "testtuple~tag~key":
		elements, err = getOutputTesttuples(db, filteredKeys)
	case "traintuple~worker~status~key", "traintuple~tag~key":
		elements, err = getOutputTraintuples(db, filteredKeys)
	case "compositeTraintuple~worker~status~key", "compositeTraintuple~tag~key":
		elements, err = getOutputCompositeTraintuples(db, filteredKeys)
	}
	return
}
