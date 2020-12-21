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

package permissions

import (
	"chaincode/errors"
	"chaincode/io"
	"chaincode/ledger"
	"chaincode/utils"
)

// Permission represents one permission based on an action type
type Permission struct {
	// Public is true if this permission is given to the asset's owner only and
	// the nodes listed in AuthorizedIDs (open to all nodes if false)
	Public bool `json:"public"`
	// AuthorizedIDs list all authorised nodes other than the asset's owner
	AuthorizedIDs []string `json:"authorized_ids"`
}

// Permissions represents all permissions associated with an asset
type Permissions struct {
	// Download define if a given node can allow its nodes to download the asset
	Download Permission `json:"download"`
	// Process define if a given node can process the sset
	Process Permission `json:"process"`
}

// CanProcess checks if a node can process the asset with the current permissions
func (perms Permissions) CanProcess(owner, node string) bool {
	if owner == node {
		return true
	}

	if perms.Process.Public {
		return true
	}

	for _, authorizedNode := range perms.Process.AuthorizedIDs {
		if node == authorizedNode {
			return true
		}
	}
	return false
}

// NewPermissions create the Permissions according to the arg received
func NewPermissions(db *ledger.LedgerDB, in io.InputPermissions) (Permissions, error) {
	nodes, err := queryNodes(db, []string{})
	if err != nil {
		return Permissions{}, err
	}

	nodesIDs := []string{}
	for _, node := range nodes {
		nodesIDs = append(nodesIDs, node.ID)
	}

	// Validate Process inputPermissions
	// @TODO Validate Download inputPermissions when implemented
	for _, authorizedID := range in.Process.AuthorizedIDs {
		if in.Process.Public {
			continue
		}

		if !utils.StringInSlice(authorizedID, nodesIDs) {
			return Permissions{}, errors.BadRequest("invalid permission input values")
		}
	}

	owner, err := utils.GetTxCreator(db.cc)
	if err != nil {
		return Permissions{}, err
	}

	permissions := Permissions{}
	process := newPermission(in.Process, owner)
	permissions.Process = process
	// Download permission is not implemented in the node server, so it is set to the process permission
	permissions.Download = process
	return permissions, nil
}

func newPermission(in io.InputPermission, owner string) Permission {
	// Owner must always be defined in the list of authorizedIDs, if the permission is private,
	// it will ease the merge of private permissions
	if !utils.StringInSlice(owner, in.AuthorizedIDs) {
		in.AuthorizedIDs = append([]string{owner}, in.AuthorizedIDs...)
	}
	return Permission(in)
}

func (priv Permission) include(other Permission) bool {
	if priv.Public {
		return true
	}
	if other.Public {
		return false
	}
	for _, node := range other.AuthorizedIDs {
		if !utils.StringInSlice(node, priv.AuthorizedIDs) {
			return false
		}
	}
	return true
}

// MergePermissions returns the intersection of input permissions
func MergePermissions(x, y Permissions) Permissions {
	perm := Permissions{}
	perm.Process = mergePermissions(x.Process, y.Process)
	perm.Download = mergePermissions(x.Download, y.Download)
	return perm
}

func mergePermissions(x, y Permission) Permission {
	priv := Permission{}
	priv.Public = x.Public && y.Public

	switch {
	case !x.Public && y.Public:
		priv.AuthorizedIDs = x.AuthorizedIDs
	case x.Public && !y.Public:
		priv.AuthorizedIDs = y.AuthorizedIDs
	default:
		priv.AuthorizedIDs = x.getNodesIntersection(y)
	}
	return priv
}

func (priv Permission) getNodesIntersection(p Permission) []string {
	nodes := []string{}
	for _, i := range priv.AuthorizedIDs {
		for _, j := range p.AuthorizedIDs {
			if i == j {
				nodes = append(nodes, i)
				break
			}
		}
	}
	return nodes
}
