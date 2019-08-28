package main

// Action is the the type of an action
type Action string

// Enum the different action types available
const (
	Process  Action = "process"
	Download Action = "download"
)

// Privilege represents one permission based on an action type
type Privilege struct {
	// Public is true if this permission is given to the asset's owner only and
	// the nodes listed in AuthorizedIDs (open to all nodes if false)
	Public bool `json:"public"`
	// AuthorizedIDs list all authorised nodes other than the asset's owner
	AuthorizedIDs []string `json:"authorizedIDs"`
}

// Permissions represents all privileges associated with an asset
type Permissions struct {
	// Download define if a given node can allow its nodes to download the asset
	Download Privilege `json:"download"`
	// Process define if a given node can process the sset
	Process Privilege `json:"process"`
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

// NewPermissions the Permissions Privilege according to the arg received
func NewPermissions(in inputPermissions, owner string) Permissions {
	perms := Permissions{}
	process := newPrivilege(in.Process, owner)
	perms.Process = process
	// Download privilege is not implemented in the node server, so it is set to the process privilege
	perms.Download = process
	return perms
}

func newPrivilege(in inputPrivilege, owner string) Privilege {
	// Owner must always be defined in the list of authorizedIDs if the privilege is private,
	// it will ease the merge of private privileges
	if !in.Public && !stringInSlice(owner, in.AuthorizedIDs) {
		in.AuthorizedIDs = append([]string{owner}, in.AuthorizedIDs...)
	}
	return Privilege{
		Public:        in.Public,
		AuthorizedIDs: in.AuthorizedIDs,
	}
}

func (priv Privilege) include(other Privilege) bool {
	if priv.Public {
		return true
	}
	if other.Public {
		return false
	}
	for _, node := range other.AuthorizedIDs {
		if !stringInSlice(node, priv.AuthorizedIDs) {
			return false
		}
	}
	return true
}

// MergePermissions returns the intersection of input permissions
func MergePermissions(x, y Permissions) Permissions {
	perm := Permissions{}
	perm.Process = mergePrivileges(x.Process, y.Process)
	perm.Download = mergePrivileges(x.Download, y.Download)
	return perm
}

func mergePrivileges(x, y Privilege) Privilege {
	priv := Privilege{}
	priv.Public = x.Public && y.Public

	if !x.Public && y.Public {
		priv.AuthorizedIDs = x.AuthorizedIDs
	} else if x.Public && !y.Public {
		priv.AuthorizedIDs = y.AuthorizedIDs
	} else {
		priv.AuthorizedIDs = x.getNodesIntersection(y)
	}
	return priv
}

func (priv Privilege) getNodesIntersection(p Privilege) []string {
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
