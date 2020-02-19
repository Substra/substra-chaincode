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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOutModelHashDress(t *testing.T) {
	scc := new(SubstraChaincode)
	mockStub := NewMockStubWithRegisterNode("substra", scc)
	db := NewLedgerDB(mockStub)

	registerItem(t, *mockStub, "compositeAlgo")
	regular, err := registerTraintuple(mockStub, TraintupleType)
	composite, err := registerTraintuple(mockStub, CompositeTraintupleType)

	// 1. Correct requests

	_, err = db.GetOutModelHashDress(regular, []AssetType{TraintupleType})
	assert.NoError(t, err, "the regular traintuple should be found when requesting regular traintuples")

	_, err = db.GetOutHeadModelHash(composite)
	assert.NoError(t, err, "the composite traintuple should be found when requesting composite traintuples")

	_, err = db.GetOutModelHashDress(composite, []AssetType{CompositeTraintupleType})
	assert.NoError(t, err, "the composite traintuple should be found when requesting composite traintuples")

	// 2. Incorrect requests

	_, err = db.GetOutModelHashDress(regular, []AssetType{CompositeTraintupleType})
	assert.Error(t, err, "the regular traintuple should not be found when requesting composite traintuples only")

	_, err = db.GetOutModelHashDress(composite, []AssetType{TraintupleType})
	assert.Error(t, err, "the composite traintuple should be found when requesting regular traintuples only")
}

func TestGzip(t *testing.T) {

	data := "mydata"
	key := "mykey"
	txID := "mytx"

	for _, enableGzip := range []bool{false, true} {
		t.Run(fmt.Sprintf("TestGzip_%v", enableGzip), func(t *testing.T) {
			var initSettings = InitSettings{
				LedgerSettings{
					EnableGzip: enableGzip,
				},
			}
			scc := new(SubstraChaincode)
			mockStub := NewMockStubWithRegisterNode("substra", scc)
			mockStub.MockInit("42", methodAndAssetToByte("init", initSettings))
			db := NewLedgerDB(mockStub)

			// put
			mockStub.MockTransactionStart(txID)
			db.Put(key, data)
			mockStub.MockTransactionEnd(txID)

			// get
			var actual string
			db.Get(key, &actual)

			assert.Equal(t, data, actual, "Data should be retrieved from ledger successfully")
		})
	}
}
