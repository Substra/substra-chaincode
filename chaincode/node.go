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
	"chaincode/errors"
)

func registerNode(db *LedgerDB, args []string) (Node, error) {
	txCreator, err := GetTxCreator(db.cc)
	if err != nil {
		return Node{}, err
	}

	node := Node{}
	node.ID = txCreator

	// Not using db.Add because we need to handle conflict as silent event without errors
	exists, err := db.KeyExists(node.ID)
	if err != nil {
		return Node{}, err
	}

	if exists {
		return node, nil
	}

	err = db.Put(node.ID, node)
	if err != nil {
		return Node{}, err
	}

	err = db.CreateIndex("node~key", []string{"node", node.ID})
	if err != nil {
		return Node{}, err
	}

	return node, nil
}

func queryNodes(db *LedgerDB, args []string) (nodes []Node, err error) {
	nodes = []Node{}

	elementsKeys, err := db.GetIndexKeys("node~key", []string{"node"})

	if len(args) != 0 {
		err = errors.BadRequest("incorrect number of arguments, expecting at most one argument")
		return
	}

	if err != nil {
		return
	}

	for _, key := range elementsKeys {
		node, err := db.GetNode(key)
		if err != nil {
			return nodes, err
		}

		nodes = append(nodes, node)
	}

	return
}
